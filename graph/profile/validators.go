package profile

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/asaskevich/govalidator"
	"gitlab.slade360emr.com/go/base"
)

// ValidateEmail returns an email if the email and verification code are valid
func ValidateEmail(email, verificationCode string, firestoreClient *firestore.Client) (string, error) {
	if !govalidator.IsEmail(email) {
		return "", fmt.Errorf("invalid email format")
	}

	query := firestoreClient.Collection(base.SuffixCollection(base.OTPCollectionName)).Where(
		"isValid", "==", true,
	).Where(
		"authorizationCode", "==", verificationCode,
	).Where(
		"email", "==", email,
	)
	ctx := context.Background()
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return "", fmt.Errorf("unable to retrieve verification codes: %v", err)
	}
	if len(docs) == 0 {
		return "", fmt.Errorf("no matching verification codes found")
	}
	for _, doc := range docs {
		otpData := doc.Data()
		otpData["isValid"] = false
		err = base.UpdateRecordOnFirestore(
			firestoreClient, base.SuffixCollection(base.OTPCollectionName), doc.Ref.ID, otpData)
		if err != nil {
			return "", fmt.Errorf("unable to save updated OTP document: %v", err)
		}
	}
	return email, nil
}

// ValidateResetPinPayload checks that the request payload supplied in the indicated request are valid
func ValidateResetPinPayload(w http.ResponseWriter, r *http.Request) (*PinRecovery, error) {
	payload := &PinRecovery{}
	base.DecodeJSONToTargetStruct(w, r, payload)
	if payload.MSISDN == "" || payload.PINNumber == "" || payload.OTP == "" {
		err := fmt.Errorf("invalid pin update payload, expected a phone number, pin and an otp")
		base.ReportErr(w, err, http.StatusBadRequest)
		return nil, err
	}
	return payload, nil
}

// ValidateCreateUserByPhonePayload checks that the request payload supplied in the indicated request are valid
func ValidateCreateUserByPhonePayload(w http.ResponseWriter, r *http.Request) (*CreateUserViaPhoneInput, error) {
	payload := &CreateUserViaPhoneInput{}
	base.DecodeJSONToTargetStruct(w, r, payload)
	if payload.MSISDN == "" {
		err := fmt.Errorf("invalid create user payload, expected a phone number")
		base.ReportErr(w, err, http.StatusBadRequest)
		return nil, err
	}
	return payload, nil
}

// ValidateMsisdn checks that the msisdn supplied in the indicated request is valid
func ValidateMsisdn(w http.ResponseWriter, r *http.Request) (*PinRecovery, error) {
	data := &PinRecovery{}
	base.DecodeJSONToTargetStruct(w, r, data)
	if data.MSISDN == "" {
		err := fmt.Errorf("invalid credentials, expected a phone number")
		base.ReportErr(w, err, http.StatusBadRequest)
		return nil, err
	}
	return data, nil
}

// ValidateUID checks that the uid supplied in the indicated request is valid
func ValidateUID(w http.ResponseWriter, r *http.Request) (*BusinessPartnerUID, error) {
	p := &BusinessPartnerUID{}
	base.DecodeJSONToTargetStruct(w, r, p)
	if p.UID == "" {
		err := fmt.Errorf("invalid credentials, expected a uid")
		return nil, err
	}
	if p == nil {
		err := fmt.Errorf(
			"nil business partner UID struct after decoding input")
		return nil, err
	}

	return p, nil
}

// ValidateUserProfileUIDs checks that the uids supplied in the indicated request are valid
func ValidateUserProfileUIDs(w http.ResponseWriter, r *http.Request) (*UserUIDs, error) {
	uids := &UserUIDs{}
	base.DecodeJSONToTargetStruct(w, r, uids)
	if len(uids.UIDs) == 0 {
		err := fmt.Errorf("invalid credentials, expected a slice of uids")
		base.ReportErr(w, err, http.StatusBadRequest)
		return nil, err
	}
	return uids, nil
}

// ValidateSendRetryOTPPayload checks the validity of the request payload
func ValidateSendRetryOTPPayload(w http.ResponseWriter, r *http.Request) (*SendRetryOTP, error) {
	payload := &SendRetryOTP{}
	base.DecodeJSONToTargetStruct(w, r, payload)
	if payload.Msisdn == "" || payload.RetryStep == 0 {
		err := fmt.Errorf("invalid generate generates and fallback otp payload")
		base.ReportErr(w, err, http.StatusBadRequest)
		return nil, err
	}
	return payload, nil
}

// ValidatePhoneSignInInput checks that the credentials supplied in the indicated request are valid
func ValidatePhoneSignInInput(w http.ResponseWriter, r *http.Request) (*PhoneSignInInput, error) {
	payload := &PhoneSignInInput{}
	base.DecodeJSONToTargetStruct(w, r, payload)
	_, err := base.NormalizeMSISDN(payload.PhoneNumber)
	if err != nil || payload.PhoneNumber == "" || payload.Pin == "" {
		err := fmt.Errorf("expected a correct value")
		base.ReportErr(w, err, http.StatusBadRequest)
		return nil, err
	}
	return payload, nil
}
