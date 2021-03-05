package client

import (
	"context"
	_ "embed"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	api "github.com/seashell/agent/api"
	state "github.com/seashell/agent/client/state"
	boltdb "github.com/seashell/agent/client/state/boltdb"
	log "github.com/seashell/agent/pkg/log"
	structs "github.com/seashell/agent/seashell/structs"
)

//go:embed "assets/drago.hcl.tmpl"
var dragoTemplateString string

//go:embed "assets/nomad.hcl.tmpl"
var nomadTemplateString string

//go:embed "assets/consul.hcl.tmpl"
var consulTemplateString string

var (
	defaultAuthenticationRetryInterval = 2 * time.Second
	defaultReconciliationRetryInterval = 5 * time.Second
	defaultReconciliationInterval      = 2 * time.Second
	defaultFirstHeartbeatDelay         = 1 * time.Second
	defaultHeartbeatInterval           = 1 * time.Second
)

// Client is the Seashell client
type Client struct {
	config *Config

	logger log.Logger

	api *api.Client

	state state.Repository

	device     *structs.Device
	deviceLock sync.Mutex

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
}

// New is used to create a new Seashell client from the
// configuration, potentially returning an error
func New(config *Config) (*Client, error) {

	rand.Seed(time.Now().Unix())

	config = DefaultConfig().Merge(config)

	c := &Client{
		config:     config,
		logger:     config.Logger.WithName("client"),
		shutdownCh: make(chan struct{}),
	}

	err := c.setupState()
	if err != nil {
		return nil, fmt.Errorf("error setting up client state: %v", err)
	}

	err = c.setupOutputDir()
	if err != nil {
		return nil, fmt.Errorf("error setting up output dir: %v", err)
	}

	err = c.setupDevice()
	if err != nil {
		return nil, fmt.Errorf("error setting up device: %v", err)
	}

	err = c.setupAPIClient()
	if err != nil {
		return nil, fmt.Errorf("error setting up api client: %v", err)
	}

	// Try to get a token
	c.tryToGetTokenUntilSuccessful()
	c.logger.Infof("successfully obtained auth token!")

	// Start goroutine for reconciling the client state
	go c.run()

	c.logger.Infof("started device %s", c.DeviceID())

	return c, nil
}

// Device returns the device managed by this client
func (c *Client) Device() *structs.Device {
	return c.device
}

// DeviceID returns the device ID
func (c *Client) DeviceID() string {
	return c.Device().ID
}

// DeviceSecret returns the device secret
func (c *Client) DeviceSecret() string {
	return c.device.Secret
}

func (c *Client) setupDevice() error {

	if c.device == nil {
		c.device = &structs.Device{}
	}

	if c.config.DeviceID == "" {
		return fmt.Errorf("invalid device ID")
	}

	if c.config.DeviceSecret == "" {
		return fmt.Errorf("invalid device secret")
	}

	c.device.ID = c.config.DeviceID
	c.device.Secret = c.config.DeviceSecret

	c.device.Meta = c.config.Meta

	c.device.Status = structs.DeviceStatusInit

	if c.device.Name == "" {
		if hostname, _ := os.Hostname(); hostname != "" {
			c.device.Name = hostname
		} else {
			c.device.Name = c.device.ID
		}
	}
	if c.device.Meta == nil {
		c.device.Meta = make(map[string]string)
	}

	return nil
}

func (c *Client) setupState() error {

	// Ensure the state dir exists. If it was not was specified,
	// create a temporary directory to store the client state.
	if c.config.StateDir != "" {
		if err := os.MkdirAll(c.config.StateDir, 0700); err != nil {
			return fmt.Errorf("failed to create state dir: %s", err)
		}
	} else {
		tmp, err := c.createTempDir("SeashellClient")
		if err != nil {
			return fmt.Errorf("failed to create tmp dir for storing state: %s", err)
		}
		c.config.StateDir = tmp
	}

	c.logger.Infof("using state directory %s", c.config.StateDir)

	repo := boltdb.NewStateRepository(path.Join(c.config.StateDir, "client.state"), c.logger)

	c.state = repo

	return nil
}

