package hs1xxplug

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/tidwall/gjson"
)

func (h *Plug) encrypt(plaintext string) []byte {
	n := len(plaintext)
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, uint32(n))

	ciphertext := []byte(buf.Bytes())
	key := byte(h.cryptKey)
	payload := make([]byte, n)

	for i := 0; i < n; i++ {
		payload[i] = plaintext[i] ^ key
		key = payload[i]
		ciphertext = append(ciphertext, payload[i])
	}

	return ciphertext
}

func (h *Plug) decrypt(in []byte) string {
	ciphertext := in[4:]
	n := len(ciphertext)
	key := byte(h.cryptKey)
	var nextKey byte

	for i := 0; i < n; i++ {
		nextKey = ciphertext[i]
		ciphertext[i] = ciphertext[i] ^ key
		key = nextKey
	}

	return string(ciphertext)
}

func (h *Plug) do(command string, expected string) (data []byte, err error) {
	ciphertext, err := h.sendRecv(h.encrypt(command))
	if err != nil {
		return nil, err
	}

	body := h.decrypt(ciphertext)
	code, msg := h.extractError(body, expected)
	if code != 0 {
		return nil, fmt.Errorf("error from device: code: %d msg: %s", code, msg)
	}

	system := gjson.Get(body, fmt.Sprintf("system.%s", expected))
	if system.Exists() {
		return []byte(system.String()), nil
	}

	emeter := gjson.Get(body, fmt.Sprintf("emeter.%s", expected))
	if emeter.Exists() {
		return []byte(emeter.String()), nil
	}

	return nil, fmt.Errorf("unknown body received")
}

func (h *Plug) extractError(body string, expected string) (code int64, msg string) {
	code = -1
	msg = ""
	key := ""

	system := gjson.Get(body, "system")
	if system.Exists() {
		key = "system"
	}

	emeter := gjson.Get(body, "emeter")
	if emeter.Exists() {
		key = "emeter"
	}

	errcode := gjson.Get(body, fmt.Sprintf("%s.%s.err_code", key, expected))
	if errcode.Exists() {
		code = errcode.Int()
	}

	errmsg := gjson.Get(body, fmt.Sprintf("%s.%s.err_msg", key, expected))
	if errmsg.Exists() {
		return code, errmsg.String()
	}

	return code, ""
}

func (h *Plug) sendRecv(payload []byte) (data []byte, err error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", h.IPAddress, h.port), h.connTimeout)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %s", err)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(h.readDeadline * time.Second))
	conn.SetWriteDeadline(time.Now().Add(h.writeDeadline * time.Second))

	_, err = conn.Write(payload)
	if err != nil {
		return nil, err
	}

	buff := make([]byte, 4096)
	n, err := conn.Read(buff)
	if err != nil {
		return nil, fmt.Errorf("read failed: %s", err)
	}

	return buff[0:n], nil
}
