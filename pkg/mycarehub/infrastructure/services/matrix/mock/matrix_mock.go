package mock

import (
	"context"
	"net/http"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/matrix"
)

// MatrixMock mocks the matrix's service
type MatrixMock struct {
	MockMakeRequestFn        func(ctx context.Context, payload matrix.RequestHelperPayload) (*http.Response, error)
	MockCreateCommunity      func(ctx context.Context, auth *domain.MatrixAuth, room *dto.CommunityInput) (string, error)
	MockRegisterUserFn       func(ctx context.Context, auth *domain.MatrixAuth, registrationPayload *domain.MatrixUserRegistration) (*dto.MatrixUserRegistrationOutput, error)
	MockLoginFn              func(ctx context.Context, username string, password string) (*domain.CommunityProfile, error)
	MockCheckIfUserIsAdminFn func(ctx context.Context, auth *domain.MatrixAuth, userID string) (bool, error)
	MockSearchUsersFn        func(ctx context.Context, limit int, searchTerm string, auth *domain.MatrixAuth) (*domain.MatrixUserSearchResult, error)
	MockDeactivateUserFn     func(ctx context.Context, userID string, auth *domain.MatrixAuth) error
}

// NewSurveysMock initializes the surveys mock service
func NewMatrixMock() *MatrixMock {
	return &MatrixMock{
		MockMakeRequestFn: func(ctx context.Context, payload matrix.RequestHelperPayload) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       nil,
			}, nil
		},
		MockCreateCommunity: func(ctx context.Context, auth *domain.MatrixAuth, room *dto.CommunityInput) (string, error) {
			return gofakeit.BeerName(), nil
		},
		MockRegisterUserFn: func(ctx context.Context, auth *domain.MatrixAuth, registrationPayload *domain.MatrixUserRegistration) (*dto.MatrixUserRegistrationOutput, error) {
			return &dto.MatrixUserRegistrationOutput{
				UserID: gofakeit.BeerName(),
			}, nil
		},
		MockLoginFn: func(ctx context.Context, username, password string) (*domain.CommunityProfile, error) {
			return &domain.CommunityProfile{
				UserID:      "@test:prohealth360.org",
				AccessToken: "sys_",
				HomeServer:  "prohealth360.org",
				DeviceID:    "OBVSIQJYBO",
				WellKnown: domain.WellKnown{
					MHomeserver: domain.MHomeserver{
						BaseURL: "https://matrix.prohealth360.org/",
					},
				},
			}, nil
		},
		MockCheckIfUserIsAdminFn: func(ctx context.Context, auth *domain.MatrixAuth, userID string) (bool, error) {
			return true, nil
		},
		MockSearchUsersFn: func(ctx context.Context, limit int, searchTerm string, auth *domain.MatrixAuth) (*domain.MatrixUserSearchResult, error) {
			return &domain.MatrixUserSearchResult{
				Limited: false,
				Results: []domain.Result{
					{
						UserID:      "@test:prohealth360.org",
						DisplayName: "test",
						AvatarURL:   "mxc://bar.com/foo",
					},
				},
			}, nil
		},
		MockDeactivateUserFn: func(ctx context.Context, userID string, auth *domain.MatrixAuth) error {
			return nil
		},
	}
}

// MakeRequest mocks the making of http request to Matrix
func (m *MatrixMock) MakeRequest(ctx context.Context, payload matrix.RequestHelperPayload) (*http.Response, error) {
	return m.MockMakeRequestFn(ctx, payload)
}

// CreateCommunity mocks the creation of a Matrix's room
func (m *MatrixMock) CreateCommunity(ctx context.Context, auth *domain.MatrixAuth, room *dto.CommunityInput) (string, error) {
	return m.MockCreateCommunity(ctx, auth, room)
}

// RegisterUser mocks the registration of user in Matrix homeserver
func (m *MatrixMock) RegisterUser(ctx context.Context, auth *domain.MatrixAuth, registrationPayload *domain.MatrixUserRegistration) (*dto.MatrixUserRegistrationOutput, error) {
	return m.MockRegisterUserFn(ctx, auth, registrationPayload)
}

// Login mocks authentication if a matrix user
func (m *MatrixMock) Login(ctx context.Context, username, password string) (*domain.CommunityProfile, error) {
	return m.MockLoginFn(ctx, username, password)
}

// CheckIfUserIsAdmin mocks the checking of whether a user is an admin or not
func (m *MatrixMock) CheckIfUserIsAdmin(ctx context.Context, auth *domain.MatrixAuth, userID string) (bool, error) {
	return m.MockCheckIfUserIsAdminFn(ctx, auth, userID)
}

// SearchUsers mocks the implementation of searching for a Matrix user
func (m *MatrixMock) SearchUsers(ctx context.Context, limit int, searchTerm string, auth *domain.MatrixAuth) (*domain.MatrixUserSearchResult, error) {
	return m.MockSearchUsersFn(ctx, limit, searchTerm, auth)
}

// DeactivateUser mocks the deactivation of a matrix user
func (m *MatrixMock) DeactivateUser(ctx context.Context, userID string, auth *domain.MatrixAuth) error {
	return m.MockDeactivateUserFn(ctx, userID, auth)
}
