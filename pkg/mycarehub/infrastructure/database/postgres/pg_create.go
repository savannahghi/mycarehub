package postgres

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// GetOrCreateFacility is responsible from creating a representation of a facility
// A facility here is the healthcare facility that are on the platform.
// A facility MFL CODE must be unique across the platform. I forms part of the unique identifiers
//
// TODO: Create a helper the checks for all required fields
// TODO: Make the create method idempotent
func (d *MyCareHubDb) GetOrCreateFacility(ctx context.Context, facility *dto.FacilityInput) (*domain.Facility, error) {
	if err := facility.Validate(); err != nil {
		return nil, fmt.Errorf("facility input validation failed: %s", err)
	}

	facilityObj := &gorm.Facility{
		Name:        facility.Name,
		Code:        facility.Code,
		Active:      facility.Active,
		County:      facility.County,
		Phone:       facility.Phone,
		Description: facility.Description,
	}

	facilitySession, err := d.create.GetOrCreateFacility(ctx, facilityObj)
	if err != nil {
		return nil, fmt.Errorf("failed to create facility: %v", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// SaveTemporaryUserPin does the actual saving of the users PIN in the database
func (d *MyCareHubDb) SaveTemporaryUserPin(ctx context.Context, pinData *domain.UserPIN) (bool, error) {
	pinObj := &gorm.PINData{
		UserID:    pinData.UserID,
		HashedPIN: pinData.HashedPIN,
		ValidFrom: pinData.ValidFrom,
		ValidTo:   pinData.ValidTo,
		IsValid:   pinData.IsValid,
		Flavour:   pinData.Flavour,
		Salt:      pinData.Salt,
	}

	_, err := d.create.SaveTemporaryUserPin(ctx, pinObj)
	if err != nil {
		return false, fmt.Errorf("failed to save user pin: %v", err)
	}

	return true, nil
}

// SavePin gets the pin details from the user and saves it in the database
func (d *MyCareHubDb) SavePin(ctx context.Context, pinInput *domain.UserPIN) (bool, error) {

	pinObj := &gorm.PINData{
		UserID:    pinInput.UserID,
		HashedPIN: pinInput.HashedPIN,
		ValidFrom: pinInput.ValidFrom,
		ValidTo:   pinInput.ValidTo,
		IsValid:   pinInput.IsValid,
		Flavour:   pinInput.Flavour,
		Salt:      pinInput.Salt,
	}

	_, err := d.create.SavePin(ctx, pinObj)
	if err != nil {
		return false, fmt.Errorf("failed to save user pin: %v", err)
	}

	return true, nil
}

// SaveOTP saves the otp to the database
func (d *MyCareHubDb) SaveOTP(ctx context.Context, otpInput *domain.OTP) error {
	otpObject := &gorm.UserOTP{
		UserID:      otpInput.UserID,
		Valid:       otpInput.Valid,
		GeneratedAt: otpInput.GeneratedAt,
		ValidUntil:  otpInput.ValidUntil,
		Channel:     otpInput.Channel,
		PhoneNumber: otpInput.PhoneNumber,
		Flavour:     otpInput.Flavour,
		OTP:         otpInput.OTP,
	}

	err := d.create.SaveOTP(ctx, otpObject)
	if err != nil {
		return fmt.Errorf("failed to save OTP")
	}

	return nil
}

// SaveSecurityQuestionResponse saves the security question response to the database
func (d *MyCareHubDb) SaveSecurityQuestionResponse(ctx context.Context, securityQuestionResponse []*dto.SecurityQuestionResponseInput) error {
	var securityQuestionResponseObj []*gorm.SecurityQuestionResponse
	for _, sqr := range securityQuestionResponse {
		response := &gorm.SecurityQuestionResponse{
			UserID:     sqr.UserID,
			QuestionID: sqr.SecurityQuestionID,
			Active:     true,
			Response:   sqr.Response,
		}
		securityQuestionResponseObj = append(securityQuestionResponseObj, response)
	}

	err := d.create.SaveSecurityQuestionResponse(ctx, securityQuestionResponseObj)
	if err != nil {
		return fmt.Errorf("failed to save security question response data")
	}

	return nil
}

// CreateHealthDiaryEntry is used to add a health diary record to the database.
func (d *MyCareHubDb) CreateHealthDiaryEntry(ctx context.Context, healthDiaryInput *domain.ClientHealthDiaryEntry) error {
	healthDiaryResponse := &gorm.ClientHealthDiaryEntry{
		Active:                healthDiaryInput.Active,
		Mood:                  healthDiaryInput.Mood,
		Note:                  healthDiaryInput.Note,
		EntryType:             healthDiaryInput.EntryType,
		ShareWithHealthWorker: healthDiaryInput.ShareWithHealthWorker,
		SharedAt:              healthDiaryInput.SharedAt,
		ClientID:              healthDiaryInput.ClientID,
	}

	err := d.create.CreateHealthDiaryEntry(ctx, healthDiaryResponse)
	if err != nil {
		return err
	}

	return nil
}

// CreateServiceRequest creates  a service request which will be handled by a staff user.
// This happens in a transaction because we do not want to
// create a health diary entry without a subsequent service request when the client's mood is "VERY_BAD"
func (d *MyCareHubDb) CreateServiceRequest(
	ctx context.Context,
	serviceRequestInput *domain.ClientServiceRequest,
) error {
	serviceRequest := &gorm.ClientServiceRequest{
		Active:       serviceRequestInput.Active,
		RequestType:  serviceRequestInput.RequestType,
		Request:      serviceRequestInput.Request,
		Status:       serviceRequestInput.Status,
		InProgressAt: serviceRequestInput.InProgressAt,
		ResolvedAt:   serviceRequestInput.ResolvedAt,
		ClientID:     serviceRequestInput.ClientID,
	}

	err := d.create.CreateServiceRequest(ctx, serviceRequest)
	if err != nil {
		return err
	}

	return nil
}
