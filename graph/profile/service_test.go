package profile

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/authorization/graph/authorization"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/mailgun/graph/mailgun"
	"gitlab.slade360emr.com/go/otp/graph/otp"
)

func GetFirestoreClient(t *testing.T) *firestore.Client {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	assert.Nil(t, err)

	ctx := base.GetAuthenticatedContext(t)
	firestoreClient, err := firebaseApp.Firestore(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, firestoreClient)
	return firestoreClient
}

func GetFirebaseAuthClient(t *testing.T) (*auth.Client, error) {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Firebase app: %w", err)
	}
	ctx := base.GetAuthenticatedContext(t)
	client, err := firebaseApp.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Firebase auth client: %w", err)
	}
	return client, nil
}

func TestNewService(t *testing.T) {
	service := NewService()
	service.checkPreconditions() // should not panic
}

func TestService_profileUpdates(t *testing.T) {
	ctx := context.Background()

	user, userErr := base.GetOrCreateFirebaseUser(ctx, base.TestUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	idTokens, idErr := base.AuthenticateCustomFirebaseToken(customToken)
	assert.Nil(t, idErr)
	assert.NotNil(t, idTokens)

	bearerToken := idTokens.IDToken
	authToken, err := base.ValidateBearerToken(ctx, bearerToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)

	authenticatedContext := context.WithValue(ctx, base.AuthTokenContextKey, authToken)

	bs, err := ioutil.ReadFile("testdata/photo.jpg")
	assert.Nil(t, err)
	photoBase64 := base64.StdEncoding.EncodeToString(bs)
	email := []string{"kithinjimkevin@gmail.com"}
	msisdn := "+254723002959"
	otpService := otp.NewService()
	otp, err := otpService.GenerateAndSendOTP(msisdn)
	assert.Nil(t, err)
	assert.NotZero(t, otp)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "bad_case",
			args: args{
				ctx: context.Background(), // no uid in this one
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "good_case",
			args: args{
				ctx: authenticatedContext, // should
			},
			want:    user.UID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := authorization.GetLoggedInUserUID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("authorization.GetLoggedInUserUID error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("authorization.GetLoggedInUserUID = %v, want %v", got, tt.want)
			}
			if got == tt.want && err == nil {
				profileSnapshot, err := s.RetrieveUserProfileFirebaseDocSnapshot(tt.args.ctx)
				assert.Nil(t, err)
				assert.NotNil(t, profileSnapshot)

				userprofile, err := s.UserProfile(tt.args.ctx)
				assert.Nil(t, err)
				assert.NotNil(t, userprofile)

				accepted, err := s.AcceptTermsAndConditions(tt.args.ctx, true)
				assert.Nil(t, err)
				assert.True(t, accepted)

				// Update the user profile
				userProfileInput := UserProfileInput{
					PhotoBase64:      photoBase64,
					PhotoContentType: base.ContentTypeJpg,
					Msisdns: []*UserProfilePhone{
						{Phone: msisdn, Otp: otp},
					},
					Emails:                     email,
					CanExperiment:              false,
					AskAgainToSetIsTester:      false,
					AskAgainToSetCanExperiment: false,
				}
				updatedProfile, err := s.UpdateUserProfile(
					tt.args.ctx, userProfileInput)
				assert.Nil(t, err)
				assert.NotNil(t, updatedProfile)

				practitionerSignupInput := PractitionerSignupInput{
					License:   "fake license",
					Cadre:     PractitionerCadreDoctor,
					Specialty: base.PractitionerSpecialtyAnaesthesia,
				}
				signedUp, err := s.PractitionerSignUp(
					tt.args.ctx, practitionerSignupInput)
				assert.Nil(t, err)
				assert.True(t, signedUp)
			}
		})
	}
}

