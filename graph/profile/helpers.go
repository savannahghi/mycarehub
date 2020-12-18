package profile

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"
	"gitlab.slade360emr.com/go/base"
	"golang.org/x/crypto/bcrypt"
)

// EncryptPIN encrypts a string
func EncryptPIN(pin string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("unable to hash PIN %w", err)
	}
	return string(bytes), nil
}

// ComparePIN compare two PINs to see if they match
func ComparePIN(hashedPin, plainPin string) (bool, error) {
	// convert hashed PIN to byte
	byteHash := []byte(hashedPin)
	plainPinHash := []byte(plainPin)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPinHash)
	if err != nil {
		return false, fmt.Errorf("PIN mismatch %w", err)
	}
	return true, nil
}

func (s Service) encryptExistingPin(
	p *PIN,
	dsnap *firestore.DocumentSnapshot,
) error {
	_, convertPinToIntErr := strconv.Atoi(p.PINNumber)
	if convertPinToIntErr == nil {
		// if the pin is converted successfully then it implies that it has numbers only
		// which means that it was probably not encrypted before
		newEncryptedPin, err := EncryptPIN(p.PINNumber)
		if err != nil {
			return fmt.Errorf("PhoneSignIn: unable to encrypt PIN: %w", err)
		}

		p.PINNumber = newEncryptedPin
		err = base.UpdateRecordOnFirestore(
			s.firestoreClient, s.GetPINCollectionName(), dsnap.Ref.ID, p,
		)
		if err != nil {
			return fmt.Errorf("PhoneSignIn: unable to update pins record: %v", err)
		}
	}
	return nil
}


func validatePIN(pin string) error {
	// make sure pin is of only digits
	_, err := strconv.ParseUint(pin, 10, 64)
	if err != nil {
		return fmt.Errorf("pin should be a valid number: %w", err)
	}

	// make sure pin length is [4-6]
	if len(pin) < 4 || len(pin) > 6 {
		return fmt.Errorf("pin should be of 4,5, or six digits")
	}
	return nil
}

// ResetPINHelper resets a user PIN to the new one supplied
func (s Service) ResetPINHelper(ctx context.Context, msisdn string, pin string) (bool, error){
	// validate the PIN
	err := validatePIN(pin)
	if err != nil {
		return false, fmt.Errorf("invalid pin: %w", err)
	}
	// ensure the phone number is valid
	phoneNumber, err := base.NormalizeMSISDN(msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to normalize the msisdn: %v", err)
	}
	// check if user has existing PIN
	exists, err := s.CheckHasPIN(ctx, msisdn)
	if err != nil {
		return false, fmt.Errorf("unable to check if the user has a PIN: %v", err)
	}
	// if the user already has an existing PIN update with the new one 
	if exists {
		return s.UpdateUserPIN(ctx, msisdn, pin)
	}
	// EncryptPIN the PIN
	encryptedPin, err := EncryptPIN(pin)
	if err != nil {
		return false, fmt.Errorf("unable to encrypt PIN: %w", err)
	}
	// prepare PIN paload
	PINPayload := PIN{
		MSISDN:    phoneNumber,
		PINNumber: encryptedPin,
		IsValid:   true,
	}
	// store the PIN
	err = s.SavePINToFirestore(PINPayload)
	if err != nil {
		return false, fmt.Errorf("unable to save PIN: %v", err)
	}
	return true, nil
}