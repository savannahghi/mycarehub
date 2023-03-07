package mock

import (
	"context"
	"net/http"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/matrix"
)

// MatrixMock mocks the matrix's service
type MatrixMock struct {
	MockMakeRequestFn   func(ctx context.Context, payload matrix.RequestHelperPayload) (*http.Response, error)
	MockCreateCommunity func(ctx context.Context, room *dto.CommunityInput) (string, error)
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
		MockCreateCommunity: func(ctx context.Context, room *dto.CommunityInput) (string, error) {
			return gofakeit.BeerName(), nil
		},
	}
}

// MakeRequest mocks the making of http request to Matrix
func (m *MatrixMock) MakeRequest(ctx context.Context, payload matrix.RequestHelperPayload) (*http.Response, error) {
	return m.MockMakeRequestFn(ctx, payload)
}

// CreateCommunity mocks the creation of a Matrix's room
func (m *MatrixMock) CreateCommunity(ctx context.Context, room *dto.CommunityInput) (string, error) {
	return m.MockCreateCommunity(ctx, room)
}
