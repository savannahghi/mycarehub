package usecases_test

import (
	"context"
	"log"

	"testing"

	"gitlab.slade360emr.com/go/base"
)

func TestMaskPhoneNumbers(t *testing.T) {

	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		phones []string
	}

	tests := []struct {
		name string
		arg  args
		want []string
	}{
		{
			name: "valid case",
			arg: args{
				phones: []string{"+254789874267"},
			},
			want: []string{"+254789***267"},
		},
		{
			name: "valid case < 10 digits",
			arg: args{
				phones: []string{"+2547898742"},
			},
			want: []string{"+2547***742"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maskedPhone := s.Onboarding.MaskPhoneNumbers(tt.arg.phones)
			if len(maskedPhone) != len(tt.want) {
				t.Errorf("returned masked phone number not the expected one, wanted: %v got: %v", tt.want, maskedPhone)
				return
			}

			for i, number := range maskedPhone {
				if tt.want[i] != number {
					t.Errorf("wanted: %v, got: %v", tt.want[i], number)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_GetUserProfileByUID(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	type args struct {
		ctx context.Context
		UID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sucess: get a user profile given their UID",
			args: args{
				ctx: ctx,
				UID: auth.UID,
			},
			wantErr: false,
		},
		{
			name: "failure: fail get a user profile given a bad UID",
			args: args{
				ctx: ctx,
				UID: "not-a-valid-uid",
			},
			wantErr: true,
		},
		{
			name: "failure: fail get a user profile given an empty UID",
			args: args{
				ctx: ctx,
				UID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := s.Onboarding.GetUserProfileByUID(tt.args.ctx, tt.args.UID)
			if tt.wantErr && profile != nil {
				t.Errorf("expected nil but got %v, since the error %v occurred",
					profile,
					err,
				)
				return
			}

			if !tt.wantErr && profile == nil {
				t.Errorf("expected a profile but got nil, since no error occurred")
				return
			}

		})
	}
}

func TestProfileUseCaseImpl_UserProfile(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *base.UserProfile
		wantErr bool
	}{
		{
			name: "valid: user profile retrieved",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context supplied",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Onboarding.UserProfile(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.UserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("nil user profile returned")
				return
			}
		})
	}
}

func TestProfileUseCaseImpl_GetProfileByID(t *testing.T) {

	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
		return
	}

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	profile, err := s.Onboarding.UserProfile(ctx)
	if err != nil {
		t.Errorf("could not retreive user profile")
		return
	}

	type args struct {
		ctx context.Context
		id  *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: user profile retreived",
			args: args{
				ctx: ctx,
				id:  &profile.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Onboarding.GetProfileByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.GetProfileByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("nil user profile returned")
				return
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdateBioData(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
		return
	}

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	dateOfBirth := base.Date{
		Day:   12,
		Year:  2000,
		Month: 2,
	}

	firstName := "Jatelo"
	lastName := "Omera"
	bioData := base.BioData{
		FirstName:   &firstName,
		LastName:    &lastName,
		DateOfBirth: &dateOfBirth,
	}

	var gender base.Gender = "female"
	updateGender := base.BioData{
		Gender: gender,
	}

	updateDOB := base.BioData{
		DateOfBirth: &dateOfBirth,
	}

	updateFirstName := base.BioData{
		FirstName: &firstName,
	}

	updateLastName := base.BioData{
		LastName: &lastName,
	}

	type args struct {
		ctx  context.Context
		data base.BioData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully update biodata",
			args: args{
				ctx:  ctx,
				data: bioData,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the firstname",
			args: args{
				ctx:  ctx,
				data: updateFirstName,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the lastname",
			args: args{
				ctx:  ctx,
				data: updateLastName,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the date of birth",
			args: args{
				ctx:  ctx,
				data: updateDOB,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the gender",
			args: args{
				ctx:  ctx,
				data: updateGender,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unauthenticated context",
			args: args{
				ctx:  context.Background(),
				data: bioData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Onboarding.UpdateBioData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.UpdateBioData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdatePhotoUploadID(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
		return
	}

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		t.Errorf("could not get the logged in user")
		return
	}

	profile, err := s.Onboarding.GetUserProfileByUID(ctx, uid)
	if err != nil {
		t.Errorf("could not retrieve user profile")
		return
	}

	uploadID := "some-photo-upload-id"
	log.Printf("THE UPLOAD ID IS %v", profile.PhotoUploadID)

	type args struct {
		ctx      context.Context
		uploadID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully update the PhotoUploadID",
			args: args{
				ctx:      ctx,
				uploadID: uploadID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Use an unauthenticated context",
			args: args{
				ctx:      context.Background(),
				uploadID: uploadID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Onboarding.UpdatePhotoUploadID(tt.args.ctx, tt.args.uploadID); (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.UpdatePhotoUploadID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
