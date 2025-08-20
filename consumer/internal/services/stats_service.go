package services

import (
	"consumer/internal/models"
	"consumer/internal/repositories"
	"consumer/internal/utils"
	"context"

	"github.com/pkg/errors"
)

var (
	ErrInterval = errors.New("erroneous interval received")
)

type StatsService struct {
	repo repositories.StatsRepo
}

func NewStatsService(r repositories.StatsRepo) *StatsService {
	return &StatsService{r}
}

func (s *StatsService) GetStats(ctx context.Context, key string) (*models.Stats, error) {
	return s.repo.GetStats(ctx, key)
}

func (s *StatsService) ProcessSwapEvent(ctx context.Context, event models.SwapEvent) error {
	err := s.repo.UpsertStats(ctx, event.TokenFrom, event.UsdValue)
	if err != nil {
		return errors.Wrapf(err, "failed to upsert stats for tokenFrom %s and usd value %v", event.TokenFrom, event.UsdValue)
	}

	err = s.repo.UpsertStats(ctx, event.TokenTo, event.UsdValue)
	if err != nil {
		return errors.Wrapf(err, "failed to upsert stats for tokenTo %s and usd value %v", event.TokenTo, event.UsdValue)
	}

	tokenPair := utils.BuildTokenPairkey(event.TokenFrom, event.TokenTo)
	err = s.repo.UpsertStats(ctx, tokenPair, event.UsdValue)
	if err != nil {
		return errors.Wrapf(err, "failed to upsert stats for token pair %s and usd value %v", tokenPair, event.UsdValue)
	}

	return nil
}
