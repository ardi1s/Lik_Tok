package video

import (
	"Lik_tok/internal/middleware/rabbitmq"
	rediscache "Lik_tok/internal/middleware/redis"
	"context"
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
	return s.repo.Like(ctx, like)
}

func (s *LikeService) Unlike(ctx context.Context, like *Like) error {
	return s.repo.Unlike(ctx, like)
}

func (s *LikeService) IsLiked(ctx context.Context, videoID, accountID uint) (bool, error) {
	return s.repo.IsLiked(ctx, videoID, accountID)
}

func (s *LikeService) ListLikedVideos(ctx context.Context, accountID uint) ([]Video, error) {
	return s.repo.ListLikedVideos(ctx, accountID)
}
