package utils

import (
	"gitlab.slade360emr.com/go/base"
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

// AddHashToCovers add a hash identifier to provided unhashed covers
func AddHashToCovers(unHashedCovers []base.Cover) []base.Cover {
	hashed := []base.Cover{}
	for _, cvr := range unHashedCovers {
		if cvr.IdentifierHash == nil {
			hashed = append(hashed, base.Cover{
				IdentifierHash: base.CreateCoverHash(base.Cover{
					PayerName:      cvr.PayerName,
					PayerSladeCode: cvr.PayerSladeCode,
					MemberNumber:   cvr.MemberNumber,
					MemberName:     cvr.MemberName,
				}),
				PayerName:      cvr.PayerName,
				PayerSladeCode: cvr.PayerSladeCode,
				MemberNumber:   cvr.MemberNumber,
				MemberName:     cvr.MemberName,
			})
		}
	}
	return hashed
}
