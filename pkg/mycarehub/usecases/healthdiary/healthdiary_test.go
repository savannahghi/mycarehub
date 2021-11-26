package healthdiary_test

import (
	"context"
	"fmt"
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

			if tt.name == "Sad Case - Fail to create service request for very sad mood" {
				fakeDB.MockCreateServiceRequestFn = func(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry, serviceRequestInput *domain.ClientServiceRequest) error {
					return fmt.Errorf("failed to create service request")
				}
			}

			h := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB)
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
