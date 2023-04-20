package mock

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

type SCSSessionManagerMock struct {
	MockPutFn      func(ctx context.Context, key string, val interface{})
	MockDestroyFn  func(ctx context.Context) error
	MockExistsFn   func(ctx context.Context, key string) bool
	MockGetBytesFn func(ctx context.Context, key string) []byte
}

func NewSCSSessionManagerMock() *SCSSessionManagerMock {
	UUID := gofakeit.UUID()
	return &SCSSessionManagerMock{
		MockPutFn: func(ctx context.Context, key string, val interface{}) {},
		MockDestroyFn: func(ctx context.Context) error {
			return nil
		},
		MockExistsFn: func(ctx context.Context, key string) bool {
			return true
		},
		MockGetBytesFn: func(ctx context.Context, key string) []byte {

			type AuthorizationSession struct {
				Page        string
				User        domain.User
				Program     domain.Program
				Facility    domain.Facility
				QueryParams url.Values
			}

			session := &AuthorizationSession{
				Page: "login",
				User: domain.User{
					ID: &UUID,
				},
				Program: domain.Program{
					ID: UUID,
				},
				Facility: domain.Facility{
					ID: &UUID,
				},
				QueryParams: map[string][]string{},
			}

			bs, _ := json.Marshal(session)
			return bs
		},
	}
}

// Put mocks the implementation of Put method
func (m *SCSSessionManagerMock) Put(ctx context.Context, key string, val interface{}) {
	m.MockPutFn(ctx, key, val)
}

// Destroy mocks the implementation of Destroy method
func (m *SCSSessionManagerMock) Destroy(ctx context.Context) error {
	return m.MockDestroyFn(ctx)
}

// Exists mocks the implementation of Exists method
func (m *SCSSessionManagerMock) Exists(ctx context.Context, key string) bool {
	return m.MockExistsFn(ctx, key)
}

// GetBytes mocks the implementation of GetBytes method
func (m *SCSSessionManagerMock) GetBytes(ctx context.Context, key string) []byte {
	return m.MockGetBytesFn(ctx, key)
}
