package hs1xxplug

import (
	"fmt"
)

// PowerOff turns the plug off
func (h *Plug) PowerOff() error {
	_, err := h.do(PowerOffCommand, "set_relay_state")
	if err != nil {
		return err
	}

	state, err := h.PowerState()
	if err != nil {
		return err
	}

	if state != PowerOff {
		return fmt.Errorf("power off was requested but device stayed on")
	}

	return nil
}

// PowerOn turns the plug on
func (h *Plug) PowerOn() error {
	_, err := h.do(PowerOnCommand, "set_relay_state")
	if err != nil {
		return err
	}

	state, err := h.PowerState()
	if err != nil {
		return err
	}

	if state != PowerOn {
		return fmt.Errorf("power on was requested but device stayed off")
	}

	return err
}

// PowerState retrieves the current power state of the plug, PowerUnknown when request failed
func (h *Plug) PowerState() (PowerState, error) {
	state, err := h.Info()
	if err != nil {
		return PowerUnknown, fmt.Errorf("could not determine if plug was turned on: %s", err)
	}

	return state.RelayState, nil
}
