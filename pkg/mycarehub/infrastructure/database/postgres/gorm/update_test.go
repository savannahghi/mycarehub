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

func TestPGInstance_UpdateUserPinChangeRequiredStatus(t *testing.T) {
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

			got, err := testingDB.UpdateUserPinChangeRequiredStatus(tt.args.ctx, tt.args.userID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserPinChangeRequiredStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UpdateUserPinChangeRequiredStatus() = %v, want %v", got, tt.want)
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

func TestPGInstance_UpdateUserLastSuccessfulLoginTime(t *testing.T) {
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
				userID: userIDToUpdateUserLastSuccessfulLoginTime,
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
			if err := testingDB.UpdateUserLastSuccessfulLoginTime(tt.args.ctx, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserLastSuccessfulLoginTime() error = %v, wantErr %v", err, tt.wantErr)
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
			},
			wantErr: true,
		},
		{
			name: "Sad case: non-existent staff",
			args: args{
				ctx:              ctx,
				staffID:          &nonExistentUUID,
				serviceRequestID: &serviceRequestID,
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid service request id",
			args: args{
				ctx:              ctx,
				staffID:          &staffID,
				serviceRequestID: &longWord,
			},
			wantErr: true,
		},
		{
			name: "Sad case: non existent service request",
			args: args{
				ctx:              ctx,
				staffID:          &staffID,
				serviceRequestID: &nonExistentUUID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.ResolveServiceRequest(tt.args.ctx, tt.args.staffID, tt.args.serviceRequestID)
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
