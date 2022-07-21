package gorm

import (
	"fmt"
	"math"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"gorm.io/gorm"
)

// paginate is a helper function that helps with querying paginated results
func paginate(value interface{}, pagination *domain.Pagination, count int64, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	pagination.Count = count

	// If no limit is specified, default to 10
	if pagination.Limit == 0 {
		pagination.Limit = pagination.GetLimit()
	}
	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	currentPage := pagination.GetPage()

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
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

// parse filter param values to map[string]interface{}
func filterParamsToMap(mapString []*domain.FiltersParam) map[string]interface{} {
	res := make(map[string]interface{})
	for _, v := range mapString {
		res[v.Name] = v.Value
	}
	return res
}

func addFilters(transaction *gorm.DB, filters []*firebasetools.FilterParam) error {
	for _, filter := range filters {
		op, err := firebasetools.OpString(filter.ComparisonOperation)
		if err != nil {
			return err
		}
		// convert firebase equal to postgres equal
		if op == "==" {
			op = "="
		}

		if filter.ComparisonOperation == enumutils.OperationIn {
			transaction.Where(fmt.Sprintf("%s %s ?", filter.FieldName, op), filter.FieldValue)
			continue
		}

		switch filter.FieldType {
		case enumutils.FieldTypeBoolean:
			value, ok := filter.FieldValue.(bool)
			if !ok {
				return fmt.Errorf("expected filter value to be true or false")
			}
			transaction.Where(fmt.Sprintf("%s %s ?", filter.FieldName, op), value)

		case enumutils.FieldTypeInteger:
			value, ok := filter.FieldValue.(int)
			if !ok {
				return fmt.Errorf("expected filter value to be an int")
			}
			transaction.Where(fmt.Sprintf("%s %s ?", filter.FieldName, op), value)

		case enumutils.FieldTypeTimestamp:
			value, ok := filter.FieldValue.(string)
			if !ok {
				return fmt.Errorf("expected filter value to be a timestamp")
			}
			transaction.Where(fmt.Sprintf("%s %s ?", filter.FieldName, op), value)

		case enumutils.FieldTypeString:
			value, ok := filter.FieldValue.(string)
			if !ok {
				return fmt.Errorf("expected filter value to be a string")
			}
			transaction.Where(fmt.Sprintf("%s %s ?", filter.FieldName, op), value)
		default:
			return fmt.Errorf("unexpected field type '%s'", filter.FieldType.String())
		}

	}

	return nil
}

// paginateQuery adds offset and limit conditions to a database query
func paginateQuery(db *gorm.DB, pagination *domain.Pagination) {
	// If no limit is specified, default to 10
	if pagination.Limit == 0 {
		pagination.Limit = pagination.GetLimit()
	}
	totalPages := int(math.Ceil(float64(pagination.Count) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	currentPage := pagination.GetPage()

	nextPage := currentPage + 1
	if nextPage <= totalPages {
		pagination.NextPage = &nextPage
	}

	previousPage := currentPage - 1
	if previousPage > 0 {
		pagination.PreviousPage = &previousPage
	}

	db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())

}
