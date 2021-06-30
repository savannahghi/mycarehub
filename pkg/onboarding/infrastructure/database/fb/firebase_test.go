package fb_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database/fb"
	extMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database/fb/mock"
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
					return nil, exceptions.InternalServerError(fmt.Errorf("unable to parse user profile as firebase snapshot"))
				}
			}

			if tt.name == "invalid:throws_internal_server_error_while_creating" {
				fakeFireStoreClientExt.GetAllFn = func(ctx context.Context, query *fb.GetAllQuery) ([]*firestore.DocumentSnapshot, error) {
					docs := []*firestore.DocumentSnapshot{}
					return docs, nil
				}

				fakeFireStoreClientExt.CreateFn = func(ctx context.Context, command *fb.CreateCommand) (*firestore.DocumentRef, error) {
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
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

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
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	type args struct {
		ctx   context.Context
		nudge *base.Nudge
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
				nudge: &base.Nudge{},
			},
			wantErr: false,
		},
		{
			name: "valid:return_internal_server_error",
			args: args{
				ctx:   ctx,
				nudge: &base.Nudge{},
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
		role base.RoleType
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:update_user_role_successful",
			args: args{
				ctx:  ctx,
				id:   "12333",
				role: base.RoleTypeEmployee,
			},
			wantErr: true,
		},
		{
			name: "invalid:update_user_role_failed_userprofile_not_found",
			args: args{
				ctx:  ctx,
				id:   "12333",
				role: base.RoleTypeAgent,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestRepository_CreateDetailedSupplierProfile(t *testing.T) {
	ctx := context.Background()
	var fireStoreClientExt fb.FirestoreClientExtension = &fakeFireStoreClientExt
	repo := fb.NewFirebaseRepository(fireStoreClientExt, fireBaseClientExt)

	prID := "c9d62c7e-93e5-44a6-b503-6fc159c1782f"

	type args struct {
		ctx       context.Context
		profileID string
		supplier  base.Supplier
	}
	tests := []struct {
		name    string
		args    args
		want    *base.Supplier
		wantErr bool
	}{
		{
			name: "valid:create_supplier_profile",
			args: args{
				ctx:       ctx,
				profileID: "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				supplier: base.Supplier{
					ProfileID: &prID,
				},
			},
			want: &base.Supplier{
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
				supplier: base.Supplier{
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
				supplier: base.Supplier{
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

			got, err := repo.CreateDetailedSupplierProfile(tt.args.ctx, tt.args.profileID, tt.args.supplier)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.CreateDetailedSupplierProfile() error = %v, wantErr %v", err, tt.wantErr)
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
		profile     base.UserProfile
	}
	tests := []struct {
		name    string
		args    args
		want    *base.UserProfile
		wantErr bool
	}{
		{
			name: "valid:create_user_profile",
			args: args{
				ctx:         ctx,
				phoneNumber: base.TestUserPhoneNumber,
				profile: base.UserProfile{
					UserBioData: base.BioData{
						FirstName: &fName,
						LastName:  &lName,
						Gender:    base.GenderMale,
					},
					Role: base.RoleTypeAgent,
				},
			},
			want: &base.UserProfile{
				ID:           "c9d62c7e-93e5-44a6-b503-6fc159c1782f",
				VerifiedUIDS: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e"},
				UserBioData: base.BioData{
					FirstName: &fName,
					LastName:  &lName,
					Gender:    base.GenderMale,
				},
				Role: base.RoleTypeAgent,
			},
			wantErr: false,
		},
		{
			name: "invalid:create_user_profile_phone_exists_error",
			args: args{
				ctx:         ctx,
				phoneNumber: base.TestUserPhoneNumber,
				profile: base.UserProfile{
					UserBioData: base.BioData{
						FirstName: &fName,
						LastName:  &lName,
						Gender:    base.GenderMale,
					},
					Role: base.RoleTypeAgent,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:create_user_profile_phone_exists",
			args: args{
				ctx:         ctx,
				phoneNumber: base.TestUserPhoneNumber,
				profile: base.UserProfile{
					UserBioData: base.BioData{
						FirstName: &fName,
						LastName:  &lName,
						Gender:    base.GenderMale,
					},
					Role: base.RoleTypeAgent,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:create_firebase_user_error",
			args: args{
				ctx:         ctx,
				phoneNumber: base.TestUserPhoneNumber,
				profile: base.UserProfile{
					UserBioData: base.BioData{
						FirstName: &fName,
						LastName:  &lName,
						Gender:    base.GenderMale,
					},
					Role: base.RoleTypeAgent,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:create_user_profile_firestore_error",
			args: args{
				ctx:         ctx,
				phoneNumber: base.TestUserPhoneNumber,
				profile: base.UserProfile{
					UserBioData: base.BioData{
						FirstName: &fName,
						LastName:  &lName,
						Gender:    base.GenderMale,
					},
					Role: base.RoleTypeAgent,
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

			got, err := repo.CreateDetailedUserProfile(tt.args.ctx, tt.args.phoneNumber, tt.args.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.CreateDetailedUserProfile() error = %v, wantErr %v", err, tt.wantErr)
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
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*base.UserProfile
		wantErr bool
	}{
		{
			name: "success:fetch_agent_user_profiles",
			args: args{
				ctx: ctx,
			},
			want:    []*base.UserProfile{},
			wantErr: false,
		},
		{
			name: "fail:fetch_agent_user_profiles_error",
			args: args{
				ctx: ctx,
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

			got, err := repo.ListAgentUserProfiles(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repository.ListAgentUserProfiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.ListAgentUserProfiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
