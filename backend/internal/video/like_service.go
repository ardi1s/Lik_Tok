package video

import (
	"Lik_tok/internal/middleware/rabbitmq"
	rediscache "Lik_tok/internal/middleware/redis"
	"context"
	"fmt"
	"log"
)

type LikeService struct {
	repo         *LikeRepository
	videoRepo    *VideoRepository
	cache        *rediscache.Client
	likeMQ       *rabbitmq.LikeMQ
	popularityMQ *rabbitmq.PopularityMQ
}

func NewLikeService(repo *LikeRepository, videoRepo *VideoRepository, cache *rediscache.Client, likeMQ *rabbitmq.LikeMQ, popularityMQ *rabbitmq.PopularityMQ) *LikeService {
	return &LikeService{
		repo:         repo,
		videoRepo:    videoRepo,
		cache:        cache,
		likeMQ:       likeMQ,
		popularityMQ: popularityMQ,
	}
}

func (s *LikeService) Like(ctx context.Context, like *Like) error {
	// 先写入数据库
	if err := s.repo.Like(ctx, like); err != nil {
		log.Printf("LikeService: Like failed: %v", err)
		return err
	}
	log.Printf("LikeService: Like success, videoID=%d, accountID=%d", like.VideoID, like.AccountID)
	// 直接更新点赞数和热度
	if err := s.videoRepo.ChangeLikesCount(ctx, like.VideoID, 1); err != nil {
		log.Printf("LikeService: ChangeLikesCount failed: %v", err)
		return err
	}
	log.Printf("LikeService: ChangeLikesCount success")
	// 清除视频缓存
	if s.cache != nil {
		entityKey := fmt.Sprintf("video:entity:%d", like.VideoID)
		_ = s.cache.Del(ctx, entityKey)
		log.Printf("LikeService: cache cleared: %s", entityKey)
	}
	// 发送 MQ 消息给 popularity worker
	if s.popularityMQ != nil {
		s.popularityMQ.Update(ctx, like.VideoID, 1)
	}
	return nil
}

func (s *LikeService) Unlike(ctx context.Context, like *Like) error {
	// 先从数据库删除
	if err := s.repo.Unlike(ctx, like); err != nil {
		return err
	}
	// 直接更新点赞数和热度
	if err := s.videoRepo.ChangeLikesCount(ctx, like.VideoID, -1); err != nil {
		return err
	}
	// 清除视频缓存
	if s.cache != nil {
		entityKey := fmt.Sprintf("video:entity:%d", like.VideoID)
		_ = s.cache.Del(ctx, entityKey)
	}
	// 发送 MQ 消息给 popularity worker
	if s.popularityMQ != nil {
		s.popularityMQ.Update(ctx, like.VideoID, -1)
	}
	return nil
}

func (s *LikeService) IsLiked(ctx context.Context, videoID, accountID uint) (bool, error) {
	return s.repo.IsLiked(ctx, videoID, accountID)
}

func (s *LikeService) ListLikedVideos(ctx context.Context, accountID uint) ([]Video, error) {
	return s.repo.ListLikedVideos(ctx, accountID)
}
