package agent

import (
	"os"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/seashell/agent/version"
)

// Config contains configurations for the Seashell agent
type Config struct {
	// UI defines whether or not Seashell's web UI will be served
	// by the agent
	UI bool `hcl:"ui,optional"`

	// Name is used to identify individual agents
	Name string `hcl:"name,optional"`

	// DataDir is the directory to store our state in
	DataDir string `hcl:"data_dir,optional"`

	// BindAddr is the address on which all of Seashell's services will
	// be bound. If not specified, this defaults to 127.0.0.1.
	BindAddr string `hcl:"bind_addr,optional"`

	// LogLevel is the level of the logs to put out
	LogLevel string `hcl:"log_level,optional"`

	// Ports is used to control the network ports we bind to.
	Ports *Ports `hcl:"ports,block"`

	// Client contains all client-specific configurations
	Client *ClientConfig `hcl:"client,block"`

	// DevMode is set by the --dev CLI flag.
	DevMode bool

	// Version information (set at compilation time)
	Version *version.VersionInfo
}

// Merge merges two Config structs, returning the result
func (c *Config) Merge(b *Config) *Config {

	if b == nil {
		return c
	}

	result := *c

	if b.UI {
		result.UI = true
	}
	if b.LogLevel != "" {
		result.LogLevel = b.LogLevel
	}
	if b.Name != "" {
		result.Name = b.Name
	}
	if b.DataDir != "" {
		result.DataDir = b.DataDir
	}
	if b.BindAddr != "" {
		result.BindAddr = b.BindAddr
	}
	if b.Version != nil {
		result.Version = b.Version
	}

	// Apply the ports config
	if result.Ports == nil && b.Ports != nil {
		ports := *b.Ports
		result.Ports = &ports
	} else if b.Ports != nil {
		result.Ports = result.Ports.Merge(b.Ports)
	}

	// Apply the client config
	if result.Client == nil && b.Client != nil {
		client := *b.Client
		result.Client = &client
	} else if b.Client != nil {
		result.Client = result.Client.Merge(b.Client)
	}

	return &result
}

// ClientConfig contains configurations for the Seashell client
type ClientConfig struct {

	// StateDir is the directory where the client state will be kep
	StateDir string `hcl:"state_dir,optional"`

	// OrganizationID
	OrganizationID string `hcl:"organization_id,optional"`

	// ProjectID
	ProjectID string `hcl:"project_id,optional"`

	// BatchID
	BatchID string `hcl:"device_batch_id,optional"`

	// DeviceID
	DeviceID string `hcl:"device_id,optional"`

	// SecretID
	SecretID string `hcl:"device_secret,optional"`

	// Meta contains metadata about the client node
	Meta map[string]string `hcl:"meta,optional"`

	// SyncInterval controls how frequently the client synchronizes its state
	SyncIntervalSeconds time.Duration `hcl:"sync_interval,optional"`

	// HeartbeatInterval controls how frequently the client issues heartbeats
	HeartbeatIntervalSeconds time.Duration `hcl:"heartbeat_interval,optional"`
}

// Merge merges two ClientConfig structs, returning the result
func (c *ClientConfig) Merge(b *ClientConfig) *ClientConfig {
	result := *c

	if b.StateDir != "" {
		result.StateDir = b.StateDir
	}
	if b.OrganizationID != "" {
		result.OrganizationID = b.OrganizationID
	}
	if b.ProjectID != "" {
		result.ProjectID = b.ProjectID
	}
	if b.BatchID != "" {
		result.BatchID = b.BatchID
	}
	if b.DeviceID != "" {
		result.DeviceID = b.DeviceID
	}
	if b.SecretID != "" {
		result.SecretID = b.SecretID
	}
	if b.SyncIntervalSeconds != 0 {
		result.SyncIntervalSeconds = b.SyncIntervalSeconds
	}
	if b.Meta != nil {
		result.Meta = b.Meta
	}

	return &result
}

// Ports encapsulates the various ports we bind to for network services. If any
// are not specified then the defaults are used instead.
type Ports struct {
	HTTP int `hcl:"http"`
	RPC  int `hcl:"rpc"`
}

// Merge merges two Ports structs, returning the result
func (c *Ports) Merge(b *Ports) *Ports {
	result := *c

	if b.HTTP != 0 {
		result.HTTP = b.HTTP
	}
	if b.RPC != 0 {
		result.RPC = b.RPC
	}

	return &result
}

// DefaultConfig returns a Config struct populated with sane defaults
func DefaultConfig() *Config {
	return &Config{
		LogLevel: "DEBUG",
		UI:       true,
		Name:     "",
		DataDir:  "/tmp/seashell",
		BindAddr: "0.0.0.0",
		Ports: &Ports{
			HTTP: 8123,
			RPC:  8124,
		},
		Client: &ClientConfig{
			Meta:                map[string]string{},
			SyncIntervalSeconds: 5,
		},
		Version: version.GetVersion(),
	}
}

// EmptyConfig returns an empty Config struct with all nested structs
// also initialized to a non-nil empty value.
func EmptyConfig() *Config {
	return &Config{
		Ports:  &Ports{},
		Client: &ClientConfig{},
	}
}

// Validate returns an error in case a Config struct is invalid.
func (c *Config) Validate() error {
	// TODO: implement validation
	return nil
}

// LoadFromFile loads the configuration from a given path
func (c *Config) LoadFromFile(path string) (*Config, error) {

	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = hclsimple.DecodeFile(path, nil, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
