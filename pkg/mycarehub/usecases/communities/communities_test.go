package communities_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
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

func TestUseCasesCommunitiesImpl_ListCommunities(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: list communities",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get logged in user",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to list communities",
			args: args{
				ctx: context.Background(),
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
					return "", errors.New("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, errors.New("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to list communities" {
				fakeDB.MockListCommunitiesFn = func(ctx context.Context, programID, organisationID string) ([]*domain.Community, error) {
					return nil, errors.New("unable to list communities")
				}
			}
			_, err := uc.ListCommunities(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.ListCommunities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_SearchUsers(t *testing.T) {
	limit := 10
	type args struct {
		ctx        context.Context
		limit      *int
		searchTerm string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: successfully search users",
			args: args{
				ctx:        context.Background(),
				limit:      &limit,
				searchTerm: "test",
			},
			wantErr: false,
		},
		{
			name: "Sad case: search term less than 3 character",
			args: args{
				ctx:        context.Background(),
				limit:      &limit,
				searchTerm: "te",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get logged in user",
			args: args{
				ctx:        context.Background(),
				limit:      &limit,
				searchTerm: "te",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get logged in user profile by id",
			args: args{
				ctx:        context.Background(),
				limit:      &limit,
				searchTerm: "test",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to search for matrix user",
			args: args{
				ctx:        context.Background(),
				limit:      &limit,
				searchTerm: "test",
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile by username",
			args: args{
				ctx:        context.Background(),
				limit:      &limit,
				searchTerm: "test",
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
					return "", errors.New("unable to get logged in user")
				}
			}
			if tt.name == "Sad case: unable to get logged in user profile by id" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, errors.New("unable to get logged in user profile by id")
				}
			}
			if tt.name == "Sad case: unable to search for matrix user" {
				fakeMatrix.MockSearchUsersFn = func(ctx context.Context, limit int, searchTerm string, auth *domain.MatrixAuth) (*domain.MatrixUserSearchResult, error) {
					return nil, errors.New("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to get user profile by username" {
				fakeDB.MockGetUserProfileByUsernameFn = func(ctx context.Context, username string) (*domain.User, error) {
					return nil, errors.New("unable to get user profile by username")
				}
			}
			_, err := uc.SearchUsers(tt.args.ctx, tt.args.limit, tt.args.searchTerm)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.SearchUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCasesCommunitiesImpl_SetPusher(t *testing.T) {
	type args struct {
		ctx     context.Context
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: successfully set a pusher - pro",
			args: args{
				ctx:     context.Background(),
				flavour: feedlib.FlavourPro,
			},
			wantErr: false,
		},
		{
			name: "Happy case: successfully set a pusher - consumer",
			args: args{
				ctx:     context.Background(),
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get logged in user",
			args: args{
				ctx:     context.Background(),
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile by user id",
			args: args{
				ctx:     context.Background(),
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to set pusher",
			args: args{
				ctx:     context.Background(),
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid flavor",
			args: args{
				ctx:     context.Background(),
				flavour: "invalid",
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
					return "", fmt.Errorf("unable to get logged in user")
				}
			}
			if tt.name == "Sad case: unable to get user profile by user id" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("unable to get user profile")
				}
			}
			if tt.name == "Sad case: unable to set pusher" {
				fakeMatrix.MockSetPusherFn = func(ctx context.Context, auth *domain.MatrixAuth, payload *domain.PusherPayload) error {
					return fmt.Errorf("unable to set pusher")
				}
			}

			_, err := uc.SetPusher(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.SetPusher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
