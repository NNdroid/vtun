# vtun

A simple VPN written in golang.

[![Travis](https://travis-ci.com/net-byte/vtun.svg?branch=master)](https://github.com/net-byte/vtun)
[![Go Report Card](https://goreportcard.com/badge/github.com/net-byte/vtun)](https://goreportcard.com/report/github.com/net-byte/vtun)
![image](https://img.shields.io/badge/License-MIT-orange)
![image](https://img.shields.io/badge/License-Anti--996-red)

# Features
* Support tun over UDP
* Support tun over WebSocket
* Support data encryption

# Usage  

```
Usage of ./vtun:
  -S    server mode
  -t    enable tls
  -c string
        tun interface CIDR (default "172.16.0.1/24")
  -k string
        encryption key (default "6w9z$C&F)J@NcRfWjXn3r4u7x!A%D*G-")
  -p string
        protocol ws/udp (default "ws")
  -l string
        local address (default "0.0.0.0:3000")
  -s string
        server address (default "0.0.0.0:3001")
```

## Build

```
sh scripts/build.sh 
```

## Client

```
sudo ./vtun-linux-amd64 -l=:3000 -s=server-addr:3001 -c=172.16.0.10/24 -k=123456

```

## Server

```
sudo ./vtun-linux-amd64 -S -l=:3001 -c=172.16.0.1/24 -k=123456

```

## Server setup

1. Add TLS for websocket,reverse proxy server(3001) via nginx/caddy(443)

2. Enable IP forwarding on server

```
  sudo echo 1 > /proc/sys/net/ipv4/ip_forward
  sudo sysctl -p
  sudo iptables -t nat -A POSTROUTING -s 172.16.0.0/24 -o ens3 -j MASQUERADE
  sudo apt-get install iptables-persistent
  sudo iptables-save > /etc/iptables/rules.v4
```

## Mobile client

### [Android](https://github.com/net-byte/vTunnel)

