package ussd_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	CRMDomain "gitlab.slade360emr.com/go/commontools/crm/pkg/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestImpl_AddAITSessionDetailsUnittest(t *testing.T) {

	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()
	if err != nil {
		t.Errorf("failed to initialize test service")
		return
	}

	validSessionId := uuid.New().String()
	phoneNumber := "+254707756919"
	level := 1
	text := "1*gabriel*were"

	validUSSDLeadDetails := &domain.USSDLeadDetails{
		ID:          "0",
		SessionID:   validSessionId,
		PhoneNumber: phoneNumber,
		Level:       level,
	}

	validData := &dto.SessionDetails{
		SessionID:   validSessionId,
		PhoneNumber: &phoneNumber,
		Level:       level,
		Text:        text,
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
		//test cases.
		{
			name: "successful_persist_data",
			args: args{
				ctx:   ctx,
				input: validData,
			},
			wantErr: false,
			want:    validUSSDLeadDetails,
		},
		{
			name: "failed_persist_data",
			args: args{
				ctx:   ctx,
				input: validData,
			},
			wantErr: true,
		},
		{
			name: "failed_msisdn_normalize",
			args: args{
				ctx:   ctx,
				input: validData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "successful_persist_data" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return &phoneNumber, nil
				}

				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) (*domain.USSDLeadDetails, error) {
					return validUSSDLeadDetails, nil
				}
			}

			if tt.name == "failed_persist_data" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return &phoneNumber, nil
				}

				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) (*domain.USSDLeadDetails, error) {
					return nil, fmt.Errorf("failed to add session details")
				}

			}
			if tt.name == "failed_msisdn_normalize" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("failed to normalize msisdn number")
				}

			}

			got, err := u.AITUSSD.AddAITSessionDetails(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.AddAITSessionDetails() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Impl.AddAITSessionDetails() = %v,want %v\n", got, tt.want)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected an error but did not get one")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("Did not expect an error but got one %v\n", err)
					return
				}
			}

		})
	}
}

