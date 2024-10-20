package utils

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/go-ping/ping"
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

// Function to get the local IP address and subnet mask
func getLocalIPWithSubnet() (string, *net.IPNet, error) {
	addrs, err := net.Interfaces()
	if err != nil {
		return "", nil, err
	}

	for _, iface := range addrs {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String(), ipnet, nil
					}
				}
			}
		}
	}
	return "", nil, fmt.Errorf("no network connection found")
}

func pingIP(ip string, wg *sync.WaitGroup, aliveIps *[]string) {
	fmt.Println("Pinging IP:", ip)
	defer wg.Done()
	pinger, err := ping.NewPinger(ip)
	// fmt.Println("Pinger created", err)
	if err != nil {
		fmt.Print("Error while creating Pinger", err)
		return
	}

	pinger.Debug = true
	// fmt.Println("Pinger created", pinger.Debug)

	pinger.Count = 3
	pinger.Timeout = time.Second * 5
	pinger.Interval = time.Microsecond * 500
	pinger.SetPrivileged(true)
	pinger.Run()
	stats := pinger.Statistics()
	fmt.Println("Stats:", stats)
	if stats.PacketsRecv > 0 {
		fmt.Printf("IP %s is alive\n", ip)
		*aliveIps = append(*aliveIps, ip)
	}
}

func ScanNetwork() {
	fmt.Println("Scanning network for alive IPs")
	localIP, ipnet, err := getLocalIPWithSubnet()
	if err != nil {
		fmt.Println("Error while getting Local IP", err)
		return

	}

	fmt.Println("Local IP:", localIP)
	fmt.Println("Subnet Mask:", ipnet.Mask)
	fmt.Println("Scanning network range:", ipnet.IP.String())

	var wg sync.WaitGroup
	var aliveIps []string

	ip := ipnet.IP.Mask(ipnet.Mask)
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		wg.Add(1)
		go pingIP(ip.String(), &wg, &aliveIps)
	}

	fmt.Println("Waiting for all the pings to complete")
	fmt.Println("Alive IPs:", aliveIps)
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
