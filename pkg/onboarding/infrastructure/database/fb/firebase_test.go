package fb_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database/fb"
	extMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database/fb/mock"
	"github.com/savannahghi/profileutils"
	"github.com/stretchr/testify/assert"
)

var fakeFireBaseClientExt extMock.FirebaseClientExtension
var fireBaseClientExt fb.FirebaseClientExtension = &fakeFireBaseClientExt
var fakeFireStoreClientExt extMock.FirestoreClientExtension

func TestRepository_UpdateUserName(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx      context.Context
		id       string
		userName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:update_user_name_failed_to_get_a user_profile",
			args: args{
				ctx:      ctx,
				id:       "12333",
				userName: "mwas",
			},
			wantErr: true,
		},
		{
			name: "invalid:user_name_already_exists",
			args: args{
				ctx:      ctx,
				id:       "12333",
				userName: "mwas",
			},
			wantErr: true,
		}, {
			name: "valid:user_name_not_found",
			args: args{
				ctx:      ctx,
				id:       "12333",
				userName: "mwas",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:update_user_name_failed_to_get_a user_profile" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}

				fakeFireStoreClientExt.UpdateFn = func(ctx context.Context, command *fb.UpdateCommand) error {
					return nil
				}
			}

			if tt.name == "invalid:user_name_already_exists" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{
						{
							Ref: &firestore.DocumentRef{
								ID: "5555",
							},
						},
					}
					return docs, nil
				}

				fakeFireStoreClientExt.UpdateFn = func(ctx context.Context, command *fb.UpdateCommand) error {
					return nil
				}
			}

			if tt.name == "valid:user_name_not_found" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
				fakeFireStoreClientExt.UpdateFn = func(ctx context.Context, command *fb.UpdateCommand) error {
					return nil
				}
			}
			err := repo.UpdateUserName(tt.args.ctx, tt.args.id, tt.args.userName)

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

func TestRepository_CheckIfExperimentParticipant(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name     string
		args     args
		expected bool
	}{
		{
			name: "valid:exists",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
			},
			expected: true,
		},
		{
			name: "valid:does_not_exist",
			args: args{
				ctx: ctx,
				id:  uuid.New().String(),
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:exists" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{
						{
							Ref: &firestore.DocumentRef{
								ID: uuid.New().String(),
							},
						},
					}
					return docs, nil
				}
			}

			if tt.name == "valid:does_not_exist" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
			}

			exists, err := repo.CheckIfExperimentParticipant(tt.args.ctx, tt.args.id)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, exists)
		})
	}
}

func TestRepository_AddUserAsExperimentParticipant(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx     context.Context
		profile *profileutils.UserProfile
	}
	tests := []struct {
		name     string
		args     args
		expected bool
		wantErr  bool
	}{
		{
			name: "valid:add",
			args: args{
				ctx: ctx,
				profile: &profileutils.UserProfile{
					ID: uuid.New().String(),
				},
			},
			expected: true,
		},
		{
			name: "valid:already_exists",
			args: args{
				ctx: ctx,
				profile: &profileutils.UserProfile{
					ID: uuid.New().String(),
				},
			},
			expected: true,
		},

		{
			name: "invalid:throws_internal_server_error_while_checking_existence",
			args: args{
				ctx: ctx,
				profile: &profileutils.UserProfile{
					ID: uuid.New().String(),
				},
			},
			expected: false,
			wantErr:  true,
		},

		{
			name: "invalid:throws_internal_server_error_while_creating",
			args: args{
				ctx: ctx,
				profile: &profileutils.UserProfile{
					ID: uuid.New().String(),
				},
			},
			expected: false,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:add" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}

				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					doc := firestore.DocumentRef{
						ID: uuid.New().String(),
					}
					return &doc, nil
				}
			}

			if tt.name == "valid:already_exists" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{
						{
							Ref: &firestore.DocumentRef{
								ID: uuid.New().String(),
							},
						},
					}
					return docs, nil
				}

			}

			if tt.name == "invalid:throws_internal_server_error_while_checking_existence" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, exceptions.InternalServerError(
						fmt.Errorf("unable to parse user profile as firebase snapshot"),
					)
				}
			}

			if tt.name == "invalid:throws_internal_server_error_while_creating" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}

				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return nil, exceptions.InternalServerError(
						fmt.Errorf("unable to add user profile of ID in experiment_participant"),
					)
				}
			}

			resp, err := repo.AddUserAsExperimentParticipant(tt.args.ctx, tt.args.profile)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.expected, resp)
		})
	}
}

