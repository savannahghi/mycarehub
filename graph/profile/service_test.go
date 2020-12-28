package profile

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"firebase.google.com/go/auth"
	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestMain(m *testing.M) {
	log.Printf("Setting tests up ...")
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("DEBUG", "true")
	os.Setenv("ROOT_COLLECTION_SUFFIX", fmt.Sprintf("profile_ci_%v", time.Now().Unix()))
	ctx := context.Background()
	s := NewService()

	purgeRecords := func() {
		collections := []string{
			s.GetPINCollectionName(),
			s.GetUserProfileCollectionName(),
			s.GetPractitionerCollectionName(),
			s.GetSignUpInfoCollectionName(),
			s.GetSupplierCollectionName(),
			s.GetSurveyCollectionName(),
			s.GetProfileNudgesCollectionName(),
			base.GetCollectionName(&TesterWhitelist{}),
		}
		for _, collection := range collections {
			ref := s.firestoreClient.Collection(collection)
			base.DeleteCollection(ctx, s.firestoreClient, ref, 10)
		}
	}
	purgeRecords()

	log.Printf("Running tests ...")
	code := m.Run()

	log.Printf("Tearing tests down ...")
	purgeRecords()

	os.Exit(code)
}

func TestNewService(t *testing.T) {
	service := NewService()
	service.checkPreconditions() // should not panic
}

func TestService_profileUpdates(t *testing.T) {
	ctx, token := base.GetAuthenticatedContextAndToken(t)
	bs, err := ioutil.ReadFile("testdata/photo.jpg")
	if err != nil {
		t.Errorf("unable to readfile: %V", err)
	}
	photoBase64 := base64.StdEncoding.EncodeToString(bs)
	email := []string{base.GenerateRandomEmail()}

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
				ctx: ctx,
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
				if err != nil {
					t.Errorf("unable to updateUserProfile: %v", err)
					return
				}
				if updatedProfile == nil {
					t.Errorf("nil updatedProfile")
				}

			}
		})
	}
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
					profile, err := s.UserProfile(tt.args.ctx)
					assert.Nil(t, err)
					assert.NotNil(t, profile)
					assert.NotZero(t, profile.VerifiedIdentifiers)
				}
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
				email: base.GenerateRandomEmail(),
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
	validTesterEmail := base.GenerateRandomEmail()
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
				email: base.GenerateRandomEmail(),
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
	validTesterEmail := base.GenerateRandomEmail()
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

func TestService_SetUserPin(t *testing.T) {
	service := NewService()
	ctx := base.GetPhoneNumberAuthenticatedContext(t)
	type args struct {
		ctx    context.Context
		msisdn string
		pin    string
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
				pin:    "1234",
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
				pin:    "5645",
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

func TestService_CheckHasPin(t *testing.T) {
	service := NewService()
	ctx := base.GetPhoneNumberAuthenticatedContext(t)
	set, err := service.SetUserPIN(ctx, base.TestUserPhoneNumber, "1234")
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
				msisdn: "+254711223344",
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

func TestService_RetrieveFireStoreSnapshotByUID(t *testing.T) {
	s := NewService()
	ctx, token := base.GetAuthenticatedContextAndToken(t)
	// ensure a user profile is created
	_, err := s.GetProfile(ctx, token.UID)
	assert.Nil(t, err)

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
				collectionName: s.GetUserProfileCollectionName(),
				field:          "verifiedIdentifiers",
			},
			wantErr: false,
		},

		{
			name: "non existent uid",
			args: args{
				ctx:            ctx,
				uid:            "122555",
				collectionName: s.GetUserProfileCollectionName(),
				field:          "verifiedIdentifiers",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.RetrieveFireStoreSnapshotByUID(
				tt.args.ctx,
				tt.args.uid,
				tt.args.collectionName,
				tt.args.field,
			)
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

func TestService_DeleteUser(t *testing.T) {
	service := NewService()
	ctx := context.Background()
	userCtx, token := createNewUser(ctx, t)

	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: delete a user",
			args: args{
				ctx: userCtx,
				uid: token.UID,
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to delete a user",
			args: args{
				ctx: context.Background(),
				uid: token.UID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service
			if err := s.DeleteUser(tt.args.ctx, tt.args.uid); (err != nil) != tt.wantErr {
				t.Errorf("Service.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_CreateUserByPhone(t *testing.T) {
	service := NewService()
	ctx := context.Background()

	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    *auth.UserRecord
		wantErr bool
	}{
		{
			name: "happy case: create a user",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254725120120",
			},
			wantErr: false,
		},
		{
			name: "create a user: invalid phone provided",
			args: args{
				ctx:         ctx,
				phoneNumber: "725120120000",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdUser, err := service.CreateUserByPhone(tt.args.ctx, tt.args.phoneNumber)
			if !tt.wantErr {
				assert.NotNil(t, createdUser.UserProfile)
				assert.NotNil(t, createdUser.CustomToken)
				assert.Equal(t, tt.args.phoneNumber, createdUser.UserProfile.Msisdns[0])
			}
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
		})
	}
}

func Test_validatePIN(t *testing.T) {
	type args struct {
		pin string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: pin legth 4",
			args: args{
				pin: "1234",
			},
			wantErr: false,
		},
		{
			name: "valid: pin length 5",
			args: args{
				pin: "12346",
			},
			wantErr: false,
		},
		{
			name: "valid: pin length 6",
			args: args{
				pin: "123465",
			},
			wantErr: false,
		},
		{
			name: "inavlid pin with letters",
			args: args{
				pin: "qwer",
			},
			wantErr: true,
		},
		{
			name: "invalid: pin with less digits",
			args: args{
				pin: "11",
			},
			wantErr: true,
		},
		{
			name: "invalid: empty pin",
			args: args{
				pin: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePIN(tt.args.pin); (err != nil) != tt.wantErr {
				t.Errorf("validatePIN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
