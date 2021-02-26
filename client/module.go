package client

import (
	state "github.com/seashell/agent/client/state"
	log "github.com/seashell/agent/pkg/log"
)

// ModuleService :
type ModuleService struct {
	config *Config
	logger log.Logger
	state  state.Repository
}

// NewModuleService ...
func NewModuleService(config *Config, logger log.Logger, state state.Repository) *ModuleService {
	return &ModuleService{
		config: config,
		logger: logger,
		state:  state,
	}
}
