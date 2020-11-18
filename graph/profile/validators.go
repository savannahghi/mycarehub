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

// ValidateUpdatePinPayload checks that the request payload supplied in the indicated request are valid
func ValidateUpdatePinPayload(w http.ResponseWriter, r *http.Request) (*PinRecovery, error) {
	payload := &PinRecovery{}
	base.DecodeJSONToTargetStruct(w, r, payload)
	if payload.MSISDN == "" || payload.PIN == "" || payload.OTP == "" {
		err := fmt.Errorf("invalid pin update payload, expected a phone number, pin and an otp")
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
	if p.UID == nil {
		err := fmt.Errorf("invalid credentials, expected a uid")
		base.ReportErr(w, err, http.StatusBadRequest)
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
