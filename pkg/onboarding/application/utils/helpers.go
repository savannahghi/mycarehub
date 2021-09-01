package utils

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/savannahghi/converterandformatter"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/profileutils"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// CheckIdentifierExists check if an identifier exists
func CheckIdentifierExists(profile *profileutils.UserProfile, UID string) bool {
	foundVerifiedUIDs := []string{}
	verifiedIDs := profile.VerifiedIdentifiers
	for _, verifiedID := range verifiedIDs {
		foundVerifiedUIDs = append(foundVerifiedUIDs, verifiedID.UID)
	}
	return converterandformatter.StringSliceContains(foundVerifiedUIDs, UID)
}

// IsFavNavAction checks if user has book marked the provided navaction
func IsFavNavAction(u *profileutils.UserProfile, title string) bool {
	if len(u.FavNavActions) == 0 {
		return false
	}
	for _, t := range u.FavNavActions {
		if t == title {
			return true
		}
	}
	return false
}

// MatchAndReturn checks is the new is not equal to old and returns it. Otherwise the old
// used when setting communication settings
func MatchAndReturn(old bool, new bool) bool {
	if (new && !old) || (!new && old) {
		return new
	}
	return old
}

// FindItem is a linear search implementation that checks if a string
// is contained in a slice and returns its index and a boolean
// The index can be used for slicing out an item from a slice.
func FindItem(slice []string, value string) (int, bool) {
	for index, val := range slice {
		if val == value {
			return index, true
		}
	}
	return -1, false
}

// UniquePermissionsArray removes duplicate permissions in an array of permissions
func UniquePermissionsArray(arr []profileutils.PermissionType) []profileutils.PermissionType {
	occurred := map[profileutils.PermissionType]bool{}
	result := []profileutils.PermissionType{}

	for e := range arr {
		// check if already the mapped
		// variable is set to true or not
		if !occurred[arr[e]] {
			occurred[arr[e]] = true

			// Append to result slice.
			result = append(result, arr[e])
		}
	}

	return result
}

// RecordSpanError is a helper function to capture errors in a span
func RecordSpanError(span trace.Span, err error) {
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}

//ServiceHealthEndPoint creates a url to the service health endpoint
func ServiceHealthEndPoint(input string) (string, error) {
	reg, err := regexp.CompilePOSIX("/graphql")
	if err != nil {
		log.Print(err)
		return "", err
	}
	str := reg.Split(input, 2)
	endpointURL := str[0] + "/" + "health"

	return endpointURL, nil
}

// CheckEmptyString check that the string isn't empty
func CheckEmptyString(text string) (*string, error) {
	inputText := text
	if inputText == "" {
		return nil, fmt.Errorf("sessionID cannot be empty")
	}
	return &inputText, nil
}

// NewActionsMapper maps the new navigation actions mapper to the old navigation actions implementation
// It is meant to  be used as a transition and ease changes made without affecting the frontend
func NewActionsMapper(ctx context.Context, grouped *dto.GroupedNavigationActions) *profileutils.NavigationActions {
	mapped := &profileutils.NavigationActions{}

	for i := 0; i < len(grouped.Primary); i++ {
		action := grouped.Primary[i]
		c := profileutils.NavAction{
			Title:      action.Title,
			OnTapRoute: action.OnTapRoute,
			Favourite:  action.Favorite,
			Icon:       feedlib.GetSVGImageLink(action.Icon, action.Title, action.Title, action.Title),
		}
		mapped.Primary = append(mapped.Primary, c)
	}

	for i := 0; i < len(grouped.Secondary); i++ {
		action := grouped.Secondary[i]
		c := profileutils.NavAction{
			Title:      action.Title,
			OnTapRoute: action.OnTapRoute,
			Favourite:  action.Favorite,
			Icon:       feedlib.GetSVGImageLink(action.Icon, action.Title, action.Title, action.Title),
		}

		if len(action.Nested) > 0 {

			for i := 0; i < len(action.Nested); i++ {
				nestedAction := (action.Nested[i]).(domain.NavigationAction)
				m := profileutils.NestedNavAction{
					Title:      nestedAction.Title,
					OnTapRoute: nestedAction.OnTapRoute,
				}

				c.Nested = append(c.Nested, m)
			}
		}

		mapped.Secondary = append(mapped.Secondary, c)
	}

	return mapped
}
