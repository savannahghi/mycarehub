package postgres

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

// MyCareHubDb struct implements the service's business specific calls to the database
type MyCareHubDb struct {
	create gorm.Create
	query  gorm.Query
	delete gorm.Delete
	update gorm.Update
}

// NewMyCareHubDb initializes a new instance of the MyCareHubDb struct
func NewMyCareHubDb(c gorm.Create, q gorm.Query, d gorm.Delete, u gorm.Update) *MyCareHubDb {
	return &MyCareHubDb{create: c, query: q, delete: d, update: u}
}
