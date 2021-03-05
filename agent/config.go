package agent

import (
	"os"
	"path"
	"time"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/seashell/agent/version"
)

const (
	defaultDataDir = "/tmp/seashell"
	defaultAPIAddr = "https://api.seashell.sh"
)

// Config contains configurations for the Seashell agent
type Config struct {

	// APIAddr contains the address of the API
	APIAddr string `hcl:"api_addr,optional"`

	// Name is used by the agent to identify itself
	Name string `hcl:"name,optional"`

	// DataDir is the directory used by the agent to store its state
	DataDir string `hcl:"data_dir,optional"`

	// LogLevel is the level of the logs to put out
	LogLevel string `hcl:"log_level,optional"`

	// Client contains all client-specific configurations
	Client *ClientConfig `hcl:"client,block"`

	// Version information (set at compilation time)
	Version *version.VersionInfo
}

// Merge merges two Config structs, returning the result
func (c *Config) Merge(b *Config) *Config {

	if b == nil {
		return c
	}

	result := *c

	if b.DataDir != "" {
		result.DataDir = b.DataDir
	}
	if b.APIAddr != "" {
		result.APIAddr = b.APIAddr
	}
	if b.Version != nil {
		result.Version = b.Version
	}
	if b.LogLevel != "" {
		result.LogLevel = b.LogLevel
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

	// APIAddr is the address of the remote API
	APIAddr string

	// StateDir is the directory used by the client to store its state
	StateDir string `hcl:"state_dir,optional"`

	// OutputDir is the directory to which the client renders the configuration
	OutputDir string `hcl:"output_dir,optional"`

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

	if b.APIAddr != "" {
		result.APIAddr = b.APIAddr
	}
	if b.OutputDir != "" {
		result.OutputDir = b.OutputDir
	}
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

// DefaultConfig returns a Config struct populated with sane defaults
func DefaultConfig() *Config {
	return &Config{
		Name:     "",
		APIAddr:  defaultAPIAddr,
		LogLevel: "DEBUG",
		DataDir:  defaultDataDir,
		Client: &ClientConfig{
			APIAddr:             defaultAPIAddr,
			OutputDir:           path.Join(defaultDataDir, "output"),
			SyncIntervalSeconds: 5,
			Meta:                map[string]string{},
		},
		Version: version.GetVersion(),
	}
}

// EmptyConfig returns an empty Config struct with all nested structs
// also initialized to a non-nil empty value.
func EmptyConfig() *Config {
	return &Config{
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
