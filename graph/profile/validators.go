package profile

import (
	"context"
	"fmt"

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
