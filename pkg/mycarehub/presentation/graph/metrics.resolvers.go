package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) CollectMetric(ctx context.Context, input domain.Metric) (bool, error) {
	return r.mycarehub.Metrics.CollectMetric(ctx, &input)
}
