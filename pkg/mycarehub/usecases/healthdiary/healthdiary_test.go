package healthdiary_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/healthdiary"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/healthdiary/mock"
)

func TestUseCasesHealthDiaryImpl_CreateHealthDiaryEntry(t *testing.T) {
	ctx := context.Background()
	note := gofakeit.HipsterSentence(20)
	type args struct {
		ctx           context.Context
		clientID      string
		note          *string
		mood          string
		reportToStaff bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create a new healthdiary entry",
			args: args{
				ctx:           ctx,
				clientID:      uuid.New().String(),
				note:          &note,
				mood:          string(enums.MoodSad),
				reportToStaff: false,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create healthdiary entry for happy mood",
			args: args{
				ctx:           ctx,
				clientID:      uuid.New().String(),
				note:          &note,
				mood:          string(enums.MoodHappy),
				reportToStaff: false,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to create healthdiary entry for very sad mood",
			args: args{
				ctx:           ctx,
				clientID:      uuid.New().String(),
				note:          &note,
				mood:          string(enums.MoodVerySad),
				reportToStaff: false,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Happy Case - Successfully create service request",
			args: args{
				ctx:           ctx,
				clientID:      uuid.New().String(),
				note:          &note,
				mood:          string(enums.MoodVerySad),
				reportToStaff: false,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create service request for very sad mood",
			args: args{
				ctx:           ctx,
				clientID:      uuid.New().String(),
				note:          &note,
				mood:          string(enums.MoodVerySad),
				reportToStaff: false,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewHealthDiaryUseCaseMock()

			if tt.name == "Sad Case - Fail to create healthdiary entry for happy mood" {
				fakeDB.MockCreateHealthDiaryEntryFn = func(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) error {
					return fmt.Errorf("failed to create health diary entry")
				}
			}

			if tt.name == "Sad Case - Fail to create healthdiary entry for very sad mood" {
				fakeDB.MockCreateHealthDiaryEntryFn = func(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) error {
					return fmt.Errorf("failed to create health diary entry")
				}
			}

			if tt.name == "Sad Case - Fail to create service request for very sad mood" {
				fakeDB.MockCreateServiceRequestFn = func(ctx context.Context, serviceRequestInput *domain.ClientServiceRequest) error {
					return fmt.Errorf("failed to create service request")
				}
			}

			h := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB, fakeDB)
			got, err := h.CreateHealthDiaryEntry(tt.args.ctx, tt.args.clientID, tt.args.note, tt.args.mood, tt.args.reportToStaff)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesHealthDiaryImpl.CreateHealthDiaryEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesHealthDiaryImpl.CreateHealthDiaryEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesHealthDiaryImpl_GetClientHealthDiaryQuote(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ClientHealthDiaryQuote
		wantErr bool
	}{
		{
			name: "Happy Case - successfully get client health diary quote",
			args: args{
				ctx: ctx,
			},
			want: &domain.ClientHealthDiaryQuote{
				Quote:  "test",
				Author: "test",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get quote",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewHealthDiaryUseCaseMock()
			h := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB, fakeDB)

			if tt.name == "Sad Case - Fail to get quote" {
				fakeDB.MockGetClientHealthDiaryQuoteFn = func(ctx context.Context) (*domain.ClientHealthDiaryQuote, error) {
					return nil, fmt.Errorf("failed to get quote")
				}
			}
			got, err := h.GetClientHealthDiaryQuote(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesHealthDiaryImpl.GetClientHealthDiaryQuote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UseCasesHealthDiaryImpl.GetClientHealthDiaryQuote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseHealthDiaryImpl_CanRecordHeathDiary(t *testing.T) {
	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: can create health diary",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing user ID",
			args: args{
				ctx: context.Background(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			healthdiary := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB, fakeDB)

			got, err := healthdiary.CanRecordHeathDiary(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseHealthDiaryImpl.CanRecordHeathDiary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UsecaseHealthDiaryImpl.CanRecordHeathDiary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesHealthDiaryImpl_GetClientHealthDiaryEntries(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get all entries",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing user ID",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeHealthDiary := mock.NewHealthDiaryUseCaseMock()
			healthdiary := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB, fakeDB)

			if tt.name == "Sad Case - Missing user ID" {
				fakeHealthDiary.MockGetClientHealthDiaryEntriesFn = func(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("failed to get client health diary entries")
				}
			}

			got, err := healthdiary.GetClientHealthDiaryEntries(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesHealthDiaryImpl.GetClientHealthDiaryEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}
