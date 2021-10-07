package infrastructure

// Interactor is an implementation of the infrastructure interface
// It combines each individual service implementation
type Interactor struct {
}

// NewInteractor initializes a new infrastructure interactor
func NewInteractor() Interactor {

	return Interactor{}
}
