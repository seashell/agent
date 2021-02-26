package agent

import (
	"errors"
	"fmt"
	"sync"

	client "github.com/seashell/agent/client"
	log "github.com/seashell/agent/pkg/log"
)

// Agent :
type Agent struct {
	config *Config
	logger log.Logger
	client *client.Client

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// New creates a new Seashell agent from the configuration,
// potentially returning an error
func New(config *Config, logger log.Logger) (*Agent, error) {

	config = DefaultConfig().Merge(config)

	if logger == nil {
		return nil, errors.New("missing logger")
	}

	a := &Agent{
		config:     config,
		logger:     logger.WithName("agent"),
		shutdownCh: make(chan struct{}),
	}

	// Setup Seashell client
	if err := a.setupClient(); err != nil {
		return nil, err
	}

	return a, nil
}

// Shutdown is used to terminate the agent.
func (a *Agent) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()

	if a.shutdown {
		return nil
	}

	a.logger.Infof("requesting shutdown")
	if a.client != nil {
		if err := a.client.Shutdown(); err != nil {
			a.logger.Errorf("client shutdown failed: %s", err.Error())
		}
	}

	a.logger.Infof("agent shutdown complete")

	a.shutdown = true
	close(a.shutdownCh)

	return nil
}

// Setup Seashell client, if enabled
func (a *Agent) setupClient() error {

	config, err := a.clientConfig()
	if err != nil {
		return fmt.Errorf("client config setup failed: %v", err)
	}

	client, err := client.New(config)
	if err != nil {
		return fmt.Errorf("client setup failed: %v", err)
	}

	a.client = client

	return nil
}

// clientConfig creates a new client.Config struct based on an
// agent.Config struct and which can be used to initialize
// a Seashell client
func (a *Agent) clientConfig() (*client.Config, error) {
	c := client.DefaultConfig()

	c.OrganizationID = a.config.Client.OrganizationID
	c.ProjectID = a.config.Client.ProjectID
	c.DeviceBatchID = a.config.Client.BatchID
	c.DeviceID = a.config.Client.DeviceID
	c.DeviceSecret = a.config.Client.SecretID

	c.StateDir = a.config.DataDir
	c.Meta = a.config.Client.Meta

	c.LogLevel = a.config.LogLevel
	c.Logger = a.logger

	return c, nil
}
