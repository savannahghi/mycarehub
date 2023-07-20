package gorm

import (
	"context"
	"fmt"

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

// ProgramScope generates a GORM scope to filter records based on the program ID.
//
// This function creates a GORM scope that filters records in a table
// based on the program ID stored in the provided context. It allows specifying the table
// name in the 'tableName' parameter to avoid ambiguity when the 'program_id' column
// exists in multiple tables involved in the query.
//
// When ambiguity occurs, an error like `pq: column reference "program_id" is ambiguous` will be returned.
// The error occurs in GORM when a column name is used in a query, but the database engine cannot determine which
// table the column belongs to due to the column being present in multiple tables involved in the query
func ProgramScope(ctx context.Context, tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		id, err := utils.GetValueFromContext(ctx, utils.ProgramContextKey)
		if err != nil {
			return db
		}
		columnReference := fmt.Sprintf("%s.program_id", tableName)
		return db.Where(columnReference, id)
	}
}
