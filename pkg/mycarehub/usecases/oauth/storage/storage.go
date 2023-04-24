package storage

import "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"

// Storage represents oauth implementation
type Storage struct {
	Update infrastructure.Update
	Query  infrastructure.Query
	Create infrastructure.Create
	Delete infrastructure.Delete
}

func NewFositeStorage(
	create infrastructure.Create,
	update infrastructure.Update,
	query infrastructure.Query,
	delete infrastructure.Delete,
) Storage {
	return Storage{
		Update: update,
		Query:  query,
		Create: create,
		Delete: delete,
	}
}
