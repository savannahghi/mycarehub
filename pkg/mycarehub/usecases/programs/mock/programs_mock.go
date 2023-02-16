package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// ProgramsUseCaseMock mocks the implementation of Programs usecase methods.
type ProgramsUseCaseMock struct {
	MockCreateProgramFn            func(ctx context.Context, input *dto.ProgramInput) (bool, error)
	MockListUserProgramsFn         func(ctx context.Context, userID string, flavour feedlib.Flavour) (*dto.ProgramOutput, error)
	MockSetCurrentProgramFn        func(ctx context.Context, programID string) (bool, error)
	MockGetProgramFacilitiesFn     func(ctx context.Context, programID string) ([]*domain.Facility, error)
	MockListProgramsFn             func(ctx context.Context, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error)
	MockSetStaffProgramFn          func(ctx context.Context, programID string) (*domain.StaffResponse, error)
	MockSetClientProgramFn         func(ctx context.Context, programID string) (*domain.ClientResponse, error)
	MockSearchProgramsFn           func(ctx context.Context, searchParameter string) ([]*domain.Program, error)
	MockGetProgramByIDFn           func(ctx context.Context, programID string) (*domain.Program, error)
	MockListOrganisationProgramsFn func(ctx context.Context, organisationID string, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error)
}

// NewProgramsUseCaseMock creates in itializes create type mocks
func NewProgramsUseCaseMock() *ProgramsUseCaseMock {
	UUID := gofakeit.UUID()
	now := time.Now()
	yesterday := now.Add(time.Hour * -24)
	organisation := domain.Organisation{
		ID:              UUID,
		Active:          false,
		Code:            gofakeit.BS(),
		Name:            gofakeit.BS(),
		Description:     gofakeit.BS(),
		EmailAddress:    gofakeit.Email(),
		PhoneNumber:     "0999999999",
		PostalAddress:   gofakeit.BS(),
		PhysicalAddress: gofakeit.BS(),
		DefaultCountry:  gofakeit.BS(),
	}
	program := domain.Program{
		ID:           UUID,
		Active:       true,
		Name:         gofakeit.BS(),
		Organisation: organisation,
	}
	programOutput := dto.ProgramOutput{
		Count:    1,
		Programs: []*domain.Program{&program},
	}

	facility := domain.Facility{
		ID:                 &UUID,
		Name:               gofakeit.BS(),
		Phone:              "0999999999",
		Active:             true,
		Country:            gofakeit.BS(),
		Description:        gofakeit.BS(),
		FHIROrganisationID: UUID,
		Identifier: domain.FacilityIdentifier{
			ID:     UUID,
			Active: true,
			Type:   enums.FacilityIdentifierTypeMFLCode,
			Value:  "32323232",
		},
		WorkStationDetails: domain.WorkStationDetails{
			Notifications:   1,
			Surveys:         1,
			Articles:        1,
			Messages:        1,
			ServiceRequests: 1,
		},
	}

	pagination := domain.Pagination{
		Limit:       1,
		CurrentPage: 1,
	}
	programPage := domain.ProgramPage{
		Pagination: pagination,
		Programs:   []*domain.Program{&program},
	}

	contact := domain.Contact{
		ID:             &UUID,
		ContactType:    "PHONE",
		ContactValue:   "0999999999",
		Active:         true,
		OptedIn:        true,
		UserID:         &UUID,
		OrganisationID: UUID,
	}

	role := domain.AuthorityRole{
		AuthorityRoleID: UUID,
		Name:            enums.UserRoleTypeCommunityManagement,
		Active:          true,
	}

	permission := domain.AuthorityPermission{
		PermissionID: UUID,
		Name:         enums.PermissionTypeCanCreateGroup,
		Active:       false,
	}
	user := domain.User{
		ID:                     new(string),
		Username:               gofakeit.BS(),
		Name:                   gofakeit.BS(),
		Gender:                 enumutils.GenderMale,
		Active:                 true,
		Contacts:               &contact,
		PushTokens:             []string{UUID},
		LastSuccessfulLogin:    &now,
		LastFailedLogin:        &yesterday,
		FailedLoginCount:       0,
		NextAllowedLogin:       &yesterday,
		PinChangeRequired:      false,
		HasSetPin:              true,
		HasSetSecurityQuestion: true,
		IsPhoneVerified:        true,
		TermsAccepted:          true,
		AcceptedTermsID:        1,
		Suspended:              false,
		Roles:                  []*domain.AuthorityRole{&role},
		Permissions:            []*domain.AuthorityPermission{&permission},
		DateOfBirth:            &yesterday,
		FailedSecurityCount:    0,
		PinUpdateRequired:      false,
		HasSetNickname:         true,
		CurrentOrganizationID:  UUID,
		CurrentProgramID:       UUID,
	}

	staffProfile := domain.StaffProfile{
		ID:              new(string),
		User:            &user,
		UserID:          UUID,
		Active:          true,
		StaffNumber:     "",
		Facilities:      []*domain.Facility{&facility},
		DefaultFacility: &facility,
		OrganisationID:  UUID,
		ProgramID:       UUID,
	}

	staffResponse := domain.StaffResponse{
		StaffProfile:   staffProfile,
		Roles:          []*domain.AuthorityRole{&role},
		Permissions:    []*domain.AuthorityPermission{&permission},
		CommunityToken: UUID,
	}

	clientProfile := domain.ClientProfile{
		ID:                      &UUID,
		User:                    &user,
		Active:                  true,
		ClientTypes:             []enums.ClientType{enums.ClientTypePmtct},
		UserID:                  UUID,
		TreatmentEnrollmentDate: &yesterday,
		FHIRPatientID:           &UUID,
		HealthRecordID:          &UUID,
		TreatmentBuddy:          UUID,
		ClientCounselled:        true,
		OrganisationID:          UUID,
		DefaultFacility:         &facility,
		CHVUserID:               &UUID,
		CHVUserName:             gofakeit.BS(),
		CaregiverID:             &UUID,
		CCCNumber:               "32323232",
		Facilities:              []*domain.Facility{&facility},
		ProgramID:               UUID,
	}

	clientResponse := &domain.ClientResponse{
		ClientProfile:  &clientProfile,
		Roles:          []*domain.AuthorityRole{&role},
		Permissions:    []*domain.AuthorityPermission{&permission},
		CommunityToken: UUID,
	}

	return &ProgramsUseCaseMock{
		MockCreateProgramFn: func(ctx context.Context, input *dto.ProgramInput) (bool, error) {
			return true, nil
		},
		MockListUserProgramsFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (*dto.ProgramOutput, error) {
			return &programOutput, nil
		},
		MockSetCurrentProgramFn: func(ctx context.Context, programID string) (bool, error) {
			return false, nil
		},
		MockGetProgramFacilitiesFn: func(ctx context.Context, programID string) ([]*domain.Facility, error) {
			return []*domain.Facility{&facility}, nil
		},
		MockListProgramsFn: func(ctx context.Context, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error) {
			return &programPage, nil
		},
		MockSetStaffProgramFn: func(ctx context.Context, programID string) (*domain.StaffResponse, error) {
			return &staffResponse, nil
		},
		MockGetProgramByIDFn: func(ctx context.Context, programID string) (*domain.Program, error) {
			return &program, nil
		},
		MockSetClientProgramFn: func(ctx context.Context, programID string) (*domain.ClientResponse, error) {
			return clientResponse, nil
		},
		MockSearchProgramsFn: func(ctx context.Context, searchParameter string) ([]*domain.Program, error) {
			return []*domain.Program{&program}, nil
		},
		MockListOrganisationProgramsFn: func(ctx context.Context, organisationID string, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error) {
			return &programPage, nil
		},
	}
}

