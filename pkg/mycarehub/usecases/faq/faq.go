package faq

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IGetFAQContent gets the faq content
type IGetFAQContent interface {
	GetFAQContent(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*domain.FAQ, error)
}

// UsecaseFAQ groups al the interfaces for the FAQ usecase
type UsecaseFAQ interface {
	IGetFAQContent
}

// UsecaseFAQImpl represents the FAQ implementation
type UsecaseFAQImpl struct {
	Query infrastructure.Query
}

// NewUsecaseFAQ is the controller function for the FAQ usecase
func NewUsecaseFAQ(
	query infrastructure.Query,
) *UsecaseFAQImpl {
	return &UsecaseFAQImpl{
		Query: query,
	}
}

// GetFAQContent gets FAQ content for a given flavour
// an optional limit can be provided
func (u *UsecaseFAQImpl) GetFAQContent(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*domain.FAQ, error) {
	defaultLimit := 10
	if !flavour.IsValid() {
		return nil, exceptions.InvalidFlavourDefinedErr(fmt.Errorf("flavour is not valid"))
	}
	if limit == nil {

		limit = &defaultLimit
	}

	faqs, err := u.Query.GetFAQContent(ctx, flavour, limit)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.GetFAQContentErr(fmt.Errorf("error getting faq content: %w", err))
	}
	return faqs, nil
}
