package utils

import (
	"context"
	"sort"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/profileutils"
)

//CheckUserHasPermission takes in the user roles and a permission and verifies that the user
//has required permissions
func CheckUserHasPermission(roles []profileutils.Role, permission profileutils.Permission) bool {
	scopes := []string{}
	for _, role := range roles {
		// only add scopes of active roles
		if role.Active {
			scopes = append(scopes, role.Scopes...)
		}
	}

	for _, scope := range scopes {
		if permission.Scope == scope {
			return true
		}
	}
	return false
}

// GetUserNavigationActions returns a sorted primary and secondary user navigation actions
func GetUserNavigationActions(
	ctx context.Context,
	user profileutils.UserProfile,
	roles []profileutils.Role,
) (*dto.GroupedNavigationActions, error) {
	//all user actions
	userNavigationActions := []domain.NavigationAction{}

	allActions := domain.AllNavigationActions
	for i := 0; i < len(allActions); i++ {
		action := allActions[i]
		if action.RequiredPermission == nil || CheckUserHasPermission(roles, *action.RequiredPermission) {
			//  check for favorite navigation actions
			if IsFavNavAction(&user, action.Title) {
				action.Favorite = true
			}

			userNavigationActions = append(userNavigationActions, action)
		}
	}

	groupNested := GroupNested(userNavigationActions)
	primary, secondary := GroupPriority(groupNested)

	navigationActions := &dto.GroupedNavigationActions{
		Primary:   primary,
		Secondary: secondary,
	}
	return navigationActions, nil
}

// GroupNested groups navigation actions into parents and children
func GroupNested(
	actions []domain.NavigationAction,
) []domain.NavigationAction {

	// Array of all parent actions i.e can have nested actions
	parents := []domain.NavigationAction{}
	for _, action := range actions {
		if !action.HasParent {
			parents = append(parents, action)
		}
	}

	// An array of properly grouped actions
	// The parent action has the nested actions
	grouped := []domain.NavigationAction{}
	for _, parent := range parents {
		for _, action := range actions {
			if action.HasParent && action.Group == parent.Group {
				parent.Nested = append(parent.Nested, action)
			}
		}
		grouped = append(grouped, parent)
	}

	return grouped
}

// GroupPriority groups navigation actions into primary and secondary actions
func GroupPriority(
	actions []domain.NavigationAction,
) (primary, secondary []domain.NavigationAction) {

	// sort actions based on priority using the sequence number
	// uses the inbuilt go sorting functionality
	// https://cs.opensource.google/go/go/+/go1.16.7:src/sort/slice.go;l=16
	sort.SliceStable(actions, func(i, j int) bool {
		return actions[i].SequenceNumber < actions[j].SequenceNumber
	})

	primary = []domain.NavigationAction{}
	secondary = []domain.NavigationAction{}

	// this helps keep track of grouped actions
	tracker := make(map[domain.NavigationGroup]bool)

	// pb is number of actions without nested actions
	pb := 0
	for _, a := range actions {
		if len(a.Nested) == 0 {
			pb++
		}
	}

	// add all the possible bottom action to primary if they are less or equal to 4
	if pb <= 4 {

		for _, action := range actions {
			if len(action.Nested) == 0 {
				primary = append(primary, action)
				tracker[action.Group] = true
			}
		}

	} else {

		// add all the high priority first
		for _, action := range actions {

			_, added := tracker[action.Group]

			// Add the primary action it was not added
			// And the action lacks nested actions
			if !added && len(action.Nested) == 0 {

				primary = append(primary, action)
				tracker[action.Group] = true

				if len(primary) == 4 {
					break
				}
			}

		}

	}

	// add all remaining items to secondary
	for _, action := range actions {

		_, added := tracker[action.Group]

		if !added {
			secondary = append(secondary, action)
			tracker[action.Group] = true
		}
	}

	// sort the primary and secondary actions based on priority again
	// this is a precautionary step since all actions were sorted before
	sort.SliceStable(primary, func(i, j int) bool {
		return primary[i].SequenceNumber < primary[j].SequenceNumber
	})

	sort.SliceStable(secondary, func(i, j int) bool {
		return secondary[i].SequenceNumber < secondary[j].SequenceNumber
	})

	return primary, secondary
}
