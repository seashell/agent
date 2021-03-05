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
	DragoConfiguration() (*structs.DragoConfiguration, error)
	SetDragoConfiguration(*structs.DragoConfiguration) error
	NomadConfiguration() (*structs.NomadConfiguration, error)
	SetNomadConfiguration(*structs.NomadConfiguration) error
	ConsulConfiguration() (*structs.ConsulConfiguration, error)
	SetConsulConfiguration(*structs.ConsulConfiguration) error
}