func (c *Client) setupOutputDir() error {

	// Ensure the output dir exists. If it was not was specified,
	// create a temporary directory to store the client output.
	if c.config.OutputDir != "" {
		if err := os.MkdirAll(c.config.OutputDir, 0700); err != nil {
			return fmt.Errorf("failed to create output dir: %s", err)
		}
	} else {
		tmp, err := c.createTempDir("SeashellClientOutput")
		if err != nil {
			return fmt.Errorf("failed to create tmp dir for storing output: %s", err)
		}
		c.config.OutputDir = tmp
	}

	c.logger.Infof("using output directory %s", c.config.OutputDir)

	return nil
}

func (c *Client) setupAPIClient() error {

	apiClient, err := api.NewClient(&api.Config{
		Address: c.config.APIAddr,
	})
	if err != nil {
		return err
	}

	c.api = apiClient

	return nil
}

func (c *Client) run() {

	c.logger.Debugf("running client")

	configurationUpdateCh := make(chan *structs.Configuration)
	go c.watchConfiguration(configurationUpdateCh)

	for {
		select {
		case desired := <-configurationUpdateCh:
			c.shutdownLock.Lock()
			if c.shutdown {
				c.shutdownLock.Unlock()
				return
			}

			c.reconcileConfiguration(desired)

			c.shutdownLock.Unlock()
		case <-c.shutdownCh:
			return
		}
	}
}

func (c *Client) reconcileConfiguration(desired *structs.Configuration) {

	c.logger.Debugf("reconciliation started...")

	if err := c.reconcileDragoConfiguration(desired); err != nil {
		c.logger.Warnf("error reconciling drago configuration : %v", err)
	}

	if err := c.reconcileNomadConfiguration(desired); err != nil {
		c.logger.Warnf("error reconciling nomad configuration : %v", err)
	}

	if err := c.reconcileConsulConfiguration(desired); err != nil {
		c.logger.Warnf("error reconciling consul configuration : %v", err)
	}

}

func (c *Client) reconcileDragoConfiguration(config *structs.Configuration) error {

	desired := &structs.DragoConfiguration{
		Name:    c.config.DeviceRemoteID,
		DataDir: path.Join(c.config.StateDir, "drago"),
		Servers: config.DragoIPAddresses,
		Secret:  config.DragoSecret,
		Meta:    config.Labels,
	}

	current, err := c.state.DragoConfiguration()
	if err != nil {
		c.logger.Errorf("could not read drago configuration: %v", err)
	}

	if current.Hash() != desired.Hash() {

		c.logger.Debugf("changes detected in drago configuration. rendering template and persisting to repository...")

		if err := renderTemplateToFile(dragoTemplateString, path.Join(c.config.OutputDir, "drago.hcl"), desired); err != nil {
			return err
		}

		// We only persist configurations that were successfully rendered so as
		// to ensure the state in the DB is synced with the configuration files.
		if err := c.state.SetDragoConfiguration(desired); err != nil {
			return err
		}

		return nil
	}

	c.logger.Debugf("no changes detected in drago configuration. skipping reconciliation...")

	return nil
}

func (c *Client) reconcileNomadConfiguration(config *structs.Configuration) error {

	desired := &structs.NomadConfiguration{
		Name:      c.config.DeviceRemoteID,
		DataDir:   path.Join(c.config.StateDir, "nomad"),
		RetryJoin: "", // TODO: get from Drago and, in the future, replace using go-connect
		Meta:      config.Labels,
	}

	current, err := c.state.NomadConfiguration()
	if err != nil {
		c.logger.Errorf("could not read nomad configuration: %v", err)
	}

	if current.Hash() != desired.Hash() {

		c.logger.Debugf("changes detected in nomad configuration. rendering template and persisting to repository...")

		if err := renderTemplateToFile(nomadTemplateString, path.Join(c.config.OutputDir, "nomad.hcl"), desired); err != nil {
			return err
		}

		// We only persist configurations that were successfully rendered so as
		// to ensure the state in the DB is synced with the configuration files.
		if err := c.state.SetNomadConfiguration(desired); err != nil {
			return err
		}

		return nil
	}

	c.logger.Debugf("no changes detected in nomad configuration. skipping reconciliation...")

	return nil
}

