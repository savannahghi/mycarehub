package communities_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	mockMatrix "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/matrix/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities"
)

func TestUseCasesCommunitiesImpl_CreateCommunity(t *testing.T) {
	genderMale := "male"
	clientType := "PMTCT"
	type args struct {
		ctx            context.Context
		communityInput *dto.CommunityInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create matrix room",
			args: args{
				ctx: context.Background(),
				communityInput: &dto.CommunityInput{
					Name:  "Test",
					Topic: "Test",
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 0,
						UpperBound: 0,
					},
					Gender:         []*enumutils.Gender{(*enumutils.Gender)(&genderMale)},
					Preset:         enums.PresetPublicChat,
					Visibility:     enums.PublicVisibility,
					ClientType:     []*enums.ClientType{(*enums.ClientType)(&clientType)},
					OrganisationID: uuid.NewString(),
					ProgramID:      uuid.NewString(),
					FacilityID:     uuid.NewString(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get logged in user",
			args: args{
				ctx: context.Background(),
				communityInput: &dto.CommunityInput{
					Name:  "Test",
					Topic: "Test",
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 0,
						UpperBound: 0,
					},
					Gender:         []*enumutils.Gender{},
					ClientType:     []*enums.ClientType{},
					OrganisationID: uuid.NewString(),
					ProgramID:      uuid.NewString(),
					FacilityID:     uuid.NewString(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create matrix room",
			args: args{
				ctx: context.Background(),
				communityInput: &dto.CommunityInput{
					Name:  "Test",
					Topic: "Test",
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 0,
						UpperBound: 0,
					},
					Gender:         []*enumutils.Gender{},
					ClientType:     []*enums.ClientType{},
					OrganisationID: uuid.NewString(),
					ProgramID:      uuid.NewString(),
					FacilityID:     uuid.NewString(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile by user id",
			args: args{
				ctx: context.Background(),
				communityInput: &dto.CommunityInput{
					Name:  "Test",
					Topic: "Test",
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 0,
						UpperBound: 0,
					},
					Gender:         []*enumutils.Gender{},
					ClientType:     []*enums.ClientType{},
					OrganisationID: uuid.NewString(),
					ProgramID:      uuid.NewString(),
					FacilityID:     uuid.NewString(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to create room in db",
			args: args{
				ctx: context.Background(),
				communityInput: &dto.CommunityInput{
					Name:  "Test",
					Topic: "Test",
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 0,
						UpperBound: 0,
					},
					Gender:         []*enumutils.Gender{},
					ClientType:     []*enums.ClientType{},
					OrganisationID: uuid.NewString(),
					ProgramID:      uuid.NewString(),
					FacilityID:     uuid.NewString(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExt := extensionMock.NewFakeExtension()
			fakeMatrix := mockMatrix.NewMatrixMock()
			uc := communities.NewUseCaseCommunitiesImpl(fakeDB, fakeDB, fakeExt, fakeMatrix)

			if tt.name == "Sad case: unable to get logged in user" {
				fakeExt.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to create matrix room" {
				fakeMatrix.MockCreateCommunity = func(ctx context.Context, auth *domain.MatrixAuth, room *dto.CommunityInput) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get user profile by user id" {
				fakeExt.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.NewString(), nil
				}

				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to create room in db" {
				fakeDB.MockCreateCommunityFn = func(ctx context.Context, community *domain.Community) (*domain.Community, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := uc.CreateCommunity(tt.args.ctx, tt.args.communityInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.CreateCommunity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
