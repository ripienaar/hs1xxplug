package hs1xxplug

import (
	"time"
)

// PowerState represents a state for the plug relay
type PowerState int

const (
	// PowerOnCommand is the command sent to turn the relay on
	PowerOnCommand = `{"system":{"set_relay_state":{"state":1}}}`

	// PowerOffCommand is the command sent to turn the relay off
	PowerOffCommand = `{"system":{"set_relay_state":{"state":0}}}`

	// InfoCommand is the command sent to request system into
	InfoCommand = `{"system":{"get_sysinfo":{}}}`

	// EnergyCommand is the command to retrieve energy info
	EnergyCommand = `{"emeter":{"get_realtime":{}}}`

	// RebootCommand is the command to reboot the plug
	RebootCommand = `{"system":{"reboot":{"delay":1}}}`

	// PowerUnknown represents an unknown power state
	PowerUnknown PowerState = -1

	// PowerOff represents the off state off the plug
	PowerOff PowerState = 0

	// PowerOn represents the off state on the plug
	PowerOn PowerState = 1
)

// Plug represents a management interface for a plug
type Plug struct {
	IPAddress string

	port          int
	cryptKey      byte
	connTimeout   time.Duration
	writeDeadline time.Duration
	readDeadline  time.Duration
}

// NewPlug creates a new management interface for the TP Link HS1xx plug
func NewPlug(ip string) *Plug {
	return &Plug{
		IPAddress:     ip,
		port:          9999,
		cryptKey:      byte(0xAB),
		connTimeout:   10 * time.Second,
		writeDeadline: 2,
		readDeadline:  2,
	}
}
