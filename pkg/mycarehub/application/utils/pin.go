package utils

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"math/big"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/xdg-go/pbkdf2"
)

const (
	// DefaultSaltLen is the length of generated salt for the user is 256
	DefaultSaltLen = 256
	// defaultIterations is the iteration count in PBKDF2 function is 10000
	defaultIterations = 10000
	// DefaultKeyLen is the length of encoded key in PBKDF2 function is 512
	DefaultKeyLen = 512
	// alphanumeric character used for generation of a `salt`
	alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	// Default length of a generated pin
	generatedPinLength = 4
	// Default min length of the date
	minPINLength = 4
	// Default max length of the date
	maxPINLength = 4
)

// DefaultHashFunction ...
var DefaultHashFunction = sha512.New

// Options is a struct for custom values of salt length, number of iterations, the encoded key's length,
// and the hash function being used. If set to `nil`, default options are used:
// &Options{ 256, 10000, 512, "sha512" }
type Options struct {
	SaltLen      int
	Iterations   int
	KeyLen       int
	HashFunction func() hash.Hash
}

func generateSalt(length int) []byte {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil
	}
	for key, val := range salt {
		salt[key] = alphanum[val%byte(len(alphanum))]
	}
	return salt
}

// EncryptPIN takes two arguments, a raw pin, and a pointer to an Options struct.
// In order to use default options, pass `nil` as the second argument.
// It returns the generated salt and encoded key for the user.
func EncryptPIN(rawPwd string, options *Options) (string, string) {
	if options == nil {
		salt := generateSalt(DefaultSaltLen)
		encodedPwd := pbkdf2.Key([]byte(rawPwd), salt, defaultIterations, DefaultKeyLen, DefaultHashFunction)
		return string(salt), hex.EncodeToString(encodedPwd)
	}
	salt := generateSalt(options.SaltLen)
	encodedPwd := pbkdf2.Key([]byte(rawPwd), salt, options.Iterations, options.KeyLen, options.HashFunction)
	return string(salt), hex.EncodeToString(encodedPwd)
}

// ComparePIN takes four arguments, the raw password, its generated salt, the encoded password,
// and a pointer to the Options struct, and returns a boolean value determining whether the password is the correct one or not.
// Passing `nil` as the last argument resorts to default options.
func ComparePIN(rawPwd string, salt string, encodedPwd string, options *Options) bool {
	if options == nil {
		return encodedPwd == hex.EncodeToString(pbkdf2.Key([]byte(rawPwd), []byte(salt), defaultIterations, DefaultKeyLen, DefaultHashFunction))
	}
	return encodedPwd == hex.EncodeToString(pbkdf2.Key([]byte(rawPwd), []byte(salt), options.Iterations, options.KeyLen, options.HashFunction))
}

// GenerateTempPIN generates a temporary One Time PIN for a user
// The PIN will have 4 digits formatted as a string
func GenerateTempPIN(ctx context.Context) (string, error) {
	var pin string

	length := 0
	for length < generatedPinLength {
		number, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}

		pin += number.String()

		length++
	}

	return pin, nil
}

// ValidatePIN is used to check for the validity of the PIN provided.
func ValidatePIN(pin string) error {
	validatePINErr := ValidatePINLength(pin)
	if validatePINErr != nil {
		return validatePINErr
	}

	pinDigitsErr := ValidatePINDigits(pin)
	if pinDigitsErr != nil {
		return pinDigitsErr
	}
	return nil
}

// ValidatePINLength ...
func ValidatePINLength(pin string) error {
	// make sure pin length is [4]
	if len(pin) < minPINLength || len(pin) > maxPINLength {
		return fmt.Errorf("PIN should be of 4 digits")
	}
	return nil
}

// ValidatePINDigits validates user pin to ensure a PIN only contains digits
func ValidatePINDigits(pin string) error {
	// ensure pin is only digits
	_, err := strconv.ParseUint(pin, 10, 64)
	if err != nil {
		return exceptions.ValidatePINDigitsErr(err)
	}
	return nil
}
