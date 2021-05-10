package permission

import (
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

// UserProfileView describes read permissions on viewing a user profile
var UserProfileView = resources.PermissionInput{
	Resource: "user_profile_view",
	Action:   "view",
}

// PrimaryPhoneUpdate describes update permissions on a user primary phonenumber
var PrimaryPhoneUpdate = resources.PermissionInput{
	Resource: "update_primary_phone",
	Action:   "edit",
}

// PrimaryEmailUpdate describes update permissions on a user's primary email address
var PrimaryEmailUpdate = resources.PermissionInput{
	Resource: "update_primary_email",
	Action:   "edit",
}

// SecondaryPhoneNumberUpdate describes update permissions on the user secondary phonenumber
var SecondaryPhoneNumberUpdate = resources.PermissionInput{
	Resource: "update_secondary_phone",
	Action:   "edit",
}

// SecondaryEmailAddressUpdate describes update permissions on the user secondary email address
var SecondaryEmailAddressUpdate = resources.PermissionInput{
	Resource: "update_secondary_email",
	Action:   "edit",
}

// VerifiedUIDUpdate describes update permissions on the user UID
var VerifiedUIDUpdate = resources.PermissionInput{
	Resource: "update_verified_uid",
	Action:   "edit",
}

// VerifiedIdentifiersUpdate describes update permissions on a user verified identifiers
var VerifiedIdentifiersUpdate = resources.PermissionInput{
	Resource: "update_verified_identifiers",
	Action:   "edit",
}

// SuspendedUpdate describes update permissions on a user suspension status
var SuspendedUpdate = resources.PermissionInput{
	Resource: "update_suspended",
	Action:   "edit",
}

// PhotoUploadIDUpdate describes update permissions on a user photo upload
var PhotoUploadIDUpdate = resources.PermissionInput{
	Resource: "update_photo_upload_id",
	Action:   "edit",
}

// PushTokensUpdate describes update permissions on a user push tokens
var PushTokensUpdate = resources.PermissionInput{
	Resource: "update_push_token",
	Action:   "edit",
}

// PermissionsUpdate describes update permissions on a user permissions
var PermissionsUpdate = resources.PermissionInput{
	Resource: "update_permissions",
	Action:   "edit",
}

// BioDataUpdate describes update permissions on a user bio data
var BioDataUpdate = resources.PermissionInput{
	Resource: "update_bio_data",
	Action:   "edit",
}

// PartnerTypeCreate describes write permissions on creating a user partner type
var PartnerTypeCreate = resources.PermissionInput{
	Resource: "add_partner_type",
	Action:   "create",
}

// SupplierCreate describes write permissions on inital creation of a supplier
var SupplierCreate = resources.PermissionInput{
	Resource: "setup_supplier",
	Action:   "create",
}

// CustomerAccountCreate describes write permissions on inital creation of a customer account
var CustomerAccountCreate = resources.PermissionInput{
	Resource: "create_customer_account",
	Action:   "create",
}

// SupplierAccountCreate describes write permissions on inital creation of a supplier account
var SupplierAccountCreate = resources.PermissionInput{
	Resource: "create_supplier_account",
	Action:   "create",
}
