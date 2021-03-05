package client

import (
	"time"

	log "github.com/seashell/agent/pkg/log"
	version "github.com/seashell/agent/version"
)

const (
	defaultLogLevel  = "DEBUG"
	defaultStateDir  = "/tmp/seashell"
	defaultOutputDir = "/etc/seashell"
)

// Config : Seashell client configuration
type Config struct {

	// APIAddr is the address of the remote API used by the client to fetch configuration
	APIAddr string

	//Logger is the logger the client will use to log.
	Logger log.Logger

	OrganizationID string
	ProjectID      string
	DeviceBatchID  string
	DeviceID       string
	DeviceSecret   string
	DeviceRemoteID string

	// StateDir is the directory used by the client to store its state.
	StateDir string

	// OutputDir is the directory to which the client will render its output.
	OutputDir string

	// ReconcileInterval is the interval between two reconciliation cycles.
	ReconcileInterval time.Duration

	// Meta contains client metadata
	Meta map[string]string

	// Version is the version of the Seashell client.
	Version *version.VersionInfo

	// LogLevel is the level at which the client should output logs.
	LogLevel string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		APIAddr:           "",
		LogLevel:          defaultLogLevel,
		StateDir:          defaultStateDir,
		OutputDir:         defaultOutputDir,
		ReconcileInterval: 5 * time.Second,
		Meta:              map[string]string{},
		Version:           version.GetVersion(),
	}
}

// Merge combines two config structs, returning the result.
func (c *Config) Merge(b *Config) *Config {
	result := *c

	if b.Logger != nil {
		result.Logger = b.Logger
	}
	if b.LogLevel != "" {
		result.LogLevel = b.LogLevel
	}
	if b.OrganizationID != "" {
		result.OrganizationID = b.OrganizationID
	}
	if b.ProjectID != "" {
		result.ProjectID = b.ProjectID
	}
	if b.DeviceBatchID != "" {
		result.DeviceBatchID = b.DeviceBatchID
	}
	if b.DeviceID != "" {
		result.DeviceID = b.DeviceID
	}
	if b.DeviceSecret != "" {
		result.DeviceSecret = b.DeviceSecret
	}
	if b.DeviceRemoteID != "" {
		result.DeviceRemoteID = b.DeviceRemoteID
	}
	if b.APIAddr != "" {
		result.APIAddr = b.APIAddr
	}
	if b.StateDir != "" {
		result.StateDir = b.StateDir
	}
	if b.OutputDir != "" {
		result.OutputDir = b.OutputDir
	}
	if b.ReconcileInterval != 0 {
		result.ReconcileInterval = b.ReconcileInterval
	}
	if b.Meta != nil {
		result.Meta = b.Meta
	}

	return &result
}
