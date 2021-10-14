package dto

import (
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

// FacilityEdge is used to serialize GraphQL Relay edges for healthcare facilities
type FacilityEdge struct {
	Cursor *string          `json:"cursor"`
	Node   *domain.Facility `json:"node"`
}

// FacilityConnection is used to serialize GraphQL Relay connections for healthcare facilities
type FacilityConnection struct {
	Edges    []*FacilityEdge         `json:"edges"`
	PageInfo *firebasetools.PageInfo `json:"pageInfo"`
}

// PaginationListResponse defines the output of a paginated list
type PaginationListResponse struct {
	Count       int    `json:"count,omitempty"`
	Next        string `json:"next,omitempty"`
	Previous    string `json:"previous,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
	CurrentPage int    `json:"current_page,omitempty"`
	TotalPages  int    `json:"total_pages,omitempty"`
	StartIndex  int    `json:"start_index,omitempty"`
	EndIndex    int    `json:"end_index,omitempty"`
}
