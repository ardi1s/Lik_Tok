package rabbitmq

import (
	"context"
	"errors"
	"time"
)

type TimelineMQ struct {
	*RabbitMQ
}

const (
	timelineExchange   = "timeline.events"
	timelineQueue      = "timeline.events"
	timelineBindingKey = "timeline.*"

	timelineVideoRK = "timeline.video"
)

type TimelineEvent struct {
	VideoID    uint      `json:"video_id"`
	CreateTime time.Time `json:"create_time"`
}

func NewTimelineMQ(base *RabbitMQ) (*TimelineMQ, error) {
	if base == nil {
		return nil, errors.New("rabbitmq base is nil")
	}
	if err := base.DeclareTopic(timelineExchange, timelineQueue, timelineBindingKey); err != nil {
		return nil, err
	}
	return &TimelineMQ{RabbitMQ: base}, nil
}

func (t *TimelineMQ) PublishVideo(ctx context.Context, videoID uint, createTime time.Time) error {
	if t == nil || t.RabbitMQ == nil {
		return errors.New("timeline mq is not initialized")
	}
	if videoID == 0 {
		return errors.New("videoID is required")
	}
	event := TimelineEvent{
		VideoID:    videoID,
		CreateTime: createTime,
	}
	return t.PublishJSON(ctx, timelineExchange, timelineVideoRK, event)
}
