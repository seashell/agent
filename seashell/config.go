package seashell

import (
	"time"

	log "github.com/seashell/agent/pkg/log"
	version "github.com/seashell/agent/version"
)

const (
	defaultDataDir     = "/tmp/seashell"
	defaultBindAddress = "0.0.0.0"
	defaultLogLevel    = "DEBUG"
	defaultHTTPPort    = 8080
	defaultRPCPort     = 8081
)

// Config : Seashell server configuration.
type Config struct {
	// UI enabled.
	UI bool

	// Version is the version of the Seashell server
	Version *version.VersionInfo

	// LogLevel is the level at which the server should output logs
	LogLevel string

	//Logger.
	Logger log.Logger

	// BindAddr.
	BindAddr string

	// RPCAdvertiseAddr is the address advertised to client nodes.
	RPCAdvertiseAddr string

	// DataDir is the directory to store our state in.
	DataDir string

	// Ports.
	Ports *Ports

	// HostGCInterval is how often we perform garbage collection of hosts.
	HostGCInterval time.Duration
}

// Ports :
type Ports struct {
	HTTP int
	RPC  int
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		UI:       true,
		Version:  version.GetVersion(),
		LogLevel: defaultLogLevel,
		BindAddr: defaultBindAddress,
		DataDir:  defaultDataDir,
		Ports: &Ports{
			HTTP: defaultHTTPPort,
			RPC:  defaultRPCPort,
		},
		HostGCInterval: 5 * time.Minute,
	}
}