func (c *Client) reconcileConsulConfiguration(config *structs.Configuration) error {

	desired := &structs.ConsulConfiguration{
		Name:      c.config.DeviceRemoteID,
		DataDir:   path.Join(c.config.StateDir, "consul"),
		RetryJoin: "", // TODO: get from Drago and, in the future, replace using go-connect
		Meta:      config.Labels,
	}

	current, err := c.state.ConsulConfiguration()
	if err != nil {
		c.logger.Errorf("could not read consul configuration: %v", err)
	}

	if current.Hash() != desired.Hash() {

		c.logger.Debugf("changes detected in consul configuration. rendering template and persisting to repository...")

		if err := renderTemplateToFile(consulTemplateString, path.Join(c.config.OutputDir, "consul.hcl"), desired); err != nil {
			return err
		}

		if err := c.state.SetConsulConfiguration(desired); err != nil {
			return err
		}

		return nil
	}

	c.logger.Debugf("no changes detected in consul configuration. skipping reconciliation...")

	return nil
}

func (c *Client) watchConfiguration(ch chan *structs.Configuration) {

	c.logger.Debugf("watching configuration")

	for {

		var err error
		var resp *structs.DeviceSyncResponse

		req := &structs.DeviceSyncRequest{
			OrganizationID: c.config.OrganizationID,
			ProjectID:      c.config.ProjectID,
			BatchID:        c.config.DeviceBatchID,
			DeviceID:       c.config.DeviceID,
			DeviceRemoteID: c.config.DeviceRemoteID,
		}

		req.QueryOptions.AuthToken = c.Device().Token

		ctx := context.TODO()

		if resp, err = c.api.Devices().SyncDevice(ctx, req); err != nil {
			c.logger.Debugf("error syncing device: %v", err)

			c.tryToGetTokenUntilSuccessful()

			retryCh := time.After(randomDuration(defaultReconciliationRetryInterval, 1*time.Second))
			select {
			case <-retryCh:
			case <-c.shutdownCh:
				return
			}

		} else {
			ch <- resp.Configuration
		}

		retryCh := time.After(randomDuration(c.config.ReconcileInterval, 1*time.Second))
		select {
		case <-c.shutdownCh:
			return
		case <-retryCh:
		}
	}
}

func (c *Client) tryToGetTokenUntilSuccessful() {

	for {
		select {
		case <-c.shutdownCh:
			return
		default:
		}

		var err error
		var resp *structs.DeviceTokenResponse

		req := &structs.DeviceGetTokenRequest{
			OrganizationID: c.config.OrganizationID,
			ProjectID:      c.config.ProjectID,
			BatchID:        c.config.DeviceBatchID,
			DeviceID:       c.config.DeviceID,
			SecretID:       c.config.DeviceSecret,
		}

		ctx := context.TODO()

		if resp, err = c.api.Devices().GetDeviceToken(ctx, req); err == nil {
			c.deviceLock.Lock()
			c.device.Token = resp.Token
			c.deviceLock.Unlock()
			return
		}

		c.logger.Debugf("error obtaining device token: %v", err)

		retryCh := time.After(randomDuration(defaultAuthenticationRetryInterval, 1*time.Second))

		select {
		case <-retryCh:
		case <-c.shutdownCh:
			return
		}
	}
}

// Shutdown is used to tear down the client
func (c *Client) Shutdown() error {
	c.shutdownLock.Lock()
	defer c.shutdownLock.Unlock()

	if c.shutdown {
		c.logger.Infof("client already shutdown")
		return nil
	}
	c.logger.Infof("shutting down")

	c.shutdown = true
	close(c.shutdownCh)

	return nil
}

func (c *Client) createTempDir(pattern string) (string, error) {
	p, err := ioutil.TempDir("", pattern)
	if err != nil {
		return "", fmt.Errorf("could not create temporary directory: %v", err)
	}
	p, err = filepath.EvalSymlinks(p)
	if err != nil {
		return "", fmt.Errorf("could not retrieve path to StateDir: %v", err)
	}
	return p, nil
}

// Generates a random duration in the interval [mean-delta, mean+delta]
func randomDuration(mean time.Duration, delta time.Duration) time.Duration {
	t := mean.Milliseconds() + int64((rand.Float32()-0.5)*float32(delta.Milliseconds()))
	return time.Duration(t * int64(time.Millisecond))
}
