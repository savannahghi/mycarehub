package communities

// UseCasesCommunities holds all interfaces required to implement the communities feature
type UseCasesCommunities interface {
}

// UseCasesCommunitiesImpl represents communities implementation
type UseCasesCommunitiesImpl struct {
}

// NewUseCaseCommunitiesImpl initializes a new communities service
func NewUseCaseCommunitiesImpl() *UseCasesCommunitiesImpl {
	return &UseCasesCommunitiesImpl{}
}
