package profile

import (
	"fmt"

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
