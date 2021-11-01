package gorm

import (
	"math"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"gorm.io/gorm"
)

// paginate is a helper function that helps with querying paginated results
func paginate(value interface{}, pagination *domain.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var count int64
	db.Model(value).Count(&count)

	pagination.Count = count
	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	currentPage := pagination.GetPage()

	// If no limit is specified, default to 10
	if pagination.Limit == 0 {
		pagination.Limit = pagination.GetLimit()
	}

	nextPage := currentPage + 1
	pagination.NextPage = &nextPage

	// if we are at the last page, reset the next page to nil
	if nextPage > totalPages {
		pagination.NextPage = nil
	}

	previousPage := currentPage - 1
	pagination.PreviousPage = &previousPage

	// reset to nil if there is no previous page to navigate to
	if previousPage == 0 {
		pagination.PreviousPage = nil
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination)
	}
}
