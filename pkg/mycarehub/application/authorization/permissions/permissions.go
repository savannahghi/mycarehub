package permissions

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"

// InviteUserToCommunity is a permission that allows a user to invite another user to a community
var InviteUserToCommunity = dto.PermissionInput{
	Object: "community.users",
	Action: "invite",
}

// ListCommunityMembers is a permission that allows a user to invite another user to a community
var ListCommunityMembers = dto.PermissionInput{
	Object: "community.users",
	Action: "list",
}

// CreateFacility is a permission that allows a user to invite another user to a community
var CreateFacility = dto.PermissionInput{
	Object: "facility",
	Action: "create",
}
