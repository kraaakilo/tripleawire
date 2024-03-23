package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type ClientPacketMessage struct {
	Message     PacketDataStruct `json:"message"`
	MessageType string           `json:"type"`
}

type ClientAlertMessage struct {
	Message     string `json:"message"`
	MessageType string `json:"type"`
}

var upgrader = websocket.Upgrader{}

func WebServer(interfaceName string) {
	fmt.Println("Starting web server on http://localhost:8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		content := []byte(`
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <link rel="stylesheet" href="https://necolas.github.io/normalize.css/8.0.1/normalize.css">
    <title>tripleawire</title>
</head>
<body>
<div>
    <h3 style="text-align:center">tripleawire monitoring</h3>
	<p style="text-align:center">The websocket server is running on ws://localhost:8080/ws</p>
</body>
</html>
`)
		r.Header.Set("Content-Type", "text/html")
		if _, err := w.Write(content); err != nil {
			return
		}
	})

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		webSocketHandler(writer, request, interfaceName)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func webSocketHandler(w http.ResponseWriter, r *http.Request, interfaceName string) {

	blacklist, err := loadBlacklist("blacklist.yaml")
	if err != nil {
		log.Fatal(err)
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	handle, err := pcap.OpenLive(interfaceName, 1600, true, pcap.BlockForever)
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		for packet := range packetSource.Packets() {
			if err := packet.ErrorLayer(); err != nil {
				fmt.Println("Error trying to decode the content of the packet.", err)
			}

			// Handle the packet

			if content, err := handleTCP(packet); err == nil {
				if marshal, err := convertPacketToJsonBytes(content, err); err == nil {
					err = c.WriteMessage(websocket.TextMessage, marshal)
				}
				// Check if the IP is blacklisted
				if ok, msg := CheckIPisBlacklisted(blacklist, content); ok {
					if marshal, err := json.Marshal(ClientAlertMessage{Message: msg, MessageType: "alert"}); err == nil {
						err = c.WriteMessage(websocket.TextMessage, marshal)
					}
				}
			}

			if content, err := handleICMPV4(packet); err == nil {
				if marshal, err := convertPacketToJsonBytes(content, err); err == nil {
					err = c.WriteMessage(websocket.TextMessage, marshal)
				}
			}

		}
	}

	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}(c)
}

func convertPacketToJsonBytes(d PacketDataStruct, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	marshal, err := json.Marshal(ClientPacketMessage{
		Message:     d,
		MessageType: "packet",
	})
	if err != nil {
		return nil, err
	}
	return marshal, nil
}
