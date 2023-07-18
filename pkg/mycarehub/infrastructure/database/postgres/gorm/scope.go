package gorm

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"gorm.io/gorm"
)

// OrganisationScope is a reusable query used for filtering out records for a specific organisation
func OrganisationScope(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		id, err := utils.GetValueFromContext(ctx, utils.OrganisationContextKey)
		if err != nil {
			return db
		}

		return db.Where("organisation_id", id)
	}
}

// ProgramScope is a reusable query used for filtering out records for a specific program
func ProgramScope(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		id, err := utils.GetValueFromContext(ctx, utils.ProgramContextKey)
		if err != nil {
			return db
		}

		return db.Where("program_id", id)
	}
}
