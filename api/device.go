package api

import (
	"context"
	"fmt"

	"github.com/seashell/agent/seashell/structs"
)

const (
	devicesPath = "/devicesync"
)

// Devices is a handle to the devices API
type Devices struct {
	client *Client
}

// Devices returns a handle on the nodes endpoints.
func (c *Client) Devices() *Devices {
	return &Devices{client: c}
}

// GetDeviceToken :
func (d *Devices) GetDeviceToken(ctx context.Context, req *structs.DeviceGetTokenRequest) (*structs.DeviceTokenResponse, error) {

	var resp structs.DeviceTokenResponse

	c := d.client.WithHeaders(map[string]string{
		"X-Organization-ID": req.OrganizationID,
		"X-Project-ID":      req.ProjectID,
		"X-Device-Batch-ID": req.BatchID,
		"X-Device-ID":       req.DeviceID,
		"X-Device-Secret":   req.SecretID,
	})

	err := c.get(devicesPath, "token", &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// SyncDevice :
func (d *Devices) SyncDevice(ctx context.Context, req *structs.DeviceSyncRequest) (*structs.DeviceSyncResponse, error) {

	var resp structs.DeviceSyncResponse

	c := d.client.WithHeaders(map[string]string{
		"X-Organization-ID":  req.OrganizationID,
		"X-Project-ID":       req.ProjectID,
		"X-Device-Batch-ID":  req.BatchID,
		"X-Device-ID":        req.DeviceID,
		"Authorization":      fmt.Sprintf("Bearer %s", req.AuthToken),
		"X-Device-Remote-ID": req.DeviceRemoteID,
	})

	err := c.get(devicesPath, "sync", &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
