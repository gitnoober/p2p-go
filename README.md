# P2P Messaging CLI Tool

This project is a simple peer-to-peer (P2P) command-line interface (CLI) tool written in Go. It enables peer discovery within a local network (Wi-Fi) using UDP broadcast messages and supports message exchange between discovered peers.

## Features

- Peer discovery using UDP broadcasts on the local network.
- Ability to send and receive messages between peers.
- Scans and pings all devices connected to the same Wi-Fi to find active peers.
- Supports both IPv4 and IPv6.
- Uses ICMP (Ping) and ARP scanning to identify active devices.
- Easily extendable for additional features such as file transfer, encryption, etc.

## Requirements

- **Go 1.20+**
- **Administrator/root access** (for ICMP-based ping functionality)
- A local Wi-Fi network with multiple connected devices.

## Installation

#### Clone the repository:

```bash
git clone git@github.com:gitnoober/p2p-go.git
cd p2p-go
```