func TestImpl_GetOrCreateSessionStateUnittest(t *testing.T) {

	ctx := context.Background()
	u, err := InitializeFakeUSSDTestService()

	if err != nil {
		t.Errorf("failed to initialize test service")
		return
	}
	validSessionId := uuid.New().String()
	phoneNumber := "+254707756919"
	level := 0
	text := "1*gabriel*were"

	validUSSDLeadDetails := &domain.USSDLeadDetails{
		ID:          "0",
		SessionID:   validSessionId,
		PhoneNumber: phoneNumber,
		Level:       level,
	}

	validSessionDetails := &dto.SessionDetails{
		SessionID:   validSessionId,
		PhoneNumber: &phoneNumber,
		Level:       level,
		Text:        text,
	}

	type args struct {
		ctx     context.Context
		payload *dto.SessionDetails
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.USSDLeadDetails
		wantErr bool
	}{
		//test cases.
		{
			name: "successful_return_session",
			args: args{
				ctx:     ctx,
				payload: validSessionDetails,
			},
			wantErr: false,
			want:    validUSSDLeadDetails,
		},
		{
			name: "failed_return_session",
			args: args{
				ctx:     ctx,
				payload: validSessionDetails,
			},
			wantErr: true,
		},
		{
			name: "successful_set_session",
			args: args{
				ctx:     ctx,
				payload: validSessionDetails,
			},
			want: &domain.USSDLeadDetails{
				Level:       0,
				ID:          "0",
				SessionID:   validSessionId,
				PhoneNumber: phoneNumber,
			},
			wantErr: false,
		},
		{
			name: "failed_set_session",
			args: args{
				ctx:     ctx,
				payload: validSessionDetails,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "successful_return_session" {
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (*domain.USSDLeadDetails, error) {
					return validUSSDLeadDetails, nil
				}
			}

			if tt.name == "failed_return_session" {
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (*domain.USSDLeadDetails, error) {
					return nil, fmt.Errorf("failed to get session details")
				}
			}

			if tt.name == "successful_set_session" {
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (*domain.USSDLeadDetails, error) {
					return nil, nil
				}
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) (*domain.USSDLeadDetails, error) {
					return &domain.USSDLeadDetails{
						Level:       0,
						ID:          "0",
						SessionID:   validSessionId,
						PhoneNumber: phoneNumber,
					}, nil
				}
			}

			if tt.name == "failed_set_session" {
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (*domain.USSDLeadDetails, error) {
					return nil, nil
				}

				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) (*domain.USSDLeadDetails, error) {
					return nil, fmt.Errorf("failed to add session details")
				}
			}

			got, err := u.AITUSSD.GetOrCreateSessionState(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.GetOrCreateSessionState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Impl.GetOrCreateSessionState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImpl_GetOrCreatePhoneNumberUserUnittest(t *testing.T) {

	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()
	if err != nil {
		t.Errorf("failed to initialize test service")
		return
	}

	UID := uuid.New().String()
	displayName := gofakeit.Name()
	email := gofakeit.Email()
	phoneNumber := "+254702215783"
	photoURL := uuid.New().String()
	providerId := uuid.New().String()

	createdUserResponse := dto.CreatedUserResponse{
		UID:         UID,
		DisplayName: displayName,
		Email:       email,
		PhoneNumber: phoneNumber,
		PhotoURL:    photoURL,
		ProviderID:  providerId,
	}

	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.CreatedUserResponse
		wantErr bool
	}{
		//test cases.
		{
			name: "successful_get_or_create",
			args: args{
				ctx:   ctx,
				phone: phoneNumber,
			},
			wantErr: false,
			want:    &createdUserResponse,
		},
		{
			name: "failed_get_or_create",
			args: args{
				ctx:   ctx,
				phone: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "successful_get_or_create" {
				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return &createdUserResponse, nil
				}
			}

			if tt.name == "failed_get_or_create" {
				fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
					return nil, fmt.Errorf("failed to get or create user")
				}
			}

			got, err := u.AITUSSD.GetOrCreatePhoneNumberUser(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.GetOrCreatePhoneNumberUser() error = %v, wantErr %v", err, tt.wantErr)

			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Impl.GetOrCreatePhoneNumberUser() = %v, want %v", got, tt.want)
			}
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected an error but did not get one\n")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("Did not expect an error but we got one\n")
					return
				}
			}
		})
	}
}

func TestImpl_CreateUserProfileUnittest(t *testing.T) {

	ctx := context.Background()
	u, err := InitializeFakeUSSDTestService()

	if err != nil {
		t.Errorf("failed to initialize test service")
		return
	}

	profileID := uuid.New().String()
	username := gofakeit.Name()
	phoneNumber := "+254702215783"
	termsAccepted := true
	suspended := false
	uid := uuid.New().String()

	userProfile := base.UserProfile{
		ID:            profileID,
		UserName:      &username,
		PrimaryPhone:  &phoneNumber,
		TermsAccepted: termsAccepted,
		Suspended:     suspended,
		VerifiedIdentifiers: []base.VerifiedIdentifier{
			{
				UID:           uid,
				LoginProvider: base.LoginProviderTypePhone,
				Timestamp:     time.Now().In(base.TimeLocation),
			},
		},
		VerifiedUIDS: []string{uid},
	}

	type args struct {
		ctx         context.Context
		phoneNumber string
		uid         string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.UserProfile
		wantErr bool
	}{
		//test cases.
		{
			name: "success_create_user_profile",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				uid:         uid,
			},
			wantErr: false,
			want:    &userProfile,
		},
		{
			name: "failed_create_user_profile",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
				uid:         uid,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "success_create_user_profile" {
				fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {
					return &userProfile, nil
				}
			}

			if tt.name == "failed_create_user_profile" {
				fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to create user profile")
				}
			}

			got, err := u.AITUSSD.CreateUserProfile(tt.args.ctx, tt.args.phoneNumber, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.CreateUserProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Impl.CreateUserProfile() = %v, want %v", got, tt.want)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected an error but did not get one")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("Did not expect an error but we got one")
					return
				}
			}
		})
	}
}