func TestRepository_RemoveUserAsExperimentParticipant(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx     context.Context
		profile *profileutils.UserProfile
	}
	tests := []struct {
		name     string
		args     args
		expected bool
		wantErr  bool
	}{
		{
			name: "valid:remove_user_as_experiment_participant",
			args: args{
				ctx: ctx,
				profile: &profileutils.UserProfile{
					ID: uuid.New().String(),
				},
			},
			expected: true,
		},

		{
			name: "invalid:throws_internal_server_error_while_removing",
			args: args{
				ctx: ctx,
				profile: &profileutils.UserProfile{
					ID: uuid.New().String(),
				},
			},
			expected: false,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:remove_user_as_experiment_participant" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{
						{
							Ref: &firestore.DocumentRef{
								ID: uuid.New().String(),
							},
						},
					}
					return docs, nil
				}

				fakeFireStoreClientExt.DeleteFn = func(ctx context.Context, command *fb.DeleteCommand) error {
					return nil
				}

			}
			if tt.name == "invalid:throws_internal_server_error_while_removing" {
				fakeFireStoreClientExt.DeleteFn = func(ctx context.Context, command *fb.DeleteCommand) error {
					return exceptions.InternalServerError(
						fmt.Errorf(
							"unable to remove user profile of ID  from experiment_participant",
						),
					)
				}
			}

			resp, err := repo.RemoveUserAsExperimentParticipant(tt.args.ctx, tt.args.profile)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.expected, resp)

		})
	}
}

func TestRepository_StageProfileNudge(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx   context.Context
		nudge *feedlib.Nudge
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:create",
			args: args{
				ctx:   ctx,
				nudge: &feedlib.Nudge{},
			},
			wantErr: false,
		},
		{
			name: "valid:return_internal_server_error",
			args: args{
				ctx:   ctx,
				nudge: &feedlib.Nudge{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:create" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					doc := firestore.DocumentRef{
						ID: uuid.New().String(),
					}
					return &doc, nil
				}
			}

			if tt.name == "valid:return_internal_server_error" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return nil, fmt.Errorf("internal server error")
				}
			}

			err := repo.StageProfileNudge(tt.args.ctx, tt.args.nudge)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

		})
	}
}

func TestRepository_StageKYCProcessingRequest(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx  context.Context
		data *domain.KYCRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:create",
			args: args{
				ctx:  ctx,
				data: &domain.KYCRequest{ID: uuid.New().String()},
			},
			wantErr: false,
		},
		{
			name: "valid:return_internal_server_error",
			args: args{
				ctx:  ctx,
				data: &domain.KYCRequest{ID: uuid.New().String()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:create" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					doc := firestore.DocumentRef{
						ID: uuid.New().String(),
					}
					return &doc, nil
				}
			}

			if tt.name == "valid:return_internal_server_error" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return nil, fmt.Errorf("internal server error")
				}
			}

			err := repo.StageKYCProcessingRequest(tt.args.ctx, tt.args.data)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

		})
	}
}

