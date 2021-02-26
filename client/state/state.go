package state

import (
	"context"

	"github.com/seashell/agent/seashell/structs"
)

// Transaction :
type Transaction interface {
	Commit() error
}

// Repository :
type Repository interface {
	Name() string
	Transaction(ctx context.Context) Transaction

	ConfigurationRepository
}

// ConfigurationRepository : Configuration repository interface
type ConfigurationRepository interface {
	Configuration() (*structs.Configuration, error)
	SetConfiguration(*structs.Configuration) error
}
