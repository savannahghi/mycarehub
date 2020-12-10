package profile

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	emailsAttribute       = "emails"
	phoneNumbersAttribute = "phonenumbers"
	fcmTokensAttribute    = "tokens"
)

// GetProfileAttributes given a slice of user uids returns the specified
// user profile attributes
func (s Service) GetProfileAttributes(
	ctx context.Context,
	uids []string,
	attribute string,
) (map[string][]string, error) {
	output := make(map[string][]string)
	values := []string{}

	for _, uid := range uids {
		profile, err := s.GetProfile(ctx, uid)
		if err != nil {
			return output, fmt.Errorf("unable to get user profile: %v", err)
		}
		if profile == nil {
			return nil, fmt.Errorf("no profile with UID %s", uid)
		}
		switch attribute {
		case emailsAttribute:
			values = append(values, profile.Emails...)
		case phoneNumbersAttribute:
			values = append(values, profile.Msisdns...)
		case fcmTokensAttribute:
			values = append(values, profile.PushTokens...)
		default:
			return nil, fmt.Errorf(
				"can't get the user profile attribute '%s'", attribute)
		}
		output[uid] = values
	}

	return output, nil
}

// GetConfirmedEmailAddresses returns verified email addresses for the uids
func (s Service) GetConfirmedEmailAddresses(
	ctx context.Context,
	uids []string,
) (map[string][]string, error) {
	s.checkPreconditions()
	return s.GetProfileAttributes(ctx, uids, emailsAttribute)
}

// GetConfirmedPhoneNumbers returns verified phone numbers for the uids
func (s Service) GetConfirmedPhoneNumbers(
	ctx context.Context,
	uids []string,
) (map[string][]string, error) {
	s.checkPreconditions()
	return s.GetProfileAttributes(ctx, uids, phoneNumbersAttribute)
}

// GetValidFCMTokens returns valid FCM push tokens for the uids provided
func (s Service) GetValidFCMTokens(
	ctx context.Context,
	uids []string,
) (map[string][]string, error) {
	s.checkPreconditions()
	return s.GetProfileAttributes(ctx, uids, fcmTokensAttribute)
}

// GetAttribute return users' email, phonenumber or fcm token attributes
func GetAttribute(
	ctx context.Context,
	r *http.Request,
	uids []string,
) (map[string][]string, error) {
	s := NewService()
	if r == nil {
		return nil, fmt.Errorf("nil request")
	}

	pathVars := mux.Vars(r)
	attribute, found := pathVars["attribute"]
	if !found {
		return nil, fmt.Errorf(
			"the request does not have a path var named `%s`", attribute)
	}
	output := make(map[string][]string)

	switch attribute {
	case emailsAttribute:
		confirmedEmails, err := s.GetConfirmedEmailAddresses(ctx, uids)
		if err != nil {
			return output, fmt.Errorf(
				"unable to get confirmed email addresses: %v", err)
		}
		output = confirmedEmails

	case phoneNumbersAttribute:
		confirmedPhoneNumbers, err := s.GetConfirmedPhoneNumbers(ctx, uids)
		if err != nil {
			return output, fmt.Errorf(
				"unable to get confirmed phone numbers: %v", err)
		}
		output = confirmedPhoneNumbers

	case fcmTokensAttribute:
		validFCMTokens, err := s.GetValidFCMTokens(ctx, uids)
		if err != nil {
			return output, fmt.Errorf(
				"unable to get valid fcm tokens: %v", err)
		}
		output = validFCMTokens

	default:
		return nil, fmt.Errorf(
			"can't get the user profile attribute %s", attribute)
	}

	return output, nil
}