func TestRepository_UpdateRole(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx  context.Context
		id   string
		role profileutils.RoleType
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid:user_profile_not_found",
			args: args{
				ctx:  ctx,
				id:   "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				role: profileutils.RoleTypeEmployee,
			},
			wantErr: true,
		},
		{
			name: "valid:update_user_role_successful",
			args: args{
				ctx:  ctx,
				id:   "12333",
				role: profileutils.RoleTypeEmployee,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "invalid:user_profile_not_found" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("unable to get user profile docs")
				}
				fakeFireBaseClientExt.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("error: unable to get profile")
				}
			}
			if tt.name == "valid:update_user_role_successful" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("unable to get user profile docs")
				}
				fakeFireBaseClientExt.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
					}, nil
				}
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
				fakeFireStoreClientExt.UpdateFn = func(ctx context.Context, command *fb.UpdateCommand) error {
					return nil
				}
			}

			err := repo.UpdateRole(tt.args.ctx, tt.args.id, tt.args.role)

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

func TestRepository_UpdateFavNavActions(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx        context.Context
		id         string
		favActions []string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid:user_profile_not_found",
			args: args{
				ctx:        ctx,
				id:         "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				favActions: []string{"home"},
			},
			wantErr: true,
		},
		{
			name: "invalid:unable_to_pass_userprofile",
			args: args{
				ctx:        ctx,
				id:         "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				favActions: []string{"home"},
			},
			wantErr: true,
		},
		{
			name: "invalid:user_profile_collection_size_0",
			args: args{
				ctx:        ctx,
				id:         "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				favActions: []string{"home"},
			},
			wantErr: true,
		},
		{
			name: "invalid:unable_update_userprofile_fav_actions",
			args: args{
				ctx:        ctx,
				id:         "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				favActions: []string{"home"},
			},
			wantErr: true,
		},
		{
			name: "valid:update_user_favorite_actions_successful",
			args: args{
				ctx:        ctx,
				id:         "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				favActions: []string{"home"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "invalid:user_profile_not_found" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("unable to get user profile docs")
				}
				fakeFireBaseClientExt.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return nil, fmt.Errorf("error: unable to get profile")
				}
			}
			if tt.name == "invalid:unable_to_pass_userprofile" {
				fakeFireBaseClientExt.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("unable to get user profile docs")
				}
			}
			if tt.name == "invalid:user_profile_collection_size_0" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("unable to get user profile docs")
				}
				fakeFireBaseClientExt.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
			}
			if tt.name == "invalid:unable_update_userprofile_fav_actions" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("unable to get user profile docs")
				}
				fakeFireBaseClientExt.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{}, nil
				}
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{
						{
							Ref: &firestore.DocumentRef{
								ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
							},
							CreateTime: time.Time{},
							UpdateTime: time.Time{},
							ReadTime:   time.Time{},
						},
					}
					return docs, nil
				}
				fakeFireStoreClientExt.UpdateFn = func(ctx context.Context, command *fb.UpdateCommand) error {
					return fmt.Errorf("unable to update user profile")
				}
			}
			if tt.name == "valid:update_user_favorite_actions_successful" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("unable to get user profile docs")
				}
				fakeFireBaseClientExt.GetUserProfileByIDFn = func(ctx context.Context, id string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
					}, nil
				}
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
				fakeFireStoreClientExt.UpdateFn = func(ctx context.Context, command *fb.UpdateCommand) error {
					return nil
				}
			}

			err := repo.UpdateFavNavActions(tt.args.ctx, tt.args.id, tt.args.favActions)

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