func TestAppleTesterPractitionerSignup(t *testing.T) {
	ctx := context.Background()

	user, userErr := base.GetOrCreateFirebaseUser(ctx, base.TestUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	idTokens, idErr := base.AuthenticateCustomFirebaseToken(customToken)
	assert.Nil(t, idErr)
	assert.NotNil(t, idTokens)

	bearerToken := idTokens.IDToken
	authToken, err := base.ValidateBearerToken(ctx, bearerToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)

	authenticatedContext := context.WithValue(ctx, base.AuthTokenContextKey, authToken)

	s := NewService()

	practitionerSignupInput := PractitionerSignupInput{
		License:   appleTesterPractitionerLicense,
		Cadre:     PractitionerCadreDoctor,
		Specialty: base.PractitionerSpecialtyAnaesthesia,
	}

	signedUp, err := s.PractitionerSignUp(
		authenticatedContext, practitionerSignupInput)
	assert.Nil(t, err)
	assert.True(t, signedUp)

	profileSnapshot, err := s.RetrieveUserProfileFirebaseDocSnapshot(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, profileSnapshot)

	userprofile, err := s.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, userprofile)
	assert.NotEqual(t, true, userprofile.PractitionerApproved)

}

func TestService_RegisterPushToken(t *testing.T) {
	ctx := context.Background()

	user, userErr := base.GetOrCreateFirebaseUser(ctx, base.TestUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	idTokens, idErr := base.AuthenticateCustomFirebaseToken(customToken)
	assert.Nil(t, idErr)
	assert.NotNil(t, idTokens)

	bearerToken := idTokens.IDToken
	authToken, err := base.ValidateBearerToken(ctx, bearerToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)

	authenticatedContext := context.WithValue(ctx, base.AuthTokenContextKey, authToken)

	msisdn := "+254723002959"
	otpService := otp.NewService()
	otp, err := otpService.GenerateAndSendOTP(msisdn)
	assert.Nil(t, err)
	assert.NotZero(t, otp)

	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "good_case",
			args: args{
				ctx:   authenticatedContext, // should
				token: "an example push token for testing",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.RegisterPushToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RegisterPushToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.RegisterPushToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_CompleteSignup(t *testing.T) {
	ctx := context.Background()

	user, userErr := base.GetOrCreateFirebaseUser(ctx, base.TestUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	idTokens, idErr := base.AuthenticateCustomFirebaseToken(customToken)
	assert.Nil(t, idErr)
	assert.NotNil(t, idTokens)

	bearerToken := idTokens.IDToken
	authToken, err := base.ValidateBearerToken(ctx, bearerToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)

	authenticatedContext := context.WithValue(ctx, base.AuthTokenContextKey, authToken)
	expectedBonus := base.Decimal(decimal.NewFromFloat(healthcashWelcomeBonusAmount))

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *base.Decimal
		wantErr bool
	}{
		{
			name: "good_case",
			args: args{
				ctx: authenticatedContext,
			},
			want:    &expectedBonus,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.CompleteSignup(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CompleteSignup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.CompleteSignup() = %v, want %v", got, tt.want)
			}
			if err == nil {
				profile, err := s.UserProfile(tt.args.ctx)
				assert.Nil(t, err)
				assert.True(t, profile.IsApproved)
			}
		})
	}
}

func TestService_RecordPostVisitSurvey(t *testing.T) {
	ctx := context.Background()

	user, userErr := base.GetOrCreateFirebaseUser(ctx, base.TestUserEmail)
	assert.Nil(t, userErr)
	assert.NotNil(t, user)

	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, user.UID)
	assert.Nil(t, tokenErr)
	assert.NotNil(t, customToken)

	idTokens, idErr := base.AuthenticateCustomFirebaseToken(customToken)
	assert.Nil(t, idErr)
	assert.NotNil(t, idTokens)

	bearerToken := idTokens.IDToken
	authToken, err := base.ValidateBearerToken(ctx, bearerToken)
	assert.Nil(t, err)
	assert.NotNil(t, authToken)

	authenticatedContext := context.WithValue(ctx, base.AuthTokenContextKey, authToken)

	type args struct {
		ctx   context.Context
		input PostVisitSurveyInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx: authenticatedContext,
				input: PostVisitSurveyInput{
					LikelyToRecommend: 0,
					Criticism:         "piece of crap",
					Suggestions:       "replace it all",
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.RecordPostVisitSurvey(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RecordPostVisitSurvey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.RecordPostVisitSurvey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ConfirmEmail(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx:   ctx,
				email: fmt.Sprintf("test-%s@healthcloud.co.ke", uuid.New()),
			},
			wantErr: false,
		},
		{
			name: "invalid emails",
			args: args{
				ctx:   ctx,
				email: "not a valid email",
			},
			wantErr: true,
		},
		{
			name: "user not logged in",
			args: args{
				ctx:   context.Background(),
				email: "test@healthcloud.co.ke",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.ConfirmEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ConfirmEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				assert.NotNil(t, got)
				profile, err := s.UserProfile(tt.args.ctx)
				assert.Nil(t, err)
				assert.NotNil(t, profile)
				assert.True(t, base.StringSliceContains(profile.Emails, tt.args.email))
			}
		})
	}
}

func Test_generatePractitionerSignupEmailTemplate(t *testing.T) {

	tests := []struct {
		name string
		want string
	}{
		{
			name: "Good case",
			want: practitionerSignupEmail,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generatePractitionerSignupEmailTemplate(); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("generatePractitionerSignupEmailTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SendPractitionerSignUpEmail(t *testing.T) {
	type args struct {
		ctx          context.Context
		emailaddress string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good case",
			args: args{
				ctx:          context.Background(),
				emailaddress: "ngure.nyaga@savannahinformatics.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			if err := s.SendPractitionerSignUpEmail(tt.args.ctx, tt.args.emailaddress); (err != nil) != tt.wantErr {
				t.Errorf("Service.SendPractitionerSignUpEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_PractitionerSignUp(t *testing.T) {

	type args struct {
		ctx   context.Context
		input PractitionerSignupInput
	}

	practitionerSignupInput := PractitionerSignupInput{
		License:   "fake license",
		Cadre:     PractitionerCadreDoctor,
		Specialty: base.PractitionerSpecialtyAnaesthesia,
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Good case",
			args: args{
				ctx:   base.GetAuthenticatedContext(t),
				input: practitionerSignupInput,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.PractitionerSignUp(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.PractitionerSignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.PractitionerSignUp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_UserProfile(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid profile, logged in user",
			args: args{
				ctx: base.GetAuthenticatedContext(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.UserProfile(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.UserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
				if got != nil {
					updated, err := s.ConfirmEmail(tt.args.ctx, base.TestUserEmail)
					assert.Nil(t, err)
					assert.NotNil(t, updated)
					assert.NotZero(t, updated.Emails)
					assert.True(t, base.StringSliceContains(updated.Emails, base.TestUserEmail))

					profile, err := s.UserProfile(tt.args.ctx)
					assert.Nil(t, err)
					assert.NotNil(t, profile)
					assert.NotZero(t, profile.Emails)
					assert.True(t, base.StringSliceContains(profile.Emails, base.TestUserEmail))
					assert.NotZero(t, profile.UID)
				}
			}
		})
	}
}

func Test_generatePractitionerWelcomeEmailTemplate(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Good case",
			want: practitionerWelcomeEmail,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generatePractitionerWelcomeEmailTemplate(); got != tt.want {
				t.Errorf("generatePractitionerWelcomeEmailTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SendPractitionerWelcomeEmail(t *testing.T) {
	type args struct {
		ctx          context.Context
		emailaddress string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good case",
			args: args{
				ctx:          context.Background(),
				emailaddress: "ngure.nyaga@savannahinformatics.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			if err := s.SendPractitionerWelcomeEmail(tt.args.ctx, tt.args.emailaddress); (err != nil) != tt.wantErr {
				t.Errorf("Service.SendPractitionerWelcomeEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_AddTester(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid test case",
			args: args{
				ctx:   ctx,
				email: gofakeit.Email(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid email",
			args: args{
				ctx:   ctx,
				email: "this is not an email",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.AddTester(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddTester() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.AddTester() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_RemoveTester(t *testing.T) {
	validTesterEmail := gofakeit.Email()
	srv := NewService()
	ctx := context.Background()
	added, err := srv.AddTester(ctx, validTesterEmail)
	assert.Nil(t, err)
	assert.True(t, added)

	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "tester that exists",
			args: args{
				ctx:   ctx,
				email: validTesterEmail,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "tester that does not exist",
			args: args{
				ctx:   ctx,
				email: fmt.Sprintf("%s@healthcloud.co.ke", uuid.New().String()),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.RemoveTester(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RemoveTester() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.RemoveTester() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ListTesters(t *testing.T) {
	validTesterEmail := gofakeit.Email()
	srv := NewService()
	ctx := context.Background()
	added, err := srv.AddTester(ctx, validTesterEmail)
	assert.Nil(t, err)
	assert.True(t, added)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.ListTesters(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListTesters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.GreaterOrEqual(t, len(got), 1)
		})
	}
}

func Test_generatePractitionerRejectionEmailTemplate(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Good case",
			want: practitionerSignupRejectionEmail,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generatePractitionerRejectionEmailTemplate(); got != tt.want {
				t.Errorf("generatePractitionerRejectionEmailTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SendPractitionerRejectionEmail(t *testing.T) {
	type args struct {
		ctx          context.Context
		emailaddress string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Good case",
			args: args{
				ctx:          context.Background(),
				emailaddress: "ngure.nyaga@savannahinformatics.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			if err := s.SendPractitionerRejectionEmail(tt.args.ctx, tt.args.emailaddress); (err != nil) != tt.wantErr {
				t.Errorf("Service.SendPractitionerRejectionEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ApprovePractitionerSignup(t *testing.T) {
	type args struct {
		ctx            context.Context
		practitionerID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Good case - Approve a profile",
			args: args{
				ctx:            base.GetAuthenticatedContext(t),
				practitionerID: "a7942fb4-61b4-4cf2-ab39-a2904d3090c3",
			},
			wantErr: false,
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.ApprovePractitionerSignup(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ApprovePractitionerSignup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.ApprovePractitionerSignup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_RejectPractitionerSignup(t *testing.T) {

	type args struct {
		ctx            context.Context
		practitionerID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Good case - Reject a profile",
			args: args{
				ctx:            base.GetAuthenticatedContext(t),
				practitionerID: "a7942fb4-61b4-4cf2-ab39-a2904d3090c3",
			},
			wantErr: false,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := s.RejectPractitionerSignup(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RejectPractitionerSignup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.RejectPractitionerSignup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetRegisteredPractitionerByLicense(t *testing.T) {
	firestoreClient := GetFirestoreClient(t)
	emailService := mailgun.NewService()
	firebaseAuth, _ := GetFirebaseAuthClient(t)
	type fields struct {
		firestoreClient *firestore.Client
		firebaseAuth    *auth.Client
		emailService    *mailgun.Service
	}
	type args struct {
		ctx     context.Context
		license string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - Retrieve a single practitioner records",
			fields: fields{
				firestoreClient: firestoreClient,
				firebaseAuth:    firebaseAuth,
				emailService:    emailService,
			},
			args: args{
				ctx:     base.GetAuthenticatedContext(t),
				license: "A8082",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				firestoreClient: tt.fields.firestoreClient,
				firebaseAuth:    tt.fields.firebaseAuth,
				emailService:    tt.fields.emailService,
			}
			_, err := s.GetRegisteredPractitionerByLicense(tt.args.ctx, tt.args.license)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetRegisteredPractitionerByLicense() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_ListKMPDURegisteredPractitioners(t *testing.T) {
	firestoreClient := GetFirestoreClient(t)
	emailService := mailgun.NewService()
	firebaseAuth, _ := GetFirebaseAuthClient(t)
	type fields struct {
		firestoreClient *firestore.Client
		firebaseAuth    *auth.Client
		emailService    *mailgun.Service
	}
	type args struct {
		ctx        context.Context
		pagination *base.PaginationInput
		filter     *base.FilterInput
		sort       *base.SortInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - Retreive all practitioner records",
			fields: fields{
				firestoreClient: firestoreClient,
				firebaseAuth:    firebaseAuth,
				emailService:    emailService,
			},
			args: args{
				ctx:        base.GetAuthenticatedContext(t),
				pagination: nil,
				filter:     nil,
				sort:       nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				firestoreClient: tt.fields.firestoreClient,
				firebaseAuth:    tt.fields.firebaseAuth,
				emailService:    tt.fields.emailService,
			}
			_, err := s.ListKMPDURegisteredPractitioners(tt.args.ctx, tt.args.pagination, tt.args.filter, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListKMPDURegisteredPractitioners() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
