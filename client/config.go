package client

import (
	"time"

	log "github.com/seashell/agent/pkg/log"
	version "github.com/seashell/agent/version"
)

const (
	defaultLogLevel = "DEBUG"
	defaultStateDir = "/etc/seashell"
)

// Config : Seashell client configuration
type Config struct {

	// DevMode indicates whether the client is running in development mode.
	DevMode bool

	// Name is used to specify the name of the client node.
	Name string

	// Version is the version of the Seashell client.
	Version *version.VersionInfo

	// LogLevel is the level at which the client should output logs.
	LogLevel string

	//Logger is the logger the client will use to log.
	Logger log.Logger

	OrganizationID string

	ProjectID string

	DeviceBatchID string

	DeviceID string

	DeviceSecret string

	// StateDir is the directory to store our state in.
	StateDir string

	// ReconcileInterval is the interval between two reconciliation cycles.
	ReconcileInterval time.Duration

	// Meta contains client metadata
	Meta map[string]string
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Name:              "",
		Version:           version.GetVersion(),
		LogLevel:          defaultLogLevel,
		StateDir:          defaultStateDir,
		ReconcileInterval: 5 * time.Second,
		Meta:              map[string]string{},
	}
}

// Merge combines two config structs, returning the result.
func (c *Config) Merge(b *Config) *Config {
	result := *c

	if b.Name != "" {
		result.Name = b.Name
	}
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
	if b.StateDir != "" {
		result.StateDir = b.StateDir
	}
	if b.ReconcileInterval != 0 {
		result.ReconcileInterval = b.ReconcileInterval
	}
	if b.Meta != nil {
		result.Meta = b.Meta
	}

	return &result
}
