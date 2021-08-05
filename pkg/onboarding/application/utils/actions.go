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

	for _, action := range domain.AllNavigationActions {

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
) map[domain.NavigationGroup]domain.NavigationAction {
	grouped := make(map[domain.NavigationGroup]domain.NavigationAction)

	// map all parents to grouped
	for _, action := range actions {
		if !action.HasParent {
			grouped[action.Group] = action
		}
	}

	//map all children to respective parents
	for _, action := range actions {
		if action.HasParent {
			parent, isAvailable := grouped[action.Group]
			if isAvailable {
				parent.Nested = append(parent.Nested, action)
				grouped[action.Group] = parent
			}
		}
	}

	return grouped
}

// GroupPriority groups navigation actions into primary and secondary actions
func GroupPriority(
	actions map[domain.NavigationGroup]domain.NavigationAction,
) (primary, secondary []domain.NavigationAction) {
	primary = []domain.NavigationAction{}
	secondary = []domain.NavigationAction{}

	added := make(map[domain.NavigationGroup]domain.NavigationAction)

	//pb is number of navactions that can possibly be on bottom navigation ie primary
	pb := 0
	for _, v := range actions {
		if len(v.Nested) == 0 {
			pb++
		}
	}

	// add all the possible bottom action to primary if they are less or equal to 4
	if pb <= 4 {
		for _, action := range actions {
			if len(action.Nested) == 0 {
				primary = append(primary, action)
				added[action.Group] = action
			}
		}
	} else {
		for {
			if len(primary) == 4 {
				break
			}
			// add all the high priority first
			for _, action := range actions {

				if action.IsHighPriority {

					_, exist := added[action.Group]
					if !exist && len(action.Nested) == 0 {

						primary = append(primary, action)
						added[action.Group] = action

						if len(primary) == 4 {
							break
						}
					}
				}
			}
			// add other actions is high priority actions are less than four
			for _, action := range actions {

				_, exist := added[action.Group]
				if !exist && len(action.Nested) == 0 {

					primary = append(primary, action)
					added[action.Group] = action

					if len(primary) == 4 {
						break
					}
				}
			}
		}
	}

	// add all remaining items to secondary
	for _, action := range actions {

		_, exists := added[action.Group]
		if !exists {

			secondary = append(secondary, action)
			added[action.Group] = action
		}
	}

	return primary, secondary
}
