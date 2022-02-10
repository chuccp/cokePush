package ukm

import (
	log "github.com/chuccp/coke-log"
	"net"
	"testing"
)

func TestStream_udk(t *testing.T) {

	listener, err := ListenerUkm("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 7676})
	if err == nil {
		for {
			conn, err1 := listener.Accept()
			if err1 == nil {
				go func() {
					data := make([]byte, 1024)
					for {
						num, err2 := conn.Read(data)
						log.InfoF("======={}====={}", num, err2)
						if err2 == nil {

						}
					}

				}()
			} else {

			}
		}

	}

}
