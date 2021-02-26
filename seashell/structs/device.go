package structs

import (
	"fmt"
	"time"
)

const (
	DeviceStatusInit  = "initializing"
	DeviceStatusReady = "ready"
	DeviceStatusDown  = "down"
)

// Device :
type Device struct {
	ID        string
	Name      string
	Secret    string
	Status    string
	Token     string
	Meta      map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate validates a structs.Device object
func (n *Device) Validate() error {

	if !IsValidDeviceStatus(n.Status) {
		return fmt.Errorf("invalid device status")
	}

	return nil
}

// IsValidDeviceStatus returns true if the status passed as argument
// corresponds to a valid device status. Otherwise returns false.
func IsValidDeviceStatus(s string) bool {

	valid := map[string]interface{}{
		DeviceStatusInit:  nil,
		DeviceStatusReady: nil,
		DeviceStatusDown:  nil,
	}

	if _, ok := valid[s]; !ok {
		return false
	}

	return true
}

// Merge :
func (n *Device) Merge(in *Device) *Device {

	result := *n

	if in.ID != "" {
		result.ID = in.ID
	}
	if in.Secret != "" {
		result.Secret = in.Secret
	}
	if in.Name != "" {
		result.Name = in.Name
	}
	if in.Meta != nil {
		for k, v := range in.Meta {
			result.Meta[k] = v
		}
	}
	if in.Status != "" {
		result.Status = in.Status
	}

	return &result
}

// Stub :
func (n *Device) Stub() *DeviceListStub {
	return &DeviceListStub{
		ID:        n.ID,
		Name:      n.Name,
		Status:    n.Status,
		Meta:      n.Meta,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

// DeviceListStub :
type DeviceListStub struct {
	ID        string
	Name      string
	Status    string
	Meta      map[string]string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DeviceGetTokenRequest :
type DeviceGetTokenRequest struct {
	OrganizationID string
	ProjectID      string
	BatchID        string
	DeviceID       string
	SecretID       string

	QueryOptions
}

// DeviceTokenResponse :
type DeviceTokenResponse struct {
	Token string `json:"token"`
	Response
}

// DeviceSyncRequest :
type DeviceSyncRequest struct {
	OrganizationID string
	ProjectID      string
	BatchID        string
	DeviceID       string
	DeviceRemoteID string

	QueryOptions
}

// DeviceSyncResponse :
type DeviceSyncResponse struct {
	*Configuration

	Response
}

// Configuration :
type Configuration struct {
	Labels           map[string]string `json:"labels"`
	DragoIPAddresses []string          `json:"dragoIpAddresses"`
}