func TestImpl_CreateEmptyCustomerProfileUnittest(t *testing.T) {

	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()
	if err != nil {
		t.Errorf("failed to initialize test service")
		return
	}

	id := uuid.New().String()
	profileId := uuid.New().String()

	customer := base.Customer{
		ID:        id,
		ProfileID: &profileId,
	}

	type args struct {
		ctx       context.Context
		profileID string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.Customer
		wantErr bool
	}{
		//test cases.
		{
			name: "success_create_empty_customer_profile",
			args: args{
				ctx:       ctx,
				profileID: profileId,
			},
			want:    &customer,
			wantErr: false,
		},
		{
			name: "failed_create_empty_customer_profile",
			args: args{
				ctx:       ctx,
				profileID: profileId,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "success_create_empty_customer_profile" {
				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					return &customer, nil
				}
			}
			if tt.name == "failed_create_empty_customer_profile" {
				fakeRepo.CreateEmptyCustomerProfileFn = func(ctx context.Context, profileID string) (*base.Customer, error) {
					return nil, fmt.Errorf("failed to create empy customer profile")
				}
			}

			got, err := u.AITUSSD.CreateEmptyCustomerProfile(tt.args.ctx, tt.args.profileID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.CreateEmptyCustomerProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Impl.CreateEmptyCustomerProfile() = %v, want %v", got, tt.want)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected an error but did not get one\n")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("Did not expect an error but we got one\n")
					return
				}
			}
		})
	}
}

func TestImpl_StageCRMPayloadUnittest(t *testing.T) {

	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()

	if err != nil {
		t.Errorf("failed to initialize test service")
		return
	}

	contactType := "phone"
	contactValue := "+254702215783"
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	dateOfBirth := base.Date{
		Day:   0,
		Month: 0,
		Year:  0,
	}
	isSync := false
	timeSync := time.Now()
	optOut := "NO"
	wantCover := false
	contactChannel := "USSD"
	isRegistered := false

	contactLeadInput := dto.ContactLeadInput{
		ContactType:    contactType,
		ContactValue:   contactValue,
		FirstName:      firstName,
		LastName:       lastName,
		DateOfBirth:    dateOfBirth,
		IsSync:         isSync,
		TimeSync:       &timeSync,
		OptOut:         CRMDomain.GeneralOptionType(optOut),
		WantCover:      wantCover,
		ContactChannel: contactChannel,
		IsRegistered:   isRegistered,
	}

	type args struct {
		ctx     context.Context
		payload *dto.ContactLeadInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		//test cases.
		{
			name: "happy case",
			args: args{
				ctx:     ctx,
				payload: &contactLeadInput,
			},
			wantErr: false,
		},
		{
			name: "sad case",
			args: args{
				ctx:     ctx,
				payload: &contactLeadInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "happy case" {
				fakeRepo.StageCRMPayloadFn = func(ctx context.Context, payload *dto.ContactLeadInput) error {
					return nil
				}
			}

			if tt.name == "sad case" {
				fakeRepo.StageCRMPayloadFn = func(ctx context.Context, payload *dto.ContactLeadInput) error {
					return fmt.Errorf("failed to stage crm payload")
				}
			}

			err := u.AITUSSD.StageCRMPayload(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.StageCRMPayload() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected an error but did not get one\n")
					return
				}
			}

			if !tt.wantErr {
				if err != nil {
					t.Errorf("Did not expect an error but we got one\n")
					return
				}
			}
		})
	}
}

