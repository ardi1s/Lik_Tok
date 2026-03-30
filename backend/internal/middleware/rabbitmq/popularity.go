package rabbitmq

import (
	"context"
	"errors"
)

type PopularityMQ struct {
	*RabbitMQ
}

const (
	popularityExchange   = "popularity.events"
	popularityQueue      = "popularity.events"
	popularityBindingKey = "popularity.*"

	popularityUpdateRK = "popularity.update"
)

type PopularityEvent struct {
	VideoID uint `json:"video_id"`
	Change  int  `json:"change"`
}

func NewPopularityMQ(base *RabbitMQ) (*PopularityMQ, error) {
	if base == nil {
		return nil, errors.New("rabbitmq base is nil")
	}
	if err := base.DeclareTopic(popularityExchange, popularityQueue, popularityBindingKey); err != nil {
		return nil, err
	}
	return &PopularityMQ{RabbitMQ: base}, nil
}

func (p *PopularityMQ) Update(ctx context.Context, videoID uint, delta int) error {
	if p == nil || p.RabbitMQ == nil {
		return errors.New("popularity mq is not initialized")
	}
	if videoID == 0 {
		return errors.New("videoID is required")
	}
	event := PopularityEvent{
		VideoID: videoID,
		Change:  delta,
	}
	return p.PublishJSON(ctx, popularityExchange, popularityUpdateRK, event)
}
