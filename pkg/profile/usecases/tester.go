package usecases

import "context"

// TesterUseCases represents all the business logic that touch the users that login to test the app
type TesterUseCases interface {
	AddTester(ctx context.Context, email string) (bool, error)
	RemoveTester(ctx context.Context, email string) (bool, error)
	ListTesters(ctx context.Context) ([]string, error)
}
