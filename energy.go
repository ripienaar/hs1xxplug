package hs1xxplug

import (
	"encoding/json"
)

// Energy is the energy usage information for a plug
type Energy struct {
	MilliVolt         int     `json:"voltage_mv"`
	Volt              float32 `json:"volt"`
	MilliAmp          int     `json:"current_ma"`
	Amp               float32 `json:"current_amp"`
	PowerUseMilliWatt int     `json:"power_mw"`
	PowerUseWatt      float32 `json:"power_w"`
	TotalMilliWatt    int     `json:"total_wh"`
	TotalWatt         float32 `json:"total_watt"`

	Error
}

// Energy retrieves the energy information
func (h *Plug) Energy() (*Energy, error) {
	energyj, err := h.do(EnergyCommand, "get_realtime")
	if err != nil {
		return nil, err
	}

	energy := &Energy{}
	err = json.Unmarshal(energyj, energy)
	if err != nil {
		return nil, err
	}

	energy.Volt = float32(energy.MilliVolt) / 1000
	energy.Amp = float32(energy.MilliAmp) / 1000
	energy.PowerUseWatt = float32(energy.PowerUseMilliWatt) / 1000
	energy.TotalWatt = float32(energy.TotalMilliWatt) / 1000

	return energy, nil
}