func TestImpl_UpdateSessionLevel_UnitTest(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()

	if err != nil {
		t.Errorf("failed to initialize test service")
		return
	}

	type args struct {
		ctx       context.Context
		level     int
		sessionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       ctx,
				level:     50,
				sessionID: "f44496b5-4f73-48f8-9f59-0ab79d3d571b",
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				level:     50,
				sessionID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (*domain.USSDLeadDetails, error) {
					return &domain.USSDLeadDetails{}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (*domain.USSDLeadDetails, error) {
					return nil, err
				}
			}

			if err := u.AITUSSD.UpdateSessionLevel(tt.args.ctx, tt.args.level, tt.args.sessionID); (err != nil) != tt.wantErr {
				t.Errorf("Impl.UpdateSessionLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImpl_UpdateOptOutCRMPayload(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()

	if err != nil {
		t.Errorf("failed to initialize test service")
		return
	}

	phoneNumber := "+254700100200"
	ContactType := "phone"
	ContactValue := phoneNumber
	FirstName := gofakeit.FirstName()
	LastName := gofakeit.LastName()
	DateOfBirth := base.Date{
		Day:   0,
		Month: 0,
		Year:  0,
	}
	IsSync := false
	TimeSync := time.Now()
	OptOut := "NO"
	WantCover := false
	ContactChannel := "USSD"
	IsRegistered := false

	contactLeadPayload := &dto.ContactLeadInput{
		ContactType:    ContactType,
		ContactValue:   ContactValue,
		FirstName:      FirstName,
		LastName:       LastName,
		DateOfBirth:    DateOfBirth,
		IsSync:         IsSync,
		TimeSync:       &TimeSync,
		OptOut:         CRMDomain.GeneralOptionType(OptOut),
		WantCover:      WantCover,
		ContactChannel: ContactChannel,
		IsRegistered:   IsRegistered,
	}

	type args struct {
		ctx         context.Context
		phoneNumber string
		contactLead *dto.ContactLeadInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:         ctx,
				phoneNumber: phoneNumber,
				contactLead: contactLeadPayload,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:         ctx,
				phoneNumber: "",
				contactLead: contactLeadPayload,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeRepo.UpdateOptOutCRMPayloadFn = func(ctx context.Context, phoneNumber string, contactLead *dto.ContactLeadInput) error {
					return nil
				}
			}

			if tt.name == "Sad case" {
				fakeRepo.UpdateOptOutCRMPayloadFn = func(ctx context.Context, phoneNumber string, contactLead *dto.ContactLeadInput) error {
					return err
				}
			}
			if err := u.AITUSSD.UpdateOptOutCRMPayload(tt.args.ctx, tt.args.phoneNumber, tt.args.contactLead); (err != nil) != tt.wantErr {
				t.Errorf("Impl.UpdateOptOutCRMPayload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImpl_UpdateBioData_Unittest(t *testing.T) {
	ctx := context.Background()

	u, err := InitializeFakeUSSDTestService()

	if err != nil {
		t.Errorf("failed to initialize test service")
		return
	}
	firstname := gofakeit.FirstName()
	lastname := gofakeit.LastName()

	biodata := &base.BioData{
		FirstName:   &firstname,
		LastName:    &lastname,
		DateOfBirth: &base.Date{},
		Gender:      "Male",
	}

	type args struct {
		ctx  context.Context
		id   string
		data base.BioData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:  ctx,
				id:   "a51b1767-ee98-40d7-bc96-24e6f1c3e8b6",
				data: *biodata,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:  ctx,
				id:   "",
				data: *biodata,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeRepo.UpdateBioDataFn = func(ctx context.Context, id string, data base.BioData) error {
					return nil
				}
			}

			if tt.name == "Sad case" {
				fakeRepo.UpdateBioDataFn = func(ctx context.Context, id string, data base.BioData) error {
					return err
				}
			}
			if err := u.AITUSSD.UpdateBioData(tt.args.ctx, tt.args.id, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Impl.UpdateBioData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
