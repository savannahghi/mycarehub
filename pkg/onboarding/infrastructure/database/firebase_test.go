package database_test

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
	extMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database/mock"
)

var fakeFireBaseClientExt extMock.FirebaseClientExtension
var fireBaseClientExt database.FirebaseClientExtension = &fakeFireBaseClientExt

var fakeFireStoreClientExt extMock.FirestoreClientExtension

func TestRepository_UpdateUserName(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt database.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := database.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:update_user_name_failed_to_get_a user_profile" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *database.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}

				fakeFireStoreClientExt.UpdateFn = func(ctx context.Context, command *database.UpdateCommand) error {
					return nil
				}
			}

			if tt.name == "invalid:user_name_already_exists" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *database.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{
						{
							Ref: &firestore.DocumentRef{
								ID: "5555",
							},
						},
					}
					return docs, nil
				}

				fakeFireStoreClientExt.UpdateFn = func(ctx context.Context, command *database.UpdateCommand) error {
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
	var fireStoreClientExt database.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := database.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

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
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *database.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
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
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *database.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
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
	var fireStoreClientExt database.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := database.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx     context.Context
		profile *base.UserProfile
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
				profile: &base.UserProfile{
					ID: uuid.New().String(),
				},
			},
			expected: true,
		},
		{
			name: "valid:already_exists",
			args: args{
				ctx: ctx,
				profile: &base.UserProfile{
					ID: uuid.New().String(),
				},
			},
			expected: true,
		},

		{
			name: "invalid:throws_internal_server_error_while_checking_existence",
			args: args{
				ctx: ctx,
				profile: &base.UserProfile{
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
				profile: &base.UserProfile{
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
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *database.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}

				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *database.CreateCommand) (*firestore.DocumentRef, error) {
					doc := firestore.DocumentRef{
						ID: uuid.New().String(),
					}
					return &doc, nil
				}
			}

			if tt.name == "valid:already_exists" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *database.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
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
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *database.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					return nil, exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot"))
				}
			}

			if tt.name == "invalid:throws_internal_server_error_while_creating" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *database.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}

				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *database.CreateCommand) (*firestore.DocumentRef, error) {
					return nil, exceptions.InternalServerError(fmt.Errorf("unable to add user profile of ID in experiment_participant"))
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
	var fireStoreClientExt database.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := database.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx     context.Context
		profile *base.UserProfile
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
				profile: &base.UserProfile{
					ID: uuid.New().String(),
				},
			},
			expected: true,
		},

		{
			name: "invalid:throws_internal_server_error_while_removing",
			args: args{
				ctx: ctx,
				profile: &base.UserProfile{
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
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *database.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{
						{
							Ref: &firestore.DocumentRef{
								ID: uuid.New().String(),
							},
						},
					}
					return docs, nil
				}

				fakeFireStoreClientExt.DeleteFn = func(ctx context.Context, command *database.DeleteCommand) error {
					return nil
				}

			}
			if tt.name == "invalid:throws_internal_server_error_while_removing" {
				fakeFireStoreClientExt.DeleteFn = func(ctx context.Context, command *database.DeleteCommand) error {
					return exceptions.InternalServerError(fmt.Errorf("unable to remove user profile of ID  from experiment_participant"))
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
	var fireStoreClientExt database.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := database.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx   context.Context
		nudge map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:create",
			args: args{
				ctx: ctx,
				nudge: map[string]interface{}{
					"name": "valid",
				},
			},
			wantErr: false,
		},
		{
			name: "valid:return_internal_server_error",
			args: args{
				ctx: ctx,
				nudge: map[string]interface{}{
					"name": "valid",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:create" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *database.CreateCommand) (*firestore.DocumentRef, error) {
					doc := firestore.DocumentRef{
						ID: uuid.New().String(),
					}
					return &doc, nil
				}
			}

			if tt.name == "valid:return_internal_server_error" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *database.CreateCommand) (*firestore.DocumentRef, error) {
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
	var fireStoreClientExt database.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := database.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

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
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *database.CreateCommand) (*firestore.DocumentRef, error) {
					doc := firestore.DocumentRef{
						ID: uuid.New().String(),
					}
					return &doc, nil
				}
			}

			if tt.name == "valid:return_internal_server_error" {
				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *database.CreateCommand) (*firestore.DocumentRef, error) {
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
