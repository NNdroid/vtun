package udp

import (
	"log"
	"net"

	"github.com/net-byte/vtun/common/cipher"
	"github.com/net-byte/vtun/common/config"
	"github.com/net-byte/vtun/tun"
	"github.com/songgao/water/waterutil"
)

// StartClient starts udp client
func StartClient(config config.Config) {
	iface := tun.CreateTun(config)
	serverAddr, err := net.ResolveUDPAddr("udp", config.ServerAddr)
	if err != nil {
		log.Fatalln("failed to resolve server addr:", err)
	}
	localAddr, err := net.ResolveUDPAddr("udp", config.LocalAddr)
	if err != nil {
		log.Fatalln("failed to get UDP socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln("failed to listen on UDP socket:", err)
	}
	defer conn.Close()
	log.Printf("vtun udp client started on %v,CIDR is %v", config.LocalAddr, config.CIDR)
	// read data from server
	go func() {
		buf := make([]byte, 1500)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil || n == 0 {
				continue
			}
			b := cipher.XOR(buf[:n])
			if !waterutil.IsIPv4(b) {
				continue
			}
			iface.Write(b)
		}
	}()
	// read data from tun
	packet := make([]byte, 1500)
	for {
		n, err := iface.Read(packet)
		if err != nil || n == 0 {
			continue
		}
		if !waterutil.IsIPv4(packet) {
			continue
		}
		b := cipher.XOR(packet[:n])
		conn.WriteToUDP(b, serverAddr)
	}
}
