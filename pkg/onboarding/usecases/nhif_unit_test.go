package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/profileutils"
)

func TestNHIFUseCaseImpl_AddNHIFDetails(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	validInput := dto.NHIFDetailsInput{
		MembershipNumber:          "123456",
		Employment:                domain.EmploymentTypeEmployed,
		NHIFCardPhotoID:           uuid.New().String(),
		IDDocType:                 enumutils.IdentificationDocTypeMilitary,
		IdentificationCardPhotoID: uuid.New().String(),
		IDNumber:                  "11111111",
	}
	type args struct {
		ctx   context.Context
		input dto.NHIFDetailsInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.NHIFDetails
		wantErr bool
	}{
		{
			name: "happy:) successfully add NHIF Details",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: false,
		},
		{
			name: "sad:( fail to add NHIF Details",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "sad:( fail to get a user profile",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "sad:( fail to get logged in user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "sad:( fail to resolve add NHIF nudge",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy:) successfully add NHIF Details" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: uuid.New().String()}, nil
				}

				fakeRepo.AddNHIFDetailsFn = func(ctx context.Context, input dto.NHIFDetailsInput, profileID string) (*domain.NHIFDetails, error) {
					return &domain.NHIFDetails{
						ID:               uuid.New().String(),
						ProfileID:        profileID,
						MembershipNumber: "12345",
						IDNumber:         "12345",
					}, nil
				}

				fakeEngagementSvs.ResolveDefaultNudgeByTitleFn = func(
					ctx context.Context,
					UID string,
					flavour feedlib.Flavour,
					nudgeTitle string,
				) error {
					return nil
				}
			}

			if tt.name == "sad:( fail to add NHIF Details" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: uuid.New().String()}, nil
				}

				fakeRepo.AddNHIFDetailsFn = func(ctx context.Context, input dto.NHIFDetailsInput, profileID string) (*domain.NHIFDetails, error) {
					return nil, fmt.Errorf("failed to add nhif details")
				}
			}

			if tt.name == "sad:( fail to get a user profile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "sad:( fail to get logged in user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "sad:( fail to resolve add NHIF nudge" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{ID: uuid.New().String()}, nil
				}

				fakeRepo.AddNHIFDetailsFn = func(ctx context.Context, input dto.NHIFDetailsInput, profileID string) (*domain.NHIFDetails, error) {
					return &domain.NHIFDetails{
						ID:               uuid.New().String(),
						ProfileID:        profileID,
						MembershipNumber: "12345",
						IDNumber:         "12345",
					}, nil
				}

				fakeEngagementSvs.ResolveDefaultNudgeByTitleFn = func(
					ctx context.Context,
					UID string,
					flavour feedlib.Flavour,
					nudgeTitle string,
				) error {
					return fmt.Errorf("an error occurred")
				}
			}
			got, err := i.NHIF.AddNHIFDetails(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NHIFUseCaseImpl.AddNHIFDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if got == nil {
					t.Errorf("nil response returned")
					return
				}
			}
		})
	}
}

func TestNHIFUseCaseImpl_NHIFDetails(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboardingInteractor()
	if err != nil {
		t.Errorf("failed to initialize fake onboarding interactor")
		return
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.NHIFDetails
		wantErr bool
	}{
		{
			name: "happy:) successfully return NHIF details",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "happy:) successfully return nil NHIF details",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "happy:) fail to return NHIF Details",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "sad:( fail to get user profile",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "sad:( fail to get logged in user",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "happy:) successfully return NHIF details" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: uuid.New().String(),
						VerifiedIdentifiers: []profileutils.VerifiedIdentifier{
							{
								UID: uid,
							},
						},
					}, nil
				}

				fakeRepo.GetNHIFDetailsByProfileIDFn = func(ctx context.Context, profileID string) (*domain.NHIFDetails, error) {
					return &domain.NHIFDetails{
						ID:               uuid.New().String(),
						ProfileID:        profileID,
						MembershipNumber: "12345",
						IDNumber:         "12345",
					}, nil
				}
			}

			if tt.name == "happy:) successfully return nil NHIF details" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: uuid.New().String(),
						VerifiedIdentifiers: []profileutils.VerifiedIdentifier{
							{
								UID: uid,
							},
						},
					}, nil
				}

				fakeRepo.GetNHIFDetailsByProfileIDFn = func(ctx context.Context, profileID string) (*domain.NHIFDetails, error) {
					return nil, nil
				}
			}

			if tt.name == "happy:) fail to return NHIF Details" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID: uuid.New().String(),
						VerifiedIdentifiers: []profileutils.VerifiedIdentifier{
							{
								UID: uid,
							},
						},
					}, nil
				}

				fakeRepo.GetNHIFDetailsByProfileIDFn = func(ctx context.Context, profileID string) (*domain.NHIFDetails, error) {
					return nil, fmt.Errorf("failed to get the user's nhif details")
				}
			}

			if tt.name == "sad:( fail to get user profile" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return &dto.UserInfo{
						UID:         "12233",
						Email:       "test@example.com",
						PhoneNumber: "0721568526",
					}, nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "sad:( fail to get logged in user" {
				fakeBaseExt.GetLoggedInUserFn = func(ctx context.Context) (*dto.UserInfo, error) {
					return nil, fmt.Errorf("failed to get logged in user")
				}
			}

			_, err := i.NHIF.NHIFDetails(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NHIFUseCaseImpl.NHIFDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}
