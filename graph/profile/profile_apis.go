package profile

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// GetProfileAttributes given a slice of user uids returns specified user profile attributes
func (s Service) GetProfileAttributes(ctx context.Context, uids []string, attribute string) (map[string][]string, error) {
	output := make(map[string][]string)
	values := []string{}

	for _, uid := range uids {
		profile, err := s.GetProfile(ctx, uid)
		if err != nil {
			return output, fmt.Errorf("unable to get user profile: %v", err)
		}
		switch attribute {
		case "emails":
			values = append(values, profile.Emails...)
		case "phoneNumbers":
			values = append(values, profile.Msisdns...)
		case "fcmTokens":
			values = append(values, profile.PushTokens...)
		default:
			log.Panicf("can't get the user profile attribute provided")
		}
		output[uid] = values
	}

	return output, nil
}

// GetConfirmedEmailAddresses returns verified email addresses for the uids provided
func (s Service) GetConfirmedEmailAddresses(ctx context.Context, uids []string) (map[string][]string, error) {
	s.checkPreconditions()

	attribute := "emails"
	return s.GetProfileAttributes(ctx, uids, attribute)
}

// GetConfirmedPhoneNumbers returns verified phone numbers for the uids provided
func (s Service) GetConfirmedPhoneNumbers(ctx context.Context, uids []string) (map[string][]string, error) {
	s.checkPreconditions()

	attribute := "phoneNumbers"
	return s.GetProfileAttributes(ctx, uids, attribute)
}

// GetValidFCMTokens returns valid FCM push tokens for the uids provided
func (s Service) GetValidFCMTokens(ctx context.Context, uids []string) (map[string][]string, error) {
	s.checkPreconditions()

	attribute := "fcmTokens"
	return s.GetProfileAttributes(ctx, uids, attribute)
}

// GetAttribute return users' email, phonenumber or fcm token attributes
func GetAttribute(ctx context.Context, r *http.Request, uids []string) (map[string][]string, error) {
	s := NewService()
	if r == nil {
		return nil, fmt.Errorf("nil request")
	}

	pathVars := mux.Vars(r)
	attribute, found := pathVars["attribute"]
	if !found {
		return nil, fmt.Errorf("the request does not have a path var named `%s`", attribute)
	}
	output := make(map[string][]string)

	switch attribute {
	case "emails":
		confirmedEmails, err := s.GetConfirmedEmailAddresses(ctx, uids)
		if err != nil {
			return output, fmt.Errorf("unable to get confirmed email addresses: %v", err)
		}
		output = confirmedEmails

	case "phonenumbers":
		confirmedPhoneNumbers, err := s.GetConfirmedPhoneNumbers(ctx, uids)
		if err != nil {
			return output, fmt.Errorf("unable to get confirmed phone numbers: %v", err)
		}
		output = confirmedPhoneNumbers

	case "tokens":
		validFCMTokens, err := s.GetValidFCMTokens(ctx, uids)
		if err != nil {
			return output, fmt.Errorf("unable to get valid fcm tokens: %v", err)
		}
		output = validFCMTokens

	default:
		log.Panicf("can't get the user profile attribute provided")
	}

	return output, nil
}
