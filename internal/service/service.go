package service

import (
	"context"
	"github.com/captainlettuce/sales-stats/pkg/types"
)

type OrderAggregationRepository interface {
	AggregateSalesPrice(ctx context.Context, c types.AggregationRequest) ([]types.AggregationResult, error)
}

type OrderAggregationService struct {
	repo OrderAggregationRepository
}

func New(repo OrderAggregationRepository) *OrderAggregationService {
	return &OrderAggregationService{repo: repo}
}

func (s *OrderAggregationService) AggregateSalesPrice(ctx context.Context, c types.AggregationRequest) ([]types.AggregationResult, error) {

	// Set default aggregation method to average
	if len(c.AggregationMethods) == 0 {
		c.AggregationMethods = []types.AggregationMethodType{types.AggregationMethodAvg}
	}

	return s.repo.AggregateSalesPrice(ctx, c)
}
