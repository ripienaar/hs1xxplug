package hs1xxplug

// Error represents an error message from the device, Code 0 is success
type Error struct {
	Code    int    `json:"err_code"`
	Message string `json:"err_msg"`
}