func TestRepository_CreateDetailedSupplierProfile(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"

	type args struct {
		ctx       context.Context
		profileID string
		supplier  profileutils.Supplier
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.Supplier
		wantErr bool
	}{
		{
			name: "valid:create_supplier_profile",
			args: args{
				ctx:       ctx,
				profileID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				supplier: profileutils.Supplier{
					ProfileID: &prID,
				},
			},
			want: &profileutils.Supplier{
				ID:        "5e6e41f4-846b-4ba5-ae3f-a92cc7a997ba",
				ProfileID: &prID,
			},
			wantErr: false,
		},
		{
			name: "invalid:create_supplier_profile_firestore_error",
			args: args{
				ctx:       ctx,
				profileID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				supplier: profileutils.Supplier{
					ProfileID: &prID,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:create_supplier_profile_firestore_error",
			args: args{
				ctx:       ctx,
				profileID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				supplier: profileutils.Supplier{
					ProfileID: &prID,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:create_supplier_profile" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}, nil
				}
			}

			if tt.name == "invalid:create_supplier_profile_firestore_error" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return nil, fmt.Errorf("cannot create supplier in firestore")
				}
			}

			got, err := repo.CreateDetailedSupplierProfile(
				tt.args.ctx,
				tt.args.profileID,
				tt.args.supplier,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Repository.CreateDetailedSupplierProfile() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Repository.CreateDetailedSupplierProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_CreateDetailedUserProfile(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	// agent 47
	fName := "Tobias"
	lName := "Rieper"

	type args struct {
		ctx         context.Context
		phoneNumber string
		profile     profileutils.UserProfile
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.UserProfile
		wantErr bool
	}{
		{
			name: "valid:create_user_profile",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				profile: profileutils.UserProfile{
					UserBioData: profileutils.BioData{
						FirstName: &fName,
						LastName:  &lName,
						Gender:    enumutils.GenderMale,
					},
					Role: profileutils.RoleTypeAgent,
				},
			},
			want: &profileutils.UserProfile{
				ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
				UserBioData: profileutils.BioData{
					FirstName: &fName,
					LastName:  &lName,
					Gender:    enumutils.GenderMale,
				},
				Role: profileutils.RoleTypeAgent,
			},
			wantErr: false,
		},
		{
			name: "invalid:create_user_profile_phone_exists_error",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				profile: profileutils.UserProfile{
					UserBioData: profileutils.BioData{
						FirstName: &fName,
						LastName:  &lName,
						Gender:    enumutils.GenderMale,
					},
					Role: profileutils.RoleTypeAgent,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:create_user_profile_phone_exists",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				profile: profileutils.UserProfile{
					UserBioData: profileutils.BioData{
						FirstName: &fName,
						LastName:  &lName,
						Gender:    enumutils.GenderMale,
					},
					Role: profileutils.RoleTypeAgent,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:create_firebase_user_error",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				profile: profileutils.UserProfile{
					UserBioData: profileutils.BioData{
						FirstName: &fName,
						LastName:  &lName,
						Gender:    enumutils.GenderMale,
					},
					Role: profileutils.RoleTypeAgent,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:create_user_profile_firestore_error",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				profile: profileutils.UserProfile{
					UserBioData: profileutils.BioData{
						FirstName: &fName,
						LastName:  &lName,
						Gender:    enumutils.GenderMale,
					},
					Role: profileutils.RoleTypeAgent,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:create_user_profile" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}

				fakeFireBaseClientExt.GetUserByPhoneNumberFn = func(ctx context.Context, phone string) (*auth.UserRecord, error) {
					return nil, nil
				}

				fakeFireBaseClientExt.CreateUserFn = func(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error) {
					return &auth.UserRecord{
						UserInfo: &auth.UserInfo{
							UID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						},
					}, nil
				}

				fakeFireBaseClientExt.GetUserByPhoneNumberFn = func(ctx context.Context, phone string) (*auth.UserRecord, error) {
					return &auth.UserRecord{
						UserInfo: &auth.UserInfo{
							UID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						},
					}, nil
				}

				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}, nil
				}
			}

			if tt.name == "invalid:create_user_profile_phone_exists" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{
						{
							Ref: &firestore.DocumentRef{
								ID: uuid.New().String(),
							},
						},
					}
					return docs, nil
				}

				fakeFireBaseClientExt.GetUserByPhoneNumberFn = func(ctx context.Context, phone string) (*auth.UserRecord, error) {
					return &auth.UserRecord{
						UserInfo: &auth.UserInfo{
							UID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						},
					}, nil
				}

			}

			if tt.name == "invalid:create_user_profile_phone_exists_error" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, fmt.Errorf("cannot profiles matching phone number")
				}
			}

			if tt.name == "invalid:create_firebase_user_error" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}

				fakeFireBaseClientExt.GetUserByPhoneNumberFn = func(ctx context.Context, phone string) (*auth.UserRecord, error) {
					return nil, nil
				}

				fakeFireBaseClientExt.CreateUserFn = func(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error) {
					return nil, fmt.Errorf("cannot create user on firebase")
				}

				fakeFireBaseClientExt.GetUserByPhoneNumberFn = func(ctx context.Context, phone string) (*auth.UserRecord, error) {
					return nil, fmt.Errorf("user doesn't exist")
				}

				fakeFireBaseClientExt.CreateUserFn = func(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error) {
					return nil, fmt.Errorf("cannot create user on firebase")
				}

				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}, nil
				}
			}

			if tt.name == "invalid:create_user_profile_firestore_error" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}

				fakeFireBaseClientExt.GetUserByPhoneNumberFn = func(ctx context.Context, phone string) (*auth.UserRecord, error) {
					return nil, nil
				}

				fakeFireBaseClientExt.CreateUserFn = func(ctx context.Context, user *auth.UserToCreate) (*auth.UserRecord, error) {
					return nil, fmt.Errorf("cannot create user on firebase")
				}

				fakeFireBaseClientExt.GetUserByPhoneNumberFn = func(ctx context.Context, phone string) (*auth.UserRecord, error) {
					return &auth.UserRecord{
						UserInfo: &auth.UserInfo{
							UID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
						},
					}, nil
				}

				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return nil, fmt.Errorf("cannot create user on firestore")
				}
			}

			got, err := repo.CreateDetailedUserProfile(
				tt.args.ctx,
				tt.args.phoneNumber,
				tt.args.profile,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Repository.CreateDetailedUserProfile() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Repository.CreateDetailedUserProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_ListAgentUserProfiles(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx  context.Context
		role profileutils.RoleType
	}
	tests := []struct {
		name    string
		args    args
		want    []*profileutils.UserProfile
		wantErr bool
	}{
		{
			name: "success:fetch_agent_user_profiles",
			args: args{
				ctx:  ctx,
				role: profileutils.RoleTypeEmployee,
			},
			want:    []*profileutils.UserProfile{},
			wantErr: false,
		},
		{
			name: "fail:fetch_agent_user_profiles_error",
			args: args{
				ctx:  ctx,
				role: profileutils.RoleTypeAgent,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "success:fetch_agent_user_profiles" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
			}

			if tt.name == "fail:fetch_agent_user_profiles_error" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {

					return nil, fmt.Errorf("cannot fetch firebase docs")
				}
			}

			got, err := repo.ListUserProfiles(tt.args.ctx, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Repository.ListAgentUserProfiles() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.ListAgentUserProfiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_AddAITSessionDetails_Unittest(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	phoneNumber := "+254700100200"
	SessionID := uuid.New().String()
	Level := 0
	Text := ""

	sessionDet := &dto.SessionDetails{
		SessionID:   SessionID,
		PhoneNumber: &phoneNumber,
		Level:       Level,
		Text:        Text,
	}

	type args struct {
		ctx   context.Context
		input *dto.SessionDetails
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.USSDLeadDetails
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:   ctx,
				input: sessionDet,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:   ctx,
				input: sessionDet,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}

				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}, nil
				}
			}

			if tt.name == "Sad case" {
				_, err := utils.ValidateUSSDDetails(sessionDet)
				if err != nil {
					t.Errorf("an error occurred")
					return
				}

				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return nil, fmt.Errorf("error")
				}

			}

			got, err := repo.AddAITSessionDetails(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Repository.AddAITSessionDetails() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.AddAITSessionDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetAITSessionDetails_Unittests(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	SessionID := uuid.New().String()

	type args struct {
		ctx       context.Context
		sessionID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.USSDLeadDetails
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				sessionID: SessionID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				sessionID: SessionID,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := repo.GetAITSessionDetails(tt.args.ctx, tt.args.sessionID)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Repository.GetAITSessionDetails() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetAITSessionDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetAITDetails_Unnittest(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	phoneNumber := "+254700100200"

	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.USSDLeadDetails
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
			}

			if tt.name == "Sad case" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := repo.GetAITDetails(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetAITDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.GetAITDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_CreateRole(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	role := &profileutils.Role{
		ID:          uuid.NewString(),
		Name:        "Accountant",
		Description: "Can handle money",
	}

	type args struct {
		ctx       context.Context
		profileID string
		input     dto.RoleInput
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.Role
		wantErr bool
	}{
		{
			name: "valid:create a new role",
			args: args{
				ctx:       ctx,
				profileID: uuid.NewString(),
				input: dto.RoleInput{
					Name:        "Accountant",
					Description: "Can handle money",
				},
			},
			want:    role,
			wantErr: false,
		},
		{
			name: "invalid: role with similar name exists",
			args: args{
				ctx:       ctx,
				profileID: uuid.NewString(),
				input: dto.RoleInput{
					Name:        "Accountant",
					Description: "Can handle money",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid: check role with similar name error",
			args: args{
				ctx:       ctx,
				profileID: uuid.NewString(),
				input: dto.RoleInput{
					Name:        "Accountant",
					Description: "Can handle money",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid: firestore create role error",
			args: args{
				ctx:       ctx,
				profileID: uuid.NewString(),
				input: dto.RoleInput{
					Name:        "Accountant",
					Description: "Can handle money",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:create a new role" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					return docRef, nil
				}
			}

			if tt.name == "invalid: role with similar name exists" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}}
					return docs, nil
				}
			}

			if tt.name == "invalid: check role with similar name error" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {

					return nil, fmt.Errorf("cannot get all")
				}
			}

			if tt.name == "invalid: firestore create role error" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {

					return nil, fmt.Errorf("cannot create role")
				}
			}

			got, err := repo.CreateRole(tt.args.ctx, tt.args.profileID, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.CreateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Repository.CreateRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_UpdateRoleDetails(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	profileID := uuid.NewString()
	role := profileutils.Role{
		ID:          uuid.NewString(),
		Name:        "Accountant",
		Description: "Can handle money",
	}

	type args struct {
		ctx       context.Context
		profileID string
		role      profileutils.Role
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid: cannot retrive role",
			args: args{
				ctx:       ctx,
				profileID: profileID,
				role:      role,
			},
			wantErr: true,
		},
		{
			name: "invalid: cannot update role",
			args: args{
				ctx:       ctx,
				profileID: profileID,
				role:      role,
			},
			wantErr: true,
		},
		{
			name: "valid:success update role",
			args: args{
				ctx:       ctx,
				profileID: profileID,
				role:      role,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "invalid: cannot retrive role" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {

					return nil, fmt.Errorf("cannot retrieve role")
				}
			}

			if tt.name == "invalid: cannot update role" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}}
					return docs, nil
				}
				fakeFireStoreClientExt.UpdateFn = func(ctx context.Context, command *fb.UpdateCommand) error {
					return fmt.Errorf("cannot update the role")
				}
			}

			if tt.name == "valid:success update role" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}}
					return docs, nil
				}
				fakeFireStoreClientExt.UpdateFn = func(ctx context.Context, command *fb.UpdateCommand) error {
					return nil
				}
			}

			_, err := repo.UpdateRoleDetails(tt.args.ctx, tt.args.profileID, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateRoleDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRepository_GetRoleByID(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx    context.Context
		roleID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "fail: cannot retrieve firestore docs",
			args: args{
				ctx:    ctx,
				roleID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "fail: multiple firestore records",
			args: args{
				ctx:    ctx,
				roleID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "fail: no firestore record",
			args: args{
				ctx:    ctx,
				roleID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "fail: cannot convert to role value",
			args: args{
				ctx:    ctx,
				roleID: uuid.NewString(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "fail: cannot retrieve firestore docs" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {

					return nil, fmt.Errorf("cannot retrieve firestore docs")
				}
			}

			if tt.name == "fail: multiple firestore records" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docRef2 := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}, {Ref: docRef2}}
					return docs, nil
				}
			}

			if tt.name == "fail: no firestore record" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
			}

			if tt.name == "success: retrieve role" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}}
					return docs, nil
				}
			}

			got, err := repo.GetRoleByID(tt.args.ctx, tt.args.roleID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetRoleByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Repository.GetAllRoles() = %v", got)
			}
		})
	}
}

