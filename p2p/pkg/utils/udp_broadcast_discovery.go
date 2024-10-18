package utils

import (
	"fmt"
	"net"
)

// Function to broadcast a message to all the peers in the network
func BroadcastDiscovery(port string) {
	addr, err := net.ResolveUDPAddr("udp", "255.255.255.255:"+port)
	if err != nil {
		fmt.Print("Error while Resolving UDP", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Print("Error while Dialing UDP", err)
		return
	}

	defer conn.Close()

	message := "DISCOVER_PEER"
	_, ferr := conn.Write([]byte(message))
	if ferr != nil {
		fmt.Print("Error while Writing to UDP", ferr)
		return
	}
	fmt.Println("Message sent")
}

// Function to listen to the broadcast messages from the peers
func ListenDiscovery(port string) {
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		fmt.Print("Error while Resolving UDP", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Print("Error while Listening UDP", err)
		return
	}

	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Print("Error while Reading UDP", err)
			continue
		}

		message := string(buf[:n])
		fmt.Println("Received message from", src, "Message:", message)
		if message == "DISCOVER_PEER" {
			sendPeerInfo(src.IP.String(), src.Port)
		}

	}
}

func sendPeerInfo(peerIP string, peerPort int) {
	responseAddr := fmt.Sprintf("%s:%d", peerIP, peerPort)
	conn, err := net.Dial("udp", responseAddr)
	if err != nil {
		fmt.Print("Error while Dialing UDP", err)
		return
	}

	defer conn.Close()

	message := fmt.Sprintf("PEER FOUND with IP: %s and port: %d", getLocalIP(), peerPort)

	_, ferr := conn.Write([]byte(message))

	if ferr != nil {
		fmt.Print("Error while Writing to UDP", ferr)
		return
	}
	fmt.Println("Peer Info sent....")
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Print("Error while getting Interface Address", err)
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			return ipnet.IP.String()
		}
	}
	return ""
}
