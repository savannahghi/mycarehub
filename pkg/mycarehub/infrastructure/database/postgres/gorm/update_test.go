package gorm_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

func TestPGInstance_InactivateFacility(t *testing.T) {

	ctx := context.Background()

	type args struct {
		ctx     context.Context
		mflCode *int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx:     ctx,
				mflCode: &mflCodeToInactivate,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx:     ctx,
				mflCode: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.InactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.InactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.InactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_ReactivateFacility(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		mflCode *int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx:     ctx,
				mflCode: &inactiveMflCode,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx:     ctx,
				mflCode: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ReactivateFacility(tt.args.ctx, tt.args.mflCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ReactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.ReactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_SetNickname(t *testing.T) {
	ctx := context.Background()

	invalidUserID := ksuid.New().String()
	invalidNickname := gofakeit.HipsterSentence(50)

	type args struct {
		ctx      context.Context
		userID   *string
		nickname *string
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
				ctx:      ctx,
				userID:   &userID,
				nickname: &userNickname,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				userID:   &invalidUserID,
				nickname: &userNickname,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:      ctx,
				userID:   nil,
				nickname: &userNickname,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no nickname",
			args: args{
				ctx:      ctx,
				userID:   &userID,
				nickname: &invalidNickname,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SetNickName(tt.args.ctx, tt.args.userID, tt.args.nickname)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SetNickName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.SetNickName() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestPGInstance_InvalidatePIN(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
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
				ctx:     ctx,
				userID:  userIDToInvalidate,
				flavour: feedlib.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: no user id provided",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.InvalidatePIN(tt.args.ctx, tt.args.userID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.InvalidatePIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.InvalidatePIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_UpdateIsCorrectSecurityQuestionResponse(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx                               context.Context
		userID                            string
		isCorrectSecurityQuestionResponse bool
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
				ctx:                               ctx,
				userID:                            userID,
				isCorrectSecurityQuestionResponse: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: invalid user id",
			args: args{
				ctx:                               ctx,
				userID:                            uuid.New().String(),
				isCorrectSecurityQuestionResponse: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: empty user id",
			args: args{
				ctx:                               ctx,
				userID:                            "",
				isCorrectSecurityQuestionResponse: true,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.UpdateIsCorrectSecurityQuestionResponse(tt.args.ctx, tt.args.userID, tt.args.isCorrectSecurityQuestionResponse)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateIsCorrectSecurityQuestionResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UpdateIsCorrectSecurityQuestionResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_AcceptTerms(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		userID  *string
		termsID *int
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
				ctx:     ctx,
				userID:  &userIDToAcceptTerms,
				termsID: &termsID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: missing args",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: no userID",
			args: args{
				ctx:     ctx,
				userID:  nil,
				termsID: &termsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: no terms",
			args: args{
				ctx:     ctx,
				userID:  &userIDToAcceptTerms,
				termsID: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.AcceptTerms(tt.args.ctx, tt.args.userID, tt.args.termsID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.AcceptTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.AcceptTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_UpdateUserFailedLoginCount(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx                 context.Context
		userID              string
		failedLoginAttempts int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				ctx:                 ctx,
				userID:              userIDToIncreaseFailedLoginCount,
				failedLoginAttempts: 1,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:                 ctx,
				userID:              "",
				failedLoginAttempts: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := testingDB.UpdateUserFailedLoginCount(tt.args.ctx, tt.args.userID, tt.args.failedLoginAttempts); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserFailedLoginCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateUserLastFailedLoginTime(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				ctx:    ctx,
				userID: userIDtoUpdateLastFailedLoginTime,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserLastFailedLoginTime(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserLastFailedLoginTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateUserNextAllowedLoginTime(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx                  context.Context
		userID               string
		nextAllowedLoginTime time.Time
	}
	tests := []struct {
		name string

		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:                  ctx,
				userID:               userIDToUpdateNextAllowedLoginTime,
				nextAllowedLoginTime: time.Now().Add(3),
			},
			wantErr: false,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:                  ctx,
				userID:               "",
				nextAllowedLoginTime: time.Now().Add(3),
			},
			wantErr: true,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserNextAllowedLoginTime(tt.args.ctx, tt.args.userID, tt.args.nextAllowedLoginTime); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserNextAllowedLoginTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestPGInstance_ShareContent(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx   context.Context
		input dto.ShareContentInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx: ctx,
				input: dto.ShareContentInput{
					UserID:    userID,
					ContentID: contentID,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing input",
			args: args{
				ctx:   ctx,
				input: dto.ShareContentInput{},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ShareContent(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ShareContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.ShareContent() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestPGInstance_BookmarkContent(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx       context.Context
		userID    string
		contentID int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				ctx:       ctx,
				userID:    userID,
				contentID: contentID,
			},
			wantErr: false,
			want:    true,
		},
		{
			// Ensures there is idepotency
			name: "bookmark already exists",
			args: args{
				ctx:       ctx,
				userID:    userID,
				contentID: contentID,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invald: missing parama",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.BookmarkContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.BookmarkContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.BookmarkContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_UnBookmarkContent(t *testing.T) {
	ctx := context.Background()

	_, err := testingDB.BookmarkContent(ctx, userID, contentID2)
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	type args struct {
		ctx       context.Context
		userID    string
		contentID int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "default case",
			args: args{
				ctx:       ctx,
				userID:    userID,
				contentID: contentID2,
			},
			wantErr: false,
			want:    true,
		},
		{
			// Ensures there is idempotency
			name: "bookmark already exists",
			args: args{
				ctx:       ctx,
				userID:    userID,
				contentID: contentID2,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invald: missing params",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.UnBookmarkContent(tt.args.ctx, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UnBookmarkContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UnBookmarkContent() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestPGInstance_CompleteOnboardingTour(t *testing.T) {
	ctx := context.Background()
	flavour := feedlib.FlavourConsumer

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name string

		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:     ctx,
				userID:  userIDUpdatePinRequireChangeStatus,
				flavour: flavour,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - Empty userID and flavour",
			args: args{
				ctx:     ctx,
				userID:  "",
				flavour: "invalid-flavour",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - Invalid flavour",
			args: args{
				ctx:     ctx,
				userID:  userIDUpdatePinRequireChangeStatus,
				flavour: "invalid-flavour",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.CompleteOnboardingTour(tt.args.ctx, tt.args.userID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CompleteOnboardingTour() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CompleteOnboardingTour() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_LikeContent(t *testing.T) {
	ctx := context.Background()

	type args struct {
		context   context.Context
		userID    string
		contentID int
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
				context:   ctx,
				userID:    userID,
				contentID: contentID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - no userID",
			args: args{
				context:   ctx,
				userID:    "",
				contentID: contentID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				context:   ctx,
				userID:    userID,
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and contentID",
			args: args{
				context:   ctx,
				userID:    "",
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.LikeContent(tt.args.context, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.LikeContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.LikeContent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_UnlikeContent(t *testing.T) {
	ctx := context.Background()

	_, err := testingDB.LikeContent(ctx, userID, contentID2)
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	type args struct {
		context   context.Context
		userID    string
		contentID int
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
				context:   ctx,
				userID:    userID,
				contentID: contentID2,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - no userID",
			args: args{
				context:   ctx,
				userID:    "",
				contentID: contentID2,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no contentID",
			args: args{
				context:   ctx,
				userID:    userID,
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and contentID",
			args: args{
				context:   ctx,
				userID:    "",
				contentID: 0,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.UnlikeContent(tt.args.context, tt.args.userID, tt.args.contentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UnlikeContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UnlikeContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestPGInstance_ViewContent(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx       context.Context
		UserID    string
		ContentID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:       ctx,
				UserID:    userID,
				ContentID: contentID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: missing input",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ViewContent(tt.args.ctx, tt.args.UserID, tt.args.ContentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != got {
				t.Errorf("expected %v, but got: %v", tt.want, got)
				return
			}
		})
	}
}

func TestPGInstance_UpdateClientCaregiver(t *testing.T) {
	ctx := context.Background()

	caregiverInput := dto.CaregiverInput{
		ClientID:      clientID,
		FirstName:     "Updated",
		LastName:      "Updated",
		PhoneNumber:   "+1234567890",
		CaregiverType: enums.CaregiverTypeMother,
	}

	type args struct {
		ctx            context.Context
		caregiverInput *dto.CaregiverInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:            ctx,
				caregiverInput: &caregiverInput,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testingDB.UpdateClientCaregiver(tt.args.ctx, tt.args.caregiverInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateClientCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_UpdateUserProfileAfterLoginSuccess(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:    ctx,
				userID: userIDToUpdateUserProfileAfterLoginSuccess,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				userID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserProfileAfterLoginSuccess(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserProfileAfterLoginSuccess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_InProgressBy(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx       context.Context
		requestID string
		staffID   string
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
				ctx:       ctx,
				requestID: clientsServiceRequestID,
				staffID:   staffID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - invalid request ID",
			args: args{
				ctx:       ctx,
				requestID: "clientsServiceRequestID",
				staffID:   staffID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - invalid staff ID",
			args: args{
				ctx:       ctx,
				requestID: clientsServiceRequestID,
				staffID:   "staffID",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid request and  staff ID",
			args: args{
				ctx:       ctx,
				requestID: "clientsServiceRequestID",
				staffID:   "staffID",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - invalid uuid for  staff ID",
			args: args{
				ctx:       ctx,
				requestID: clientsServiceRequestID,
				staffID:   uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SetInProgressBy(tt.args.ctx, tt.args.requestID, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SetInProgressBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.SetInProgressBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_ResolveServiceRequest(t *testing.T) {
	ctx := context.Background()
	longWord := gofakeit.HipsterSentence(10)
	nonExistentUUID := uuid.New().String()

	type args struct {
		ctx              context.Context
		staffID          *string
		serviceRequestID *string
		status           string
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
				ctx:              ctx,
				staffID:          &staffID,
				serviceRequestID: &serviceRequestID,
				status:           enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case: invalid staff id",
			args: args{
				ctx:              ctx,
				staffID:          &longWord,
				serviceRequestID: &serviceRequestID,
				status:           enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: non-existent staff",
			args: args{
				ctx:              ctx,
				staffID:          &nonExistentUUID,
				serviceRequestID: &serviceRequestID,
				status:           enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid service request id",
			args: args{
				ctx:              ctx,
				staffID:          &staffID,
				serviceRequestID: &longWord,
				status:           enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: non existent service request",
			args: args{
				ctx:              ctx,
				staffID:          &staffID,
				serviceRequestID: &nonExistentUUID,
				status:           enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.ResolveServiceRequest(tt.args.ctx, tt.args.staffID, tt.args.serviceRequestID, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ResolveServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.ResolveServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_ResolveStaffServiceRequest(t *testing.T) {
	ctx := context.Background()
	fakeString := gofakeit.HipsterSentence(10)
	serviceRequestID := "26b20a42-cbb8-4553-aedb-c539602d04fc"
	badUID := "BadUID"

	type args struct {
		ctx                context.Context
		staffID            *string
		serviceRequestID   *string
		verificationStatus string
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
				ctx:                ctx,
				staffID:            &staffID,
				serviceRequestID:   &serviceRequestID,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case: invalid staff id",
			args: args{
				ctx:                ctx,
				staffID:            &fakeString,
				serviceRequestID:   &serviceRequestID,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: non-existent staff",
			args: args{
				ctx:                ctx,
				staffID:            &badUID,
				serviceRequestID:   &serviceRequestID,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid service request id",
			args: args{
				ctx:                ctx,
				staffID:            &staffID,
				serviceRequestID:   &fakeString,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: non existent service request",
			args: args{
				ctx:                ctx,
				staffID:            &staffID,
				serviceRequestID:   &badUID,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ResolveStaffServiceRequest(tt.args.ctx, tt.args.staffID, tt.args.serviceRequestID, tt.args.verificationStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ResolveStaffServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.ResolveStaffServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_AssignRoles(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx    context.Context
		userID string
		roles  []enums.UserRoleType
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
				ctx:    ctx,
				userID: userID,
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator, enums.UserRoleTypeContentManagement},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid: invalid user ID",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator, enums.UserRoleTypeContentManagement},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Invalid: invalid role",
			args: args{
				ctx:    ctx,
				userID: userID,
				roles:  []enums.UserRoleType{enums.UserRoleType("invalid"), enums.UserRoleTypeContentManagement},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.AssignRoles(tt.args.ctx, tt.args.userID, tt.args.roles)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.AssignRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.AssignRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_RevokeRoles(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx    context.Context
		userID string
		roles  []enums.UserRoleType
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
				ctx:    ctx,
				userID: userID,
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator, enums.UserRoleTypeContentManagement},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid: invalid user ID",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator, enums.UserRoleTypeContentManagement},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Invalid: invalid role",
			args: args{
				ctx:    ctx,
				userID: userID,
				roles:  []enums.UserRoleType{enums.UserRoleType("invalid"), enums.UserRoleTypeContentManagement},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.RevokeRoles(tt.args.ctx, tt.args.userID, tt.args.roles)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RevokeRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.RevokeRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_InvalidateScreeningToolResponse(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		clientID   string
		questionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:        ctx,
				clientID:   clientID,
				questionID: screeningToolsQuestionID,
			},
			wantErr: false,
		},
		{
			name: "Invalid: invalid client ID",
			args: args{
				ctx:        ctx,
				clientID:   uuid.New().String(),
				questionID: screeningToolsQuestionID,
			},
			wantErr: true,
		},
		{
			name: "Invalid: invalid question ID",
			args: args{
				ctx:        ctx,
				clientID:   clientID,
				questionID: "12345",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.InvalidateScreeningToolResponse(tt.args.ctx, tt.args.clientID, tt.args.questionID); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.InvalidateScreeningToolResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateServiceRequestsFromKenyaEMR(t *testing.T) {
	ctx := context.Background()
	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}

	serviceReq := &gorm.ClientServiceRequest{
		Base:           gorm.Base{},
		ID:             &serviceRequestID,
		Active:         true,
		RequestType:    "RED_FLAG",
		Request:        "VERY SAD",
		Status:         "IN PROGRESS",
		InProgressAt:   &time.Time{},
		ResolvedAt:     &time.Time{},
		ClientID:       clientID,
		InProgressByID: &staffID,
		OrganisationID: uuid.New().String(),
		ResolvedByID:   &staffID,
		FacilityID:     facilityID,
		Meta:           `{}`,
	}

	badServiceRequestID := "badServiceRequestID"
	invalidServiceReq := &gorm.ClientServiceRequest{
		ID: &badServiceRequestID,
	}

	err = pg.DB.Create(serviceReq).Error
	if err != nil {
		t.Errorf("Create securityQuestionResponse failed: %v", err)
		return
	}

	type args struct {
		ctx     context.Context
		payload []*gorm.ClientServiceRequest
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
				ctx:     ctx,
				payload: []*gorm.ClientServiceRequest{serviceReq},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				payload: []*gorm.ClientServiceRequest{invalidServiceReq},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.UpdateServiceRequests(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UpdateServiceRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_UpdateAppointment(t *testing.T) {

	type args struct {
		ctx        context.Context
		payload    *gorm.Appointment
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: update an appointment using id",
			args: args{
				ctx: context.Background(),
				payload: &gorm.Appointment{
					ID: appointmentID,
				},
				updateData: map[string]interface{}{
					"appointment_type": "Dental",
					"status":           enums.AppointmentStatusCompleted.String(),
					"client_id":        clientID,
					"reason":           "Knocked up",
					"date":             time.Now().Add(time.Duration(100)),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: update non-existent appointment",
			args: args{
				ctx: context.Background(),
				payload: &gorm.Appointment{
					ID: gofakeit.UUID(),
				},
				updateData: map[string]interface{}{
					"appointment_type": "Dental",
					"status":           enums.AppointmentStatusCompleted.String(),
					"client_id":        clientID,
					"reason":           "Knocked up",
					"date":             time.Now().Add(time.Duration(100)),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: update appointment missing ids",
			args: args{
				ctx:     context.Background(),
				payload: &gorm.Appointment{},
				updateData: map[string]interface{}{
					"appointment_type": "Dental",
					"status":           enums.AppointmentStatusCompleted.String(),
					"client_id":        clientID,
					"reason":           "Knocked up",
					"date":             time.Now().Add(time.Duration(100)),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.UpdateAppointment(tt.args.ctx, tt.args.payload, tt.args.updateData)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateAppointment() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("PGInstance.UpdateAppointment() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateUserPinUpdateRequiredStatus(t *testing.T) {
	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
		status  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: update user pin update required status",
			args: args{
				ctx:     context.Background(),
				userID:  userID2,
				flavour: feedlib.FlavourConsumer,
				status:  true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserPinUpdateRequiredStatus(tt.args.ctx, tt.args.userID, tt.args.flavour, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserPinUpdateRequiredStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateClient(t *testing.T) {
	type args struct {
		ctx     context.Context
		client  *gorm.Client
		updates map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.Client
		wantErr bool
	}{
		{
			name: "Happy case: update client profile",
			args: args{
				ctx: context.Background(),
				client: &gorm.Client{
					ID: &clientID,
				},
				updates: map[string]interface{}{
					"fhir_patient_id": gofakeit.UUID(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: update client missing ID",
			args: args{
				ctx:    context.Background(),
				client: &gorm.Client{},
				updates: map[string]interface{}{
					"fhir_patient_id": gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: update client invalid field",
			args: args{
				ctx:    context.Background(),
				client: &gorm.Client{},
				updates: map[string]interface{}{
					"invalid_field_id": gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.UpdateClient(tt.args.ctx, tt.args.client, tt.args.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected client to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil && got.FHIRPatientID == nil {
				t.Errorf("expected client not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_UpdateHealthDiary(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		payload *gorm.ClientHealthDiaryEntry
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
				ctx: ctx,
				payload: &gorm.ClientHealthDiaryEntry{
					ClientHealthDiaryEntryID: &clientsHealthDiaryEntryID,
					ClientID:                 clientID,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				payload: &gorm.ClientHealthDiaryEntry{
					ClientHealthDiaryEntryID: &clientsHealthDiaryEntryID,
					ClientID:                 "clientID",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.UpdateHealthDiary(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateHealthDiary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UpdateHealthDiary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_UpdateUserPinChangeRequiredStatus(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
		status  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				userID:  userID2,
				flavour: "CONSUMER",
				status:  true,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				userID:  "userID2",
				flavour: "CONSUMER",
				status:  true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserPinChangeRequiredStatus(tt.args.ctx, tt.args.userID, tt.args.flavour, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserPinChangeRequiredStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateFailedSecurityQuestionsAnsweringAttempts(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx       context.Context
		userID    string
		failCount int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: reset failed security attempts",
			args: args{
				ctx:       ctx,
				userID:    userFailedSecurityCountID,
				failCount: 0,
			},
			wantErr: false,
		},
		{
			name: "Sad case: user not found",
			args: args{
				ctx:    ctx,
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid user ID",
			args: args{
				ctx:    ctx,
				userID: "32354",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateFailedSecurityQuestionsAnsweringAttempts(tt.args.ctx, tt.args.userID, tt.args.failCount); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateFailedSecurityQuestionsAnsweringAttempts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateUser(t *testing.T) {
	ctx := context.Background()

	invalidUserID := "invalid user"

	type args struct {
		ctx        context.Context
		user       *gorm.User
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				user: &gorm.User{
					UserID: &userID,
				},
				updateData: map[string]interface{}{
					"next_allowed_login": time.Now(),
					"failed_login_count": 0,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				user: &gorm.User{
					UserID: &invalidUserID,
				},
				updateData: map[string]interface{}{
					"next_allowed_login": time.Now(),
					"failed_login_count": 0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUser(tt.args.ctx, tt.args.user, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateUserActiveStatus(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
		active  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				userID:  userID,
				flavour: feedlib.FlavourConsumer,
				active:  true,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				userID:  "userID",
				flavour: feedlib.FlavourConsumer,
				active:  true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserActiveStatus(tt.args.ctx, tt.args.userID, tt.args.flavour, tt.args.active); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserActiveStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateFacility(t *testing.T) {
	ctx := context.Background()

	invalidFacilityID := "invalid facility"

	type args struct {
		ctx        context.Context
		facility   *gorm.Facility
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				facility: &gorm.Facility{
					FacilityID: &facilityID,
				},
				updateData: map[string]interface{}{
					"fhir_organization_id": uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				facility: &gorm.Facility{
					FacilityID: &invalidFacilityID,
				},
				updateData: map[string]interface{}{
					"fhir_organization_id": uuid.New().String(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateFacility(tt.args.ctx, tt.args.facility, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateFacility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