func TestRepository_CheckIfRoleNameExists(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "fail: cannot retrieve firestore documents",
			args: args{
				ctx:  ctx,
				name: "accountant",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "success: similar role exists",
			args: args{
				ctx:  ctx,
				name: "accountant",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "success: similar role doesn't exist",
			args: args{
				ctx:  ctx,
				name: "accountant",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "fail: cannot retrieve firestore documents" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}}
					return docs, fmt.Errorf("cannot get firestore docs")
				}
			}

			if tt.name == "success: similar role exists" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}}
					return docs, nil
				}
			}

			if tt.name == "success: similar role doesn't exist" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
			}

			got, err := repo.CheckIfRoleNameExists(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Repository.CheckIfRoleNameExists() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if got != tt.want {
				t.Errorf("Repository.CheckIfRoleNameExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetRolesByIDs(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx     context.Context
		roleIDs []string
	}
	tests := []struct {
		name    string
		args    args
		want    *[]profileutils.Role
		wantErr bool
	}{

		{
			name: "fail: cannot retrieve role by id",
			args: args{
				ctx:     ctx,
				roleIDs: []string{uuid.NewString(), uuid.NewString()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		if tt.name == "success: retrieve role by id" {
			fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
				docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
				docs := []*firestore.DocumentSnapshot{{Ref: docRef}}
				return docs, nil
			}
		}

		if tt.name == "fail: cannot retrieve role by id" {
			fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {

				return nil, fmt.Errorf("cannot retrieve role")
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			_, err := repo.GetRolesByIDs(tt.args.ctx, tt.args.roleIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetRolesByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRepository_UpdateUserRoleIDs(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx     context.Context
		id      string
		roleIDs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "fail:cannot retrieve user profile",
			args: args{
				ctx:     ctx,
				id:      uuid.NewString(),
				roleIDs: []string{uuid.NewString()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "fail:cannot retrieve user profile" {

				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					ref := firestore.DocumentRef{ID: "123"}
					docs := []*firestore.DocumentSnapshot{{Ref: &ref}}
					return docs, nil
				}

				fakeFireStoreClientExt.UpdateFn = func(ctx context.Context, command *fb.UpdateCommand) error {
					return nil
				}
			}

			if err := repo.UpdateUserRoleIDs(tt.args.ctx, tt.args.id, tt.args.roleIDs); (err != nil) != tt.wantErr {
				t.Errorf("Repository.UpdateUserRoleIDs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepository_GetAllRoles(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *[]profileutils.Role
		wantErr bool
	}{
		{
			name: "fail: cannot retrieve firestore docs",
			args: args{
				ctx: ctx,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "fail: cannot retrieve firestore docs" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {

					return nil, fmt.Errorf("cannot retrieve firestore docs")
				}
			}

			got, err := repo.GetAllRoles(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetAllRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Repository.GetAllRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_CheckIfUserHasPermission(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx                context.Context
		UID                string
		requiredPermission profileutils.Permission
	}
	input := args{
		ctx:                ctx,
		UID:                uuid.NewString(),
		requiredPermission: profileutils.CanAssignRole,
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "sad unable to get user profile",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "sad unable to get roles by ids",
			args:    input,
			want:    false,
			wantErr: true,
		},
		{
			name:    "sad required permission not among user roles",
			args:    input,
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad unable to get user profile" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("unable to get userprofile by uid")
				}
			}
			if tt.name == "sad unable to get roles by ids" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}}
					return docs, nil
				}
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("unable to get userprofile by uid")
				}
			}
			if tt.name == "sad required permission not among user roles" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}}
					return docs, nil
				}
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{
						ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
					}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}}
					return docs, nil
				}
			}
			got, err := repo.CheckIfUserHasPermission(
				tt.args.ctx,
				tt.args.UID,
				tt.args.requiredPermission,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf(
					"Repository.CheckIfUserHasPermission() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if got != tt.want {
				t.Errorf("Repository.CheckIfUserHasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetRoleByName(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx      context.Context
		roleName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "fail: cannot retrieve firestore docs",
			args: args{
				ctx:      ctx,
				roleName: "Test Role",
			},
			wantErr: true,
		},
		{
			name: "fail: multiple firestore records",
			args: args{
				ctx:      ctx,
				roleName: "Test Role",
			},
			wantErr: true,
		},
		{
			name: "fail: no firestore record",
			args: args{
				ctx:      ctx,
				roleName: "Test Role",
			},
			wantErr: true,
		},
		{
			name: "fail: cannot convert to role value",
			args: args{
				ctx:      ctx,
				roleName: "Test Role",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "fail: cannot retrieve firestore docs" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {

					return nil, fmt.Errorf("cannot retrieve firestore docs")
				}
			}

			if tt.name == "fail: multiple firestore records" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docRef2 := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}, {Ref: docRef2}}
					return docs, nil
				}
			}

			if tt.name == "fail: no firestore record" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}
			}

			if tt.name == "success: retrieve role" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}}
					return docs, nil
				}
			}

			got, err := repo.GetRoleByName(tt.args.ctx, tt.args.roleName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetRoleByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Repository.GetRoleByName() = %v", got)
			}
		})
	}
}

