package healthdiary_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/healthdiary"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/healthdiary/mock"
	serviceRequestMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest/mock"
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
				mood:          enums.MoodSad.String(),
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
				mood:          enums.MoodHappy.String(),
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
				mood:          enums.MoodVerySad.String(),
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
				mood:          enums.MoodVerySad.String(),
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
				mood:          enums.MoodVerySad.String(),
				reportToStaff: false,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Failed to get client profile by client id",
			args: args{
				ctx:           ctx,
				clientID:      uuid.New().String(),
				note:          &note,
				mood:          enums.MoodHappy.String(),
				reportToStaff: false,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()
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
				fakeServiceRequest.MockCreateServiceRequestFn = func(ctx context.Context, input *dto.ServiceRequestInput) (bool, error) {
					return false, fmt.Errorf("failed to create service request")
				}
			}

			if tt.name == "Sad Case - Failed to get client profile by client id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			h := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB, fakeDB, fakeDB, fakeServiceRequest)
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
		ctx   context.Context
		limit int
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ClientHealthDiaryQuote
		wantErr bool
	}{
		{
			name: "Happy Case - successfully get client health diary quote",
			args: args{
				ctx:   ctx,
				limit: 10,
			},
			want: []*domain.ClientHealthDiaryQuote{
				{
					Quote:  "test",
					Author: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get quote",
			args: args{
				ctx:   ctx,
				limit: 10,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewHealthDiaryUseCaseMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()
			h := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB, fakeDB, fakeDB, fakeServiceRequest)

			if tt.name == "Sad Case - Fail to get quote" {
				fakeDB.MockGetClientHealthDiaryQuoteFn = func(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error) {
					return nil, fmt.Errorf("failed to get quote")
				}
			}
			got, err := h.GetClientHealthDiaryQuote(tt.args.ctx, tt.args.limit)
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
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()
			h := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB, fakeDB, fakeDB, fakeServiceRequest)

			got, err := h.CanRecordHeathDiary(tt.args.ctx, tt.args.clientID)
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
		moodType enums.Mood
		shared   bool
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
				moodType: enums.MoodSad,
				shared:   true,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Missing client ID",
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
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()
			h := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB, fakeDB, fakeDB, fakeServiceRequest)

			if tt.name == "Sad Case - Missing user ID" {
				fakeHealthDiary.MockGetClientHealthDiaryEntriesFn = func(ctx context.Context, clientID string, moodType *enums.Mood, shared *bool) ([]*domain.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("failed to get client health diary entries")
				}
			}

			got, err := h.GetClientHealthDiaryEntries(tt.args.ctx, tt.args.clientID, &tt.args.moodType, &tt.args.shared)
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

func TestUseCasesHealthDiaryImpl_GetFacilityHealthDiaryEntries(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.FetchHealthDiaryEntries
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Get facility health diary entries",
			args: args{
				ctx: context.Background(),
				input: dto.FetchHealthDiaryEntries{
					MFLCode:      1234,
					LastSyncTime: &time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Failed to check if facility exists",
			args: args{
				ctx: context.Background(),
				input: dto.FetchHealthDiaryEntries{
					MFLCode: 1212,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Non-existent facility",
			args: args{
				ctx: context.Background(),
				input: dto.FetchHealthDiaryEntries{
					MFLCode: 9932,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to retrieve facility by MFLCODE",
			args: args{
				ctx: context.Background(),
				input: dto.FetchHealthDiaryEntries{
					MFLCode: 1234,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get clients in a facility",
			args: args{
				ctx: context.Background(),
				input: dto.FetchHealthDiaryEntries{
					MFLCode: 12345,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get recent health diary entries",
			args: args{
				ctx: context.Background(),
				input: dto.FetchHealthDiaryEntries{
					MFLCode:      1234,
					LastSyncTime: &time.Time{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()
			h := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB, fakeDB, fakeDB, fakeServiceRequest)

			if tt.name == "Sad Case - Failed to check if facility exists" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, fmt.Errorf("failed to check if facility exists")
				}
			}

			if tt.name == "Sad Case - Non-existent facility" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "Sad Case - Fail to retrieve facility by MFLCODE" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("failed to retrieve facility")
				}
			}

			if tt.name == "Sad Case - Fail to get clients in a facility" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockGetClientsInAFacilityFn = func(ctx context.Context, facilityID string) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get clients within a facility")
				}
			}

			if tt.name == "Sad Case - Fail to get recent health diary entries" {
				fakeDB.MockGetRecentHealthDiaryEntriesFn = func(ctx context.Context, lastSyncTime time.Time, client *domain.ClientProfile) ([]*domain.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("failed to get recent health diary entries")
				}
			}

			got, err := h.GetFacilityHealthDiaryEntries(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesHealthDiaryImpl.GetFacilityHealthDiaryEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesHealthDiaryImpl_ShareHealthDiaryEntry(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()
	fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()
	h := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB, fakeDB, fakeDB, fakeServiceRequest)

	type args struct {
		ctx                    context.Context
		healthDiaryEntryID     string
		shareEntireHealthDiary bool
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
				ctx:                    ctx,
				healthDiaryEntryID:     uuid.New().String(),
				shareEntireHealthDiary: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - unable to create service request",
			args: args{
				ctx:                ctx,
				healthDiaryEntryID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - unable to get health diary by id",
			args: args{
				ctx:                ctx,
				healthDiaryEntryID: "",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - unable to update health diary",
			args: args{
				ctx:                ctx,
				healthDiaryEntryID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - unable to create service request" {
				fakeServiceRequest.MockCreateServiceRequestFn = func(ctx context.Context, input *dto.ServiceRequestInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - unable to get health diary by id" {
				fakeDB.MockGetHealthDiaryEntryByIDFn = func(ctx context.Context, healthDiaryEntryID string) (*domain.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - unable to update health diary" {
				fakeDB.MockGetHealthDiaryEntryByIDFn = func(ctx context.Context, healthDiaryEntryID string) (*domain.ClientHealthDiaryEntry, error) {
					ID := uuid.New().String()
					return &domain.ClientHealthDiaryEntry{
						ID: &ID,
					}, nil
				}
				fakeDB.MockUpdateHealthDiaryFn = func(ctx context.Context, clientHealthDiaryEntry *domain.ClientHealthDiaryEntry, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			got, err := h.ShareHealthDiaryEntry(tt.args.ctx, tt.args.healthDiaryEntryID, tt.args.shareEntireHealthDiary)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesHealthDiaryImpl.ShareHealthDiaryEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesHealthDiaryImpl.ShareHealthDiaryEntry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesHealthDiaryImpl_GetSharedHealthDiaryEntry(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()
	fakeServiceRequest := serviceRequestMock.NewServiceRequestUseCaseMock()
	h := healthdiary.NewUseCaseHealthDiaryImpl(fakeDB, fakeDB, fakeDB, fakeServiceRequest)

	type args struct {
		ctx        context.Context
		clientID   string
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ClientHealthDiaryEntry
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:        ctx,
				clientID:   uuid.New().String(),
				facilityID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:        ctx,
				clientID:   uuid.New().String(),
				facilityID: "",
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty client ID",
			args: args{
				ctx:        ctx,
				clientID:   uuid.New().String(),
				facilityID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeDB.MockGetSharedHealthDiaryEntriesFn = func(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty client ID" {
				fakeDB.MockGetSharedHealthDiaryEntriesFn = func(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := h.GetSharedHealthDiaryEntries(tt.args.ctx, tt.args.clientID, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesHealthDiaryImpl.GetSharedHealthDiaryEntries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected shared health diary entries to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected shared health diary entries not to be nil for %v", tt.name)
				return
			}
		})
	}
}
