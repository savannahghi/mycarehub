//go:build wireinject
// +build wireinject

package wire

import (
	"context"

	"github.com/google/wire"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
)

// InitializeUseCases is an injector that initializes the use cases
func InitializeUseCases(ctx context.Context) (*usecases.MyCareHub, error) {
	wire.Build(WireSet)
	return &usecases.MyCareHub{}, nil
}
