package hs1xxplug

import (
	"reflect"
	"time"

	"errors"

	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"testing"
)

func TestPowerOff(t *testing.T) {
	testCases := []struct {
		test                   string
		plugRelayStateFeedback int
		want                   error
	}{
		{"PowerOff succeed to turn the power off", 0, nil},
		{"PowerOff fails to turn the power off", 1, errors.New("power off was requested but device stayed on")},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			randPort := rand.Intn(9999) // This reduce the risks of collisions on parallel runs
			p := &Plug{
				IPAddress:     "localhost",
				port:          randPort,
				cryptKey:      byte(0xAB),
				connTimeout:   10 * time.Second,
				writeDeadline: 2,
				readDeadline:  2,
			}
			listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", randPort))
			if err != nil {
				log.Fatal(err)
			}
			defer listener.Close()

			go func(l net.Listener) {
				// Handler for the test
				for {
					conn, err := l.Accept()
					if err != nil || conn == nil {
						break
					}
					go func(c net.Conn) {
						buf := make([]byte, 4096)
						reqLen, err := c.Read(buf)
						if err != nil {
							if err != io.EOF {
								fmt.Println("Error:" + err.Error())
							}
						}

						// fmt.Printf("Message contents: %s\n", p.decrypt(buf[:reqLen]))

						rw := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
						switch {
						// TurnOn command
						case reflect.DeepEqual(buf[:reqLen], p.encrypt(PowerOffCommand)):
							rw.Write(p.encrypt("{\"system\":{\"set_relay_state\":{\"err_code\":0}}}"))
						default:
							rw.Write(p.encrypt(
								fmt.Sprintf("{\"system\":{\"get_sysinfo\":{\"sw_ver\":\"1.5.3 Build 180619 Rel.094609\",\"hw_ver\":\"2.0\",\"type\":\"IOT.SMARTPLUGSWITCH\",\"model\":\"HS110(EN)\",\"mac\":\"86:DA:C4:BE:E8:C6\",\"dev_name\":\"Smart Wi-Fi Plug With Energy Monitoring\",\"alias\":\"thylong\",\"relay_state\":%d,\"on_time\":0,\"active_mode\":\"none\",\"feature\":\"TIM:ENE\",\"updating\":0,\"icon_hash\":\"\",\"rssi\":-30,\"led_off\":0,\"longitude_i\":22856,\"latitude_i\":488207,\"hwId\":\"AAA2979A05ERDD4ED8509CA93007EC550\",\"fwId\":\"00000000000000000000000000000000\",\"deviceId\":\"800692DD412DE6D91FC8EA3A3E55827D1BBFF79E\",\"oemId\":\"CF482E3482754DDAAA09BA72423467CE\",\"next_action\":{\"type\":-1},\"err_code\":0}}}", tc.plugRelayStateFeedback),
							))
						}

						rw.Flush()
						c.Close()
					}(conn)
				}
			}(listener)

			err = p.PowerOff()

			listener.Close()
			if err != nil && tc.want != nil && (tc.want.Error() != err.Error()) {
				t.Errorf("unexpected error returned, expecting \"%s\" instead of \"%s\"", tc.want, err)
			}
		})
	}
}

func TestPowerOn(t *testing.T) {
	testCases := []struct {
		test                   string
		plugRelayStateFeedback int
		want                   error
	}{
		{"PowerOn succeed to turn the power on", 1, nil},
		{"PowerOn fails to turn the power on", 0, errors.New("power on was requested but device stayed off")},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			randPort := rand.Intn(9999) // This reduce the risks of collisions on parallel runs
			p := &Plug{
				IPAddress:     "localhost",
				port:          randPort,
				cryptKey:      byte(0xAB),
				connTimeout:   10 * time.Second,
				writeDeadline: 2,
				readDeadline:  2,
			}
			listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", randPort))
			if err != nil {
				log.Fatal(err)
			}
			defer listener.Close()

			go func(l net.Listener) {
				p := NewPlug("localhost")

				// Handler for the test
				for {
					conn, err := l.Accept()
					if err != nil || conn == nil {
						break
					}
					go func(c net.Conn) {
						buf := make([]byte, 4096)
						reqLen, err := c.Read(buf)
						if err != nil {
							if err != io.EOF {
								fmt.Println("Error:" + err.Error())
							}
						}

						rw := bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
						switch {
						// TurnOn command
						case reflect.DeepEqual(buf[:reqLen], p.encrypt(PowerOnCommand)):
							rw.Write(p.encrypt("{\"system\":{\"set_relay_state\":{\"err_code\":0}}}"))
						default:
							rw.Write(p.encrypt(
								fmt.Sprintf("{\"system\":{\"get_sysinfo\":{\"sw_ver\":\"1.5.3 Build 180619 Rel.094609\",\"hw_ver\":\"2.0\",\"type\":\"IOT.SMARTPLUGSWITCH\",\"model\":\"HS110(EN)\",\"mac\":\"86:DA:C4:BE:E8:C6\",\"dev_name\":\"Smart Wi-Fi Plug With Energy Monitoring\",\"alias\":\"thylong\",\"relay_state\":%d,\"on_time\":0,\"active_mode\":\"none\",\"feature\":\"TIM:ENE\",\"updating\":0,\"icon_hash\":\"\",\"rssi\":-30,\"led_off\":0,\"longitude_i\":22856,\"latitude_i\":488207,\"hwId\":\"AAA2979A05ERDD4ED8509CA93007EC550\",\"fwId\":\"00000000000000000000000000000000\",\"deviceId\":\"800692DD412DE6D91FC8EA3A3E55827D1BBFF79E\",\"oemId\":\"CF482E3482754DDAAA09BA72423467CE\",\"next_action\":{\"type\":-1},\"err_code\":0}}}", tc.plugRelayStateFeedback),
							))
						}

						rw.Flush()
						c.Close()
					}(conn)
				}
			}(listener)

			err = p.PowerOn()

			// listener.Close()
			if err != nil && tc.want != nil && (tc.want.Error() != err.Error()) {
				t.Errorf("unexpected error returned, expecting \"%s\" instead of \"%s\"", tc.want, err)
			}
		})
	}
}
