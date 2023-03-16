package terms_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms/mock"
	"github.com/segmentio/ksuid"
)

func TestTermsOfServiceImpl_GetCurrentTerms_Unittest(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "Sad case - nil context",
			args: args{
				ctx: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewTermsUseCaseMock()

			j := terms.NewUseCasesTermsOfService(fakeDB, fakeDB, fakeDB)

			if tt.name == "Sad case - nil context" {
				fakeDB.MockGetCurrentTermsFn = func(ctx context.Context) (*domain.TermsOfService, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := j.GetCurrentTerms(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TermsOfServiceImpl.GetCurrentTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestServiceTermsImpl_AcceptTerms(t *testing.T) {
	ctx := context.Background()

	userID := ksuid.New().String()
	termsID := gofakeit.Number(1, 100000)
	negativeTermsID := gofakeit.Number(-100000, -1)

	type args struct {
		ctx     context.Context
		userID  *string
		termsID *int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				userID:  &userID,
				termsID: &termsID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				userID:  &userID,
				termsID: &termsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no termsID",
			args: args{
				ctx:     ctx,
				userID:  &userID,
				termsID: nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:     ctx,
				userID:  nil,
				termsID: &termsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and termsID",
			args: args{
				ctx:     ctx,
				userID:  nil,
				termsID: nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - negative termsID",
			args: args{
				ctx:     ctx,
				userID:  nil,
				termsID: &negativeTermsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - userID and negative termsID",
			args: args{
				ctx:     ctx,
				userID:  &userID,
				termsID: &negativeTermsID,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewTermsUseCaseMock()

			j := terms.NewUseCasesTermsOfService(fakeDB, fakeDB, fakeDB)

			if tt.name == "Sad case" {
				fakeDB.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no termsID" {
				fakeDB.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID" {
				fakeDB.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID and termsID" {
				fakeDB.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - negative termsID" {
				fakeDB.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - userID and negative termsID" {
				fakeDB.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := j.AcceptTerms(tt.args.ctx, tt.args.userID, tt.args.termsID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceTermsImpl.AcceptTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceTermsImpl.AcceptTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceTermsImpl_CreateTermsOfService(t *testing.T) {
	dummyText := gofakeit.BS()
	type args struct {
		ctx            context.Context
		termsOfService *domain.TermsOfService
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create terms of service",
			args: args{
				ctx: context.Background(),
				termsOfService: &domain.TermsOfService{
					TermsID:   1,
					Text:      &dummyText,
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: failed to create security questions",
			args: args{
				ctx: context.Background(),
				termsOfService: &domain.TermsOfService{
					TermsID:   1,
					Text:      &dummyText,
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()

			if tt.name == "Sad Case: failed to create security questions" {
				fakeDB.MockCreateTermsOfServiceFn = func(ctx context.Context, termsOfService *domain.TermsOfService) (*domain.TermsOfService, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			tr := terms.NewUseCasesTermsOfService(fakeDB, fakeDB, fakeDB)
			_, err := tr.CreateTermsOfService(tt.args.ctx, tt.args.termsOfService)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceTermsImpl.CreateTermsOfService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