func TestRepository_GetUserProfilesByRoleID(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx    context.Context
		roleID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "fail: cannot retrieve firestore docs",
			args: args{
				ctx:    ctx,
				roleID: uuid.NewString(),
			},
			wantErr: true,
		},
		{
			name: "fail: cannot parse user profile",
			args: args{
				ctx:    ctx,
				roleID: uuid.NewString(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "fail: cannot retrieve firestore docs" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, fmt.Errorf("cannot retrieve docs")
				}
			}

			if tt.name == "fail: cannot parse user profile" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docRef := &firestore.DocumentRef{ID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f"}
					docs := []*firestore.DocumentSnapshot{{Ref: docRef}}

					return docs, nil
				}
			}

			got, err := repo.GetUserProfilesByRoleID(tt.args.ctx, tt.args.roleID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.GetUserProfilesByRoleID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("Repository.GetUserProfilesByRoleID() = %v", got)
			}
		})
	}
}

func TestRepository_SaveRoleRevocation(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx        context.Context
		userID     string
		revocation dto.RoleRevocationInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "fail: cannot save a role revocation",
			args: args{
				ctx:    ctx,
				userID: uuid.NewString(),
				revocation: dto.RoleRevocationInput{
					ProfileID: uuid.NewString(),
					RoleID:    uuid.NewString(),
					Reason:    "test reason",
				},
			},
			wantErr: true,
		},
		{
			name: "success: save a role revocation",
			args: args{
				ctx:    ctx,
				userID: uuid.NewString(),
				revocation: dto.RoleRevocationInput{
					ProfileID: uuid.NewString(),
					RoleID:    uuid.NewString(),
					Reason:    "test reason",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "fail: cannot save a role revocation" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return nil, fmt.Errorf("cannot create firestore document")
				}
			}

			if tt.name == "success: save a role revocation" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
					return nil, nil
				}
			}

			if err := repo.SaveRoleRevocation(tt.args.ctx, tt.args.userID, tt.args.revocation); (err != nil) != tt.wantErr {
				t.Errorf("Repository.SaveRoleRevocation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
