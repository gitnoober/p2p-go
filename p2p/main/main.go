package main

import (
	"time"

	"github.com/gitnoober/p2p-go/pkg/utils"
)

func main() {
	port := "9999"

	// Start listening to the broadcast messages
	go utils.ListenDiscovery(port)

	time.Sleep(2 * time.Second)

	for {
		utils.BroadcastDiscovery(port)
		time.Sleep(5 * time.Second)
	}
}
