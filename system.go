package hs1xxplug

// Reboot restarts the plug after 1 second delay
func (h *Plug) Reboot() error {
	_, err := h.do(RebootCommand, "reboot")
	return err
}
