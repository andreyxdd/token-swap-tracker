package services

import (
	"consumer/internal/models"
	"consumer/internal/repositories"
	"consumer/internal/utils"
	"context"
	"encoding/json"

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

func (s *StatsService) ProcessSwapEvent(
	ctx context.Context,
	event models.SwapEvent,
	broadcast chan []byte,
) error {
	data, err := s.repo.UpsertStats(ctx, event.TokenFrom, event.UsdValue)
	if err != nil {
		return errors.Wrapf(err, "failed to upsert stats for tokenFrom %s and usd value %v", event.TokenFrom, event.UsdValue)
	}
	statsTokenFrom, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal stats data for tokenFrom %s and usd value %v", event.TokenFrom, event.UsdValue)
	}
	broadcast <- statsTokenFrom

	data, err = s.repo.UpsertStats(ctx, event.TokenTo, event.UsdValue)
	if err != nil {
		return errors.Wrapf(err, "failed to upsert stats for tokenTo %s and usd value %v", event.TokenTo, event.UsdValue)
	}
	statsTokenTo, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal stats data for tokenTo %s and usd value %v", event.TokenTo, event.UsdValue)
	}
	broadcast <- statsTokenTo

	tokenPair := utils.BuildHyphenKey(event.TokenFrom, event.TokenTo)
	data, err = s.repo.UpsertStats(ctx, tokenPair, event.UsdValue)
	if err != nil {
		return errors.Wrapf(err, "failed to upsert stats for token pair %s and usd value %v", tokenPair, event.UsdValue)
	}
	statsTokenPair, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal stats data for token pair %s and usd value %v", tokenPair, event.UsdValue)
	}
	broadcast <- statsTokenPair

	return nil
}
