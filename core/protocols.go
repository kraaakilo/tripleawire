package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"time"
)

var Protocols = []string{"tcp", "icmpv4"}

type PacketDataStruct struct {
	Source      string `json:"source,omitempty"`
	Destination string `json:"destination,omitempty"`
	Date        string `json:"date,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	Content     string `json:"content,omitempty"`
}

func (receiver PacketDataStruct) String() string {
	return fmt.Sprintf("---\nSource: %s\nDestination: %s\nDate: %s\nProtocol: %s\nContent: %s\n---", receiver.Source, receiver.Destination, receiver.Date, receiver.Protocol, receiver.Content)
}

func handleTCP(packet gopacket.Packet) (PacketDataStruct, error) {
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, ok := tcpLayer.(*layers.TCP)
		if ok {
			if netLayer := packet.NetworkLayer(); netLayer != nil {
				source, destination := netLayer.NetworkFlow().Endpoints()
				return PacketDataStruct{
					Source:      source.String(),
					Destination: destination.String(),
					Date:        time.Now().Format("2006-01-02 15:04:05"),
					Protocol:    tcp.LayerType().String(),
					Content:     packet.String(),
				}, nil
			}
		}
	}
	return PacketDataStruct{}, fmt.Errorf("error when processing packet")
}

func handleICMPV4(packet gopacket.Packet) (PacketDataStruct, error) {
	if layer := packet.Layer(layers.LayerTypeICMPv4); layer != nil {
		_, ok := layer.(*layers.ICMPv4)
		if ok {
			if netLayer := packet.NetworkLayer(); netLayer != nil {
				source, destination := netLayer.NetworkFlow().Endpoints()
				return PacketDataStruct{
					Source:      source.String(),
					Destination: destination.String(),
					Date:        time.Now().Format("2006-01-02 15:04:05"),
					Protocol:    "ICMPv4",
					Content:     packet.String(),
				}, nil
			}
		}
	}
	return PacketDataStruct{}, fmt.Errorf("error when processing packet")
}
