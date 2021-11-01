package domain

// Pagination contains the struct fields for performing pagination.
type Pagination struct {
	Limit        int   `json:"limit"`
	CurrentPage  int   `json:"currentPage"`
	Count        int64 `json:"count"`
	TotalPages   int   `json:"totalPages"`
	NextPage     *int  `json:"nextPage"`
	PreviousPage *int  `json:"previousPage"`
}

// GetOffset calculates the deviation in pages that come before
func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

// GetLimit calculates the maximum number of items to be shown per page
func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

// GetPage gets the current page
func (p *Pagination) GetPage() int {
	if p.CurrentPage == 0 {
		p.CurrentPage = 1
	}
	return p.CurrentPage
}
