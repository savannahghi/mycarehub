package infrastructure

import libDatabase "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database"

// Infrastructure ...
type Infrastructure interface {
	libDatabase.Repository
}

//Interactor is used to combine both our internal and open source
// infrastructure implementation
type Interactor struct {
	*libDatabase.DbService
}

// NewInfrastructureInteractor initializes new combined interactor
func NewInfrastructureInteractor() (*Interactor, error) {
	db := libDatabase.NewDbService()
	impl := &Interactor{
		db,
	}

	return impl, nil
}
