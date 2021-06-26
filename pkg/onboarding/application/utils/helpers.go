package utils

import (
	"fmt"
	"strconv"

	"gitlab.slade360emr.com/go/base"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// IfCoverExistsInSlice checks is a cover is in a slice.
// In a slice of covers of a given length (srcCovers), the cover is question will only
// match one if it does exist in the srcCovers. Hence;
// - if the difference between `coversThatExists` and `coversThatDontExist` is less than 0, it means the cover is question  does not exist in srcCovers .
// - if the difference between `coversThatExists` and `coversThatDontExist` is equal than 0, it means the cover is question does  exist in srcCovers.
func IfCoverExistsInSlice(srcCovers []base.Cover, cover base.Cover) bool {
	existCount := 0
	doesNotExistCount := 0

	// return early. There is nothing to check here
	if len(srcCovers) == 0 {
		return false
	}

	// check if the identifier in the current cover matches that to cover we are try to check for existence
	// if it does not match, increase doesNotExistCount.
	// if it does match, increase existCount.
	for _, current := range srcCovers {
		if *current.IdentifierHash != *cover.IdentifierHash {
			doesNotExistCount++
		} else {
			existCount++
		}
	}

	// if the difference between `existCount` and `doesNotExistCount` is less than 0, it means the cover is question  does not exist in srcCovers
	if (existCount - doesNotExistCount) < 0 {
		return false
	}

	// if the difference between `existCount` and `doesNotExistCount` is equal than 0, it means the cover is question does  exist in srcCovers
	// we do this explicitly instead of relaying on the fallback
	if (existCount - doesNotExistCount) == 0 {
		return true
	}

	// fallback.
	return true
}

// CheckIdentifierExists check if an identifier exists
func CheckIdentifierExists(profile *base.UserProfile, UID string) bool {
	foundVerifiedUIDs := []string{}
	verifiedIDs := profile.VerifiedIdentifiers
	for _, verifiedID := range verifiedIDs {
		foundVerifiedUIDs = append(foundVerifiedUIDs, verifiedID.UID)
	}
	return base.StringSliceContains(foundVerifiedUIDs, UID)
}

// CheckUserHasFavNavAction checks if user has book marked the provided navaction
func CheckUserHasFavNavAction(u *base.UserProfile, title string) bool {
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

// AddHashToCovers add a hash identifier to provided unhashed covers
func AddHashToCovers(unHashedCovers []base.Cover) []base.Cover {
	hashed := []base.Cover{}
	for _, cvr := range unHashedCovers {
		if cvr.IdentifierHash == nil {
			hashed = append(hashed, base.Cover{
				IdentifierHash: base.CreateCoverHash(base.Cover{
					PayerName:             cvr.PayerName,
					PayerSladeCode:        cvr.PayerSladeCode,
					MemberNumber:          cvr.MemberNumber,
					MemberName:            cvr.MemberName,
					BeneficiaryID:         cvr.BeneficiaryID,
					EffectivePolicyNumber: cvr.EffectivePolicyNumber,
					ValidFrom:             cvr.ValidFrom,
					ValidTo:               cvr.ValidTo,
				}),
				PayerName:             cvr.PayerName,
				PayerSladeCode:        cvr.PayerSladeCode,
				MemberNumber:          cvr.MemberNumber,
				MemberName:            cvr.MemberName,
				BeneficiaryID:         cvr.BeneficiaryID,
				EffectivePolicyNumber: cvr.EffectivePolicyNumber,
				ValidFrom:             cvr.ValidFrom,
				ValidTo:               cvr.ValidTo,
			})
		}
	}
	return hashed
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

// ParseUSSDDateInput ...
func ParseUSSDDateInput(date string) string {
	dayEntered, _ := strconv.Atoi(date[0:2])
	monthEntered, _ := strconv.Atoi(date[2:4])
	yearEntered, _ := strconv.Atoi(date[4:8])

	return fmt.Sprintf("%v-%v-%v", dayEntered, monthEntered, yearEntered)
}

// UniquePermissionsArray removes duplicate permissions in a slice of permissions
func UniquePermissionsArray(arr []base.PermissionType) []base.PermissionType {
	occured := map[base.PermissionType]bool{}
	result := []base.PermissionType{}

	for e := range arr {
		// check if already the mapped
		// variable is set to true or not
		if !occured[arr[e]] {
			occured[arr[e]] = true

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
