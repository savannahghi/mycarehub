package profile

import (
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
