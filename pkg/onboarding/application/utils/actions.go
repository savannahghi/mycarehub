package utils

import (
	"context"

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
	parents := []domain.NavigationAction{}
	grouped := []domain.NavigationAction{}
	// map all parents to grouped
	for i := 0; i < len(actions); i++ {
		action := actions[i]

		if !action.HasParent {
			parents = append(parents, action)
		}
	}

	//map all children to respective parents
	for _, parent := range parents {
		for i := 0; i < len(actions); i++ {
			action := actions[i]
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
	primary = []domain.NavigationAction{}
	secondary = []domain.NavigationAction{}

	mapped := make(map[domain.NavigationGroup]bool)

	//pb is number of navactions that can possibly be on bottom navigation ie primary
	pb := 0
	for _, v := range actions {
		if len(v.Nested) == 0 {
			pb++
		}
	}

	// add all the possible bottom action to primary if they are less or equal to 4
	if pb <= 4 {
		for i := 0; i < len(actions); i++ {
			action := actions[i]
			if len(action.Nested) == 0 {
				primary = append(primary, action)
				mapped[action.Group] = true
			}
		}
	} else {
		for {
			if len(primary) == 4 {
				break
			}
			// add all the high priority first
			for i := 0; i < len(actions); i++ {
				action := actions[i]
				if action.IsHighPriority {

					_, wasMapped := mapped[action.Group]
					if !wasMapped && len(action.Nested) == 0 {

						primary = append(primary, action)
						mapped[action.Group] = true

						if len(primary) == 4 {
							break
						}
					}
				}
			}
			// add other actions is high priority actions are less than four
			for i := 0; i < len(actions); i++ {
				action := actions[i]
				_, wasMapped := mapped[action.Group]
				if !wasMapped && len(action.Nested) == 0 {

					primary = append(primary, action)
					mapped[action.Group] = true

					if len(primary) == 4 {
						break
					}
				}
			}
		}
	}

	// add all remaining items to secondary
	for i := 0; i < len(actions); i++ {
		action := actions[i]
		_, wasMapped := mapped[action.Group]
		if !wasMapped {

			secondary = append(secondary, action)
			mapped[action.Group] = true
		}
	}

	return primary, secondary
}
