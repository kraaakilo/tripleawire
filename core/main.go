package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func main() {
	app := &cli.App{
		Name: "triplewire",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases: []string{"i"},
				Name:    "interface",
				Usage:   "interface to listen",
			},
			&cli.StringFlag{
				Aliases: []string{"m"},
				Name:    "mode",
				Usage:   "mode to run the tool in. Ex: cli, web",
			},
		},
		Usage:       "triplewire is a tool to capture packets and display them in a web interface.",
		Description: "The tool offers a web socket server that listens to packets on a specified interface and emits them. It also offers a cli to real-time capture packets.",
		Action: func(cCtx *cli.Context) error {
			if len(cCtx.String("interface")) == 0 {
				fmt.Println("Please specify an interface to listen to. Ex: triplewire --interface wlan0")
				return nil
			}
			mode := cCtx.String("mode")
			if mode != "cli" && mode != "web" {
				fmt.Println("Please specify a mode to run the tool in. Ex: triplewire --mode web")
				return nil
			}
			if mode == "cli" {
				startLiveCapture(cCtx.String("interface"))
			} else if mode == "web" {
				WebServer(cCtx.String("interface"))
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	return
}

func startLiveCapture(i string) {
	fmt.Println("Starting live capture on interface", i, " at : ", time.Now().Format("2006-01-02 15:04:05"))
	blacklist, err := loadBlacklist("blacklist.yaml")
	if err != nil {
		log.Fatal(err)
	}

	config := Config{
		protocols: Protocols,
	}
	handle, err := pcap.OpenLive(i, 1600, true, pcap.BlockForever)
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for {
			select {
			case packet := <-packetSource.Packets():
				if packet == nil {
					return
				}
				if err := packet.ErrorLayer(); err != nil {
					fmt.Println("Error trying to decode the content of the packet.", err)
				}

				for _, protocol := range config.protocols {
					switch protocol {
					case "tcp":
						content, err := handleTCP(packet)
						if err == nil {
							fmt.Println(content)
							ok, msg := CheckIPisBlacklisted(blacklist, content)
							if ok {
								fmt.Println(msg)
							}
						}
					case "icmpv4":
						content, err := handleICMPV4(packet)
						if err == nil {
							fmt.Println(content)
							ok, msg := CheckIPisBlacklisted(blacklist, content)
							if ok {
								fmt.Println(msg)
							}
						}
					}
				}
			}
		}
	}
}
