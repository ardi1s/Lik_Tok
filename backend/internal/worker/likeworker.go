package worker

import (
	"Lik_tok/internal/middleware/rabbitmq"
	rediscache "Lik_tok/internal/middleware/redis"
	"Lik_tok/internal/video"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type LikeWorker struct {
	ch     *amqp.Channel
	likes  *video.LikeRepository
	videos *video.VideoRepository
	cache  *rediscache.Client
	queue  string
}

func NewLikeWorker(ch *amqp.Channel, likes *video.LikeRepository, videos *video.VideoRepository, cache *rediscache.Client, queue string) *LikeWorker {
	return &LikeWorker{ch: ch, likes: likes, videos: videos, cache: cache, queue: queue}
}

func (w *LikeWorker) Run(ctx context.Context) error {
	if w == nil || w.ch == nil || w.likes == nil || w.videos == nil {
		return errors.New("like worker is not initialized")
	}
	if w.queue == "" {
		return errors.New("queue is required")
	}

	deliveries, err := w.ch.Consume(
		w.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-deliveries:
			if !ok {
				return errors.New("deliveries channel closed")
			}
			w.handleDelivery(ctx, d)
		}
	}
}

func (w *LikeWorker) handleDelivery(ctx context.Context, d amqp.Delivery) {
	if err := w.process(ctx, d.Body); err != nil {
		log.Printf("like worker: failed to process message: %v", err)
		_ = d.Nack(false, true)
		return
	}
	_ = d.Ack(false)
}

func (w *LikeWorker) process(ctx context.Context, body []byte) error {
	var evt rabbitmq.LikeEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		// 解析事件失败，直接丢弃
		log.Printf("like worker: failed to unmarshal event: %v", err)
		return nil
	}
	log.Printf("like worker: received event: action=%s, userID=%d, videoID=%d", evt.Action, evt.UserID, evt.VideoID)
	if evt.UserID == 0 || evt.VideoID == 0 {
		log.Printf("like worker: invalid event, skipping")
		return nil
	}

	switch evt.Action {
	case "like":
		return w.applyLike(ctx, evt.UserID, evt.VideoID)
	case "unlike":
		return w.applyUnlike(ctx, evt.UserID, evt.VideoID)
	default:
		log.Printf("like worker: unknown action: %s", evt.Action)
		return nil
	}
}

func (w *LikeWorker) applyLike(ctx context.Context, userID, videoID uint) error {
	ok, err := w.videos.IsExist(ctx, videoID)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	created, err := w.likes.LikeIgnoreDuplicate(ctx, &video.Like{
		VideoID:   videoID,
		AccountID: userID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	// 只有新创建的记录才更新点赞数和热度
	if !created {
		return nil
	}

	if err := w.videos.ChangeLikesCount(ctx, videoID, 1); err != nil {
		return err
	}
	// 清除视频缓存
	w.clearVideoCache(videoID)
	return w.videos.ChangePopularity(ctx, videoID, 1)
}

func (w *LikeWorker) applyUnlike(ctx context.Context, userID, videoID uint) error {
	log.Printf("like worker: applyUnlike started: userID=%d, videoID=%d", userID, videoID)
	ok, err := w.videos.IsExist(ctx, videoID)
	if err != nil {
		log.Printf("like worker: IsExist failed: %v", err)
		return err
	}
	if !ok {
		log.Printf("like worker: video %d does not exist", videoID)
		return nil
	}

	deleted, err := w.likes.DeleteByVideoAndAccount(ctx, videoID, userID)
	if err != nil {
		log.Printf("like worker: DeleteByVideoAndAccount failed: %v", err)
		return err
	}
	if !deleted {
		log.Printf("like worker: no like record found for user %d and video %d", userID, videoID)
		return nil
	}

	log.Printf("like worker: like record deleted, updating likes_count")
	if err := w.videos.ChangeLikesCount(ctx, videoID, -1); err != nil {
		log.Printf("like worker: ChangeLikesCount failed: %v", err)
		return err
	}
	// 清除视频缓存
	w.clearVideoCache(videoID)
	log.Printf("like worker: applyUnlike completed successfully")
	return w.videos.ChangePopularity(ctx, videoID, -1)
}

func (w *LikeWorker) clearVideoCache(videoID uint) {
	if w.cache == nil {
		return
	}
	// 清除视频实体缓存
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	entityKey := fmt.Sprintf("video:entity:%d", videoID)
	_ = w.cache.Del(ctx, entityKey)
}
