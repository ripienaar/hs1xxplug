package hs1xxplug

import (
	"encoding/json"
)

// Info is the system
type Info struct {
	SoftwareVersion string     `json:"sw_ver"`
	HardwareVersion string     `json:"hw_ver"`
	Type            string     `json:"type"`
	Model           string     `json:"model"`
	MAC             string     `json:"mac"`
	DeviceName      string     `json:"dev_name"`
	Alias           string     `json:"alias"`
	RelayState      PowerState `json:"relay_state"`
	OnTime          int        `json:"on_time"`
	ActiveMode      string     `json:"active_mode"`
	Features        string     `json:"feature"`
	Updating        int        `json:"updating"`
	SignalStrength  int        `json:"rssi"`
	LEDOff          int        `json:"led_off"`
	Lon             int        `json:"longitude_i"`
	Lat             int        `json:"latitude_i"`
	HardwareID      string     `json:"hwId"`
	FirmwareID      string     `json:"fwId"`
	DeviceID        string     `json:"deviceId"`
	OEMID           string     `json:"oemId"`
	NTCState        int        `json:"ntc_state"`
	Error

	On        bool            `json:"power_on"`
	Off       bool            `json:"power_off"`
	Address   string          `json:"address"`
	RawStatus json.RawMessage `json:"-"`
}

// Info retrieves the system information of the plug
func (h *Plug) Info() (info *Info, err error) {
	infoj, err := h.do(InfoCommand, "get_sysinfo")
	if err != nil {
		return nil, err
	}

	info = &Info{
		RawStatus: infoj,
		Address:   h.IPAddress,
	}

	err = json.Unmarshal(infoj, info)
	if err != nil {
		return nil, err
	}

	info.On = info.RelayState == PowerOn
	info.Off = info.RelayState == PowerOff

	return info, err
}
