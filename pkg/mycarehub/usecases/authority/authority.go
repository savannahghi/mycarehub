package authority

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification"
)

// UsecaseAuthority groups al the interfaces for the Authority usecase
type UsecaseAuthority interface {
}

// UsecaseAuthorityImpl represents the Authority implementation
type UsecaseAuthorityImpl struct {
	Query        infrastructure.Query
	Update       infrastructure.Update
	ExternalExt  extension.ExternalMethodsExtension
	Notification notification.UseCaseNotification
}

// NewUsecaseAuthority is the controller function for the Authority usecase
func NewUsecaseAuthority(
	query infrastructure.Query,
	update infrastructure.Update,
	externalExt extension.ExternalMethodsExtension,
	notification notification.UseCaseNotification,
) *UsecaseAuthorityImpl {
	return &UsecaseAuthorityImpl{
		Query:        query,
		Update:       update,
		ExternalExt:  externalExt,
		Notification: notification,
	}
}
