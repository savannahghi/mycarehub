package profile

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"google.golang.org/api/iterator"
)

func deleteCollection(
	ctx context.Context,
	client *firestore.Client,
	ref *firestore.CollectionRef,
	batchSize int) error {
	for {
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0
		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}

func TestMain(m *testing.M) {
	log.Printf("Setting tests up ...")
	os.Setenv("ENVIRONMENT", "testing")
	os.Setenv("ROOT_COLLECTION_SUFFIX", "onboarding_testing")
	ctx := context.Background()
	s := NewService()

	log.Printf("Running tests ...")
	code := m.Run()

	log.Printf("Tearing tests down ...")
	collections := []string{
		s.GetPINCollectionName(),
		s.GetUserProfileCollectionName(),
		s.GetPractitionerCollectionName(),
	}
	for _, collection := range collections {
		ref := s.firestoreClient.Collection(collection)
		deleteCollection(ctx, s.firestoreClient, ref, 10)
	}

	os.Exit(code)
}

func TestNewService(t *testing.T) {
	service := NewService()
	service.checkPreconditions() // should not panic
}

func TestService_profileUpdates(t *testing.T) {
	ctx, token := base.GetAuthenticatedContextAndToken(t)
	bs, err := ioutil.ReadFile("testdata/photo.jpg")
	assert.Nil(t, err)
	photoBase64 := base64.StdEncoding.EncodeToString(bs)
	email := []string{gofakeit.Email()}

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
				ctx: ctx, // should
			},
			want:    token.UID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewService()
			got, err := base.GetLoggedInUserUID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("base.GetLoggedInUserUID error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("base.GetLoggedInUserUID = %v, want %v", got, tt.want)
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
					PhotoBase64:                photoBase64,
					PhotoContentType:           base.ContentTypeJpg,
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
				ctx:   base.GetAuthenticatedContext(t), // should
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

	type args struct {
		ctx context.Context
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
				ctx: authenticatedContext,
			},
			want:    true,
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
					assert.NotZero(t, profile.VerifiedIdentifiers)
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
	ctx := base.GetAuthenticatedContext(t)
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
	ctx := base.GetAuthenticatedContext(t)
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
	ctx := base.GetAuthenticatedContext(t)
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
	service := NewService()
	type args struct {
		ctx     context.Context
		license string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - Retrieve a single practitioner records",
			args: args{
				ctx:     base.GetAuthenticatedContext(t),
				license: "A8082",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			_, err := s.GetRegisteredPractitionerByLicense(tt.args.ctx, tt.args.license)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetRegisteredPractitionerByLicense() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_ListKMPDURegisteredPractitioners(t *testing.T) {
	service := NewService()
	type args struct {
		ctx        context.Context
		pagination *base.PaginationInput
		filter     *base.FilterInput
		sort       *base.SortInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case - Retreive all practitioner records",
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
			s := service
			_, err := s.ListKMPDURegisteredPractitioners(tt.args.ctx, tt.args.pagination, tt.args.filter, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListKMPDURegisteredPractitioners() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_IsUnderAge(t *testing.T) {
	service := NewService()
	ctx := base.GetAuthenticatedContext(t)
	profile, err := service.UserProfile(ctx)
	if err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}
	date := &base.Date{
		Year:  1997,
		Month: 12,
		Day:   13,
	}
	profile.DateOfBirth = date
	dsnap, err := service.RetrieveUserProfileFirebaseDocSnapshot(ctx)
	if err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}
	err = base.UpdateRecordOnFirestore(
		service.firestoreClient, service.GetUserProfileCollectionName(),
		dsnap.Ref.ID, profile,
	)
	if err != nil {
		t.Errorf("got %v, want %v", err, nil)
	}

	type args struct {
		ctx context.Context
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
				ctx: ctx,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.IsUnderAge(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.IsUnderAge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.IsUnderAge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SetUserPin(t *testing.T) {
	service := NewService()
	ctx := base.GetPhoneNumberAuthenticatedContext(t)
	type args struct {
		ctx    context.Context
		msisdn string
		pin    int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// expectation is creation of the new pin
		{
			name: "Happy case: successfully set a user pin",
			args: args{
				ctx:    ctx,
				msisdn: base.TestUserPhoneNumber,
				pin:    1234,
			},
			want:    true,
			wantErr: false,
		},
		// expectation is the return of the existing PIN
		// since they have created it on the first place
		{
			name: "ensure PIN is only one",
			args: args{
				ctx:    ctx,
				msisdn: base.TestUserPhoneNumber,
				pin:    5645,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			isPINSet, err := s.SetUserPIN(tt.args.ctx, tt.args.msisdn, tt.args.pin)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SetUserPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if isPINSet != tt.want {
				t.Errorf("Service.SetUserPIN() = %v, want %v", isPINSet, tt.want)
			}
			// retrieve the logged in user
			// they should have a PIN set
			profile, err := s.UserProfile(ctx)
			if err == nil {
				assert.True(t, profile.HasPin)
			}
		})
	}
}

func TestService_RetrievePINFirebaseDocSnapshotByMSISDN(t *testing.T) {
	service := NewService()
	ctx := base.GetPhoneNumberAuthenticatedContext(t)
	set, err := service.SetUserPIN(ctx, "+254703754685", 1234)
	if !set {
		t.Errorf("setting a pin for test user failed. It returned false")
	}
	if err != nil {
		t.Errorf("setting a pin for test user failed: %v", err)
	}
	type args struct {
		ctx    context.Context
		msisdn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: retreive a pin dsnap that exists",
			args: args{
				ctx:    ctx,
				msisdn: base.TestUserPhoneNumber,
			},
			wantErr: false,
		},
		{
			name: "Sad case: retreive a pin that does not exist",
			args: args{
				ctx:    ctx,
				msisdn: "ain't no such a number",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			_, err := s.RetrievePINFirebaseDocSnapshotByMSISDN(tt.args.ctx, tt.args.msisdn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RetrievePINFirebaseDocSnapshotByMSISDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_VerifyMSISDNandPin(t *testing.T) {
	service := NewService()
	type args struct {
		ctx    context.Context
		msisdn string
		pin    int
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
				ctx:    base.GetPhoneNumberAuthenticatedContext(t),
				msisdn: base.TestUserPhoneNumber,
				pin:    1234,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx:    base.GetPhoneNumberAuthenticatedContext(t),
				msisdn: "not even close to an msisdn",
				pin:    1256,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.VerifyMSISDNandPIN(tt.args.ctx, tt.args.msisdn, tt.args.pin)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.VerifyMSISDNandPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.VerifyMSISDNandPIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_CheckHasPin(t *testing.T) {
	service := NewService()
	ctx := base.GetPhoneNumberAuthenticatedContext(t)
	set, err := service.SetUserPIN(ctx, base.TestUserPhoneNumber, 1234)
	if !set {
		t.Errorf("setting a pin for test user failed. It returned false")
	}
	if err != nil {
		t.Errorf("setting a pin for test user failed: %v", err)
	}

	type args struct {
		ctx    context.Context
		msisdn string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: the user has a pin",
			args: args{
				ctx:    ctx,
				msisdn: base.TestUserPhoneNumber,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: data with a bad phone number",
			args: args{
				ctx:    ctx,
				msisdn: "not a valid phone number",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: the user does not have a pin",
			args: args{
				ctx:    context.Background(),
				msisdn: "+254712345678",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.CheckHasPIN(tt.args.ctx, tt.args.msisdn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CheckHasPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.CheckHasPIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_SetLanguagePreference(t *testing.T) {
	service := NewService()
	type args struct {
		ctx      context.Context
		language base.Language
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "An allowed language/enum type",
			args: args{
				ctx:      base.GetAuthenticatedContext(t),
				language: base.LanguageEn,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "An allowed language/enum type",
			args: args{
				ctx:      base.GetAuthenticatedContext(t),
				language: "not a language",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.SetLanguagePreference(tt.args.ctx, tt.args.language)
			if err == nil {
				assert.NotNil(t, got)
			}
			if err != nil {
				assert.Empty(t, got)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.SetLanguagePreference() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.SetLanguagePreference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_VerifyEmailOtp(t *testing.T) {
	service := NewService()
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	assert.Nil(t, err)

	ctx := base.GetAuthenticatedContext(t)
	firestoreClient, err := firebaseApp.Firestore(ctx)
	assert.Nil(t, err)

	validOtpCode := rand.Int()
	validOtpData := map[string]interface{}{
		"authorizationCode": strconv.Itoa(validOtpCode),
		"isValid":           true,
		"message":           "Testing email OTP message",
		"timestamp":         time.Now(),
		"email":             "ngure.nyaga@healthcloud.co.ke",
	}
	_, err = base.SaveDataToFirestore(firestoreClient,
		base.SuffixCollection(base.OTPCollectionName), validOtpData)

	assert.Nil(t, err)
	type args struct {
		ctx   context.Context
		email string
		otp   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case - sent otp code",
			args: args{
				ctx:   ctx,
				email: "ngure.nyaga@healthcloud.co.ke",
				otp:   strconv.Itoa(validOtpCode),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - non existent otp code",
			args: args{
				ctx:   ctx,
				email: "ngure.nyaga@healthcloud.co.ke",
				otp:   "029837",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.VerifyEmailOtp(tt.args.ctx, tt.args.email, tt.args.otp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.VerifyEmailOtp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.VerifyEmailOtp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_CreateSignUpMethod(t *testing.T) {
	service := NewService()
	type args struct {
		ctx          context.Context
		signUpMethod SignUpMethod
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
				ctx:          base.GetAuthenticatedContext(t),
				signUpMethod: "google",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid case",
			args: args{
				ctx:          base.GetAuthenticatedContext(t),
				signUpMethod: "not a sign up method",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "No logged in user case",
			args: args{
				ctx:          context.Background(),
				signUpMethod: "not a sign up method",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.CreateSignUpMethod(tt.args.ctx, tt.args.signUpMethod)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateSignUpMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.CreateSignUpMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_GetSignUpMethod(t *testing.T) {
	service := NewService()
	ctx, authToken := base.GetAuthenticatedContextAndToken(t)
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    SignUpMethod
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx: ctx,
				id:  authToken.UID,
			},
			want:    "google",
			wantErr: false,
		},
		{
			name: "sad case - sign up method not found",
			args: args{
				ctx: ctx,
				id:  "invalid uid",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.GetSignUpMethod(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetSignUpMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.GetSignUpMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_AddPractitionerServices(t *testing.T) {
	service := NewService()
	ctx, _ := base.GetAuthenticatedContextAndToken(t)

	type args struct {
		ctx           context.Context
		services      PractitionerServiceInput
		otherServices *OtherPractitionerServiceInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case - without other option",
			args: args{
				ctx: ctx,
				services: PractitionerServiceInput{
					Services: []PractitionerService{"PHARMACY"},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "happy case - with other option",
			args: args{
				ctx: base.GetAuthenticatedContext(t),
				services: PractitionerServiceInput{
					Services: []PractitionerService{"OUTPATIENT_SERVICES", "OTHER"},
				},
				otherServices: &OtherPractitionerServiceInput{
					OtherServices: []string{"other-services"},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case - invalid enums",
			args: args{
				ctx: base.GetAuthenticatedContext(t),
				services: PractitionerServiceInput{
					Services: []PractitionerService{"not a valid enum"},
				},
				otherServices: &OtherPractitionerServiceInput{
					OtherServices: []string{"other-services"},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case - others specified but no data entered",
			args: args{
				ctx: base.GetAuthenticatedContext(t),
				services: PractitionerServiceInput{
					Services: []PractitionerService{"OTHER"},
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.AddPractitionerServices(tt.args.ctx, tt.args.services, tt.args.otherServices)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.AddPractitionerServices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Service.AddPractitionerServices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_RetrieveFireStoreSnapshotByUID(t *testing.T) {
	service := NewService()
	ctx, token := base.GetAuthenticatedContextAndToken(t)

	type args struct {
		ctx            context.Context
		uid            string
		collectionName string
		field          string
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
				uid:            token.UID,
				collectionName: service.GetPractitionerCollectionName(),
				field:          "profile.verifiedIdentifiers",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			got, err := s.RetrieveFireStoreSnapshotByUID(tt.args.ctx, tt.args.uid, tt.args.collectionName, tt.args.field)
			if err == nil {
				assert.NotNil(t, got)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.RetrieveFireStoreSnapshotByUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_SaveMemberCoverToFirestore(t *testing.T) {

	ctx := base.GetAuthenticatedContext(t)
	assert.NotNil(t, ctx, "context is nil")

	srv := NewService()
	assert.NotNil(t, srv, "service is nil")

	type args struct {
		payerName      string
		memberNumber   string
		memberName     string
		PayerSladeCode int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid case",
			args: args{
				memberName:     "Jakaya",
				memberNumber:   "144",
				payerName:      "Jubilee",
				PayerSladeCode: 136,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := srv.SaveMemberCoverToFirestore(ctx, tt.args.payerName, tt.args.memberNumber, tt.args.memberName, tt.args.PayerSladeCode); (err != nil) != tt.wantErr {
				t.Errorf("Service.SaveMemberCoverToFirestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func createNewUser(ctx context.Context, t *testing.T) (context.Context, *auth.Token) {
	authClient, err := base.GetFirebaseAuthClient(ctx)
	if err != nil {
		return nil, nil
	}
	params := (&auth.UserToCreate{}).
		EmailVerified(false).
		Disabled(false)
	newUser, createErr := authClient.CreateUser(ctx, params)
	if createErr != nil {
		return nil, nil
	}

	customToken, tokenErr := base.CreateFirebaseCustomToken(ctx, newUser.UID)
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

	return authenticatedContext, authToken
}
func TestService_GetOrCreateUserProfile(t *testing.T) {
	service := NewService()
	ctx := context.Background()
	newCtx, newToken := createNewUser(ctx, t)
	if newToken == nil {
		t.Errorf("unable to create new user token")
		return
	}
	if newCtx == nil {
		t.Errorf("unable to create new user context")
		return
	}
	emailCtx, emailToken := base.GetAuthenticatedContextAndToken(t)
	if emailToken == nil {
		t.Errorf("unable to create new email token")
		return
	}
	if emailCtx == nil {
		t.Errorf("unable to create new email context")
		return
	}

	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Case 1: creating a new user profile",
			args: args{
				ctx:   emailCtx,
				phone: "+254716862585",
			},
			wantErr: false,
		},
		{
			name: "Case 1: Linking a new uid to the existing profile",
			args: args{
				ctx:   newCtx,
				phone: "+254716862585",
			},
			wantErr: false,
		},
		{
			name: "Bad Case: non existent user profile",
			args: args{
				ctx:   ctx,
				phone: "+254716862585",
			},
			wantErr: true,
		},
		{
			name: "Bad Case: bad phone nnumber",
			args: args{
				ctx:   ctx,
				phone: "not a valid phone number",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			profile, err := s.GetOrCreateUserProfile(tt.args.ctx, tt.args.phone)
			if err == nil && tt.wantErr == false {
				if profile == nil {
					t.Errorf("empty profile was found")
					return
				}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetOrCreateUserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_FindProfile(t *testing.T) {
	service := NewService()
	ctx := context.Background()
	emailCtx := base.GetAuthenticatedContext(t)
	if emailCtx == nil {
		t.Errorf("unable to create email user context")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sad case: the document does not exist",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "happy case: the document does exist",
			args: args{
				ctx: emailCtx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			_, err := s.FindProfile(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.FindProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_GetProfileByID(t *testing.T) {
	s := NewService()
	if s == nil {
		t.Errorf("nil profile service, can't proceed with tests")
		return
	}

	authenticatedContext := base.GetAuthenticatedContext(t)
	if authenticatedContext == nil {
		t.Errorf("can't initialize authenticated context for testing")
		return
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "authenticated context",
			args: args{
				ctx: authenticatedContext,
				id:  ksuid.New().String(),
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "unauthenticated context",
			args: args{
				ctx: context.Background(),
				id:  ksuid.New().String(),
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := s.GetProfileByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetProfileByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil && profile == nil {
				t.Errorf("got a nil profile, did not expect one")
				return
			}
			if !tt.wantNil && !tt.wantErr && profile != nil {
				profileID := profile.ID

				// fetch again with the same profile ID
				refetchedProfile, err := s.GetProfileByID(tt.args.ctx, profileID)
				if err != nil {
					t.Errorf("unable to re-fetch newly created profile by ID: %v", err)
					return
				}

				if refetchedProfile == nil {
					t.Errorf("newly created profile is nil when re-fetched")
					return
				}
			}
		})
	}
}
