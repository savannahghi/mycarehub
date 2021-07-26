package permission

import (
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
)

// UserProfileView describes read permissions on viewing a user profile
var UserProfileView = dto.PermissionInput{
	Resource: "user_profile_view",
	Action:   "view",
}

// PrimaryPhoneUpdate describes update permissions on a user primary phonenumber
var PrimaryPhoneUpdate = dto.PermissionInput{
	Resource: "update_primary_phone",
	Action:   "edit",
}

// PrimaryEmailUpdate describes update permissions on a user's primary email address
var PrimaryEmailUpdate = dto.PermissionInput{
	Resource: "update_primary_email",
	Action:   "edit",
}

// SecondaryPhoneNumberUpdate describes update permissions on the user secondary phonenumber
var SecondaryPhoneNumberUpdate = dto.PermissionInput{
	Resource: "update_secondary_phone",
	Action:   "edit",
}

// SecondaryEmailAddressUpdate describes update permissions on the user secondary email address
var SecondaryEmailAddressUpdate = dto.PermissionInput{
	Resource: "update_secondary_email",
	Action:   "edit",
}

// VerifiedUIDUpdate describes update permissions on the user UID
var VerifiedUIDUpdate = dto.PermissionInput{
	Resource: "update_verified_uid",
	Action:   "edit",
}

// VerifiedIdentifiersUpdate describes update permissions on a user verified identifiers
var VerifiedIdentifiersUpdate = dto.PermissionInput{
	Resource: "update_verified_identifiers",
	Action:   "edit",
}

// SuspendedUpdate describes update permissions on a user suspension status
var SuspendedUpdate = dto.PermissionInput{
	Resource: "update_suspended",
	Action:   "edit",
}

// PhotoUploadIDUpdate describes update permissions on a user photo upload
var PhotoUploadIDUpdate = dto.PermissionInput{
	Resource: "update_photo_upload_id",
	Action:   "edit",
}

// PushTokensUpdate describes update permissions on a user push tokens
var PushTokensUpdate = dto.PermissionInput{
	Resource: "update_push_token",
	Action:   "edit",
}

// PermissionsUpdate describes update permissions on a user permissions
var PermissionsUpdate = dto.PermissionInput{
	Resource: "update_permissions",
	Action:   "edit",
}

// BioDataUpdate describes update permissions on a user bio data
var BioDataUpdate = dto.PermissionInput{
	Resource: "update_bio_data",
	Action:   "edit",
}

// PartnerTypeCreate describes write permissions on creating a user partner type
var PartnerTypeCreate = dto.PermissionInput{
	Resource: "add_partner_type",
	Action:   "create",
}

// SupplierCreate describes write permissions on initial creation of a supplier
var SupplierCreate = dto.PermissionInput{
	Resource: "setup_supplier",
	Action:   "create",
}

// CustomerAccountCreate describes write permissions on initial creation of a customer account
var CustomerAccountCreate = dto.PermissionInput{
	Resource: "create_customer_account",
	Action:   "create",
}

// SupplierAccountCreate describes write permissions on initial creation of a supplier account
var SupplierAccountCreate = dto.PermissionInput{
	Resource: "create_supplier_account",
	Action:   "create",
}

// MicroserviceDelete describes delete permissions on a microservice
var MicroserviceDelete = dto.PermissionInput{
	Resource: "microservice_delete",
	Action:   "delete",
}

// MicroserviceCreate describes create permissions on a microservice
var MicroserviceCreate = dto.PermissionInput{
	Resource: "microservice_create",
	Action:   "create",
}
