package main

import (
	"log/slog"
	"net"

	freeD "github.com/jwetzell/free-d-go"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:3333")
	if err != nil {
		slog.Error("error making UDP address", "err", err)
		return
	}

	client, err := net.ListenUDP("udp", addr)
	if err != nil {
		slog.Error("error listening to UDP", "err", err)
		return
	}
	defer client.Close()

	for {
		buffer := make([]byte, 2048)

		length, _, err := client.ReadFromUDP(buffer)
		if err != nil {
			slog.Error("error reading from UDP", "err", err)
		} else if length > 0 {
			message, err := freeD.Decode(buffer[0:length])
			if err != nil {
				slog.Error("error decoding", "err", err)
			}
			slog.Info("decoded", "message", message)
		}
	}
}