// ListUserPrograms mock the implementation of the ListUserPrograms method
func (gm *ProgramsUseCaseMock) ListUserPrograms(ctx context.Context, userID string, flavour feedlib.Flavour) (*dto.ProgramOutput, error) {
	return gm.MockListUserProgramsFn(ctx, userID, flavour)
}

// SetCurrentProgram mock the implementation of the SetCurrentProgram method
func (gm *ProgramsUseCaseMock) SetCurrentProgram(ctx context.Context, programID string) (bool, error) {
	return gm.MockSetCurrentProgramFn(ctx, programID)
}

// GetProgramFacilities mock the implementation of the GetProgramFacilities method
func (gm *ProgramsUseCaseMock) GetProgramFacilities(ctx context.Context, programID string) ([]*domain.Facility, error) {
	return gm.MockGetProgramFacilitiesFn(ctx, programID)
}

// ListPrograms mock the implementation of the ListPrograms method
func (gm *ProgramsUseCaseMock) ListPrograms(ctx context.Context, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error) {
	return gm.MockListProgramsFn(ctx, paginationsInput)
}

// SetStaffProgram mock the implementation of the SetStaffProgram method
func (gm *ProgramsUseCaseMock) SetStaffProgram(ctx context.Context, programID string) (*domain.StaffResponse, error) {
	return gm.MockSetStaffProgramFn(ctx, programID)
}

// SetClientProgram mock the implementation of the SetClientProgram method
func (gm *ProgramsUseCaseMock) SetClientProgram(ctx context.Context, programID string) (*domain.ClientResponse, error) {
	return gm.MockSetClientProgramFn(ctx, programID)
}

// GetProgramByID mock the implementation of the GetProgramByID method
func (gm *ProgramsUseCaseMock) GetProgramByID(ctx context.Context, programID string) (*domain.Program, error) {
	return gm.MockGetProgramByIDFn(ctx, programID)
}

// CreateProgram mock the implementation of the CreateProgram method
func (gm *ProgramsUseCaseMock) CreateProgram(ctx context.Context, input *dto.ProgramInput) (bool, error) {
	return gm.MockCreateProgramFn(ctx, input)
}

// SearchPrograms mock the implementation of the SearchPrograms method
func (gm *ProgramsUseCaseMock) SearchPrograms(ctx context.Context, searchParameter string) ([]*domain.Program, error) {
	return gm.MockSearchProgramsFn(ctx, searchParameter)
}

// ListOrganisationPrograms mock the implementation of the ListOrganisationPrograms method
func (gm *ProgramsUseCaseMock) ListOrganisationPrograms(ctx context.Context, organisationID string, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error) {
	return gm.MockListOrganisationProgramsFn(ctx, organisationID, paginationsInput)
}
