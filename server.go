package main

import (
	"log"
	"net"
)

func HandleServer(clientCh <-chan ClientMessage) {
	clients := make(map[string]Client)

	for clientMsg := range clientCh {
		clientKey := clientMsg.client.addr.String()
		switch clientMsg.msgType {
		case Connected:
			clients[clientKey] = clientMsg.client
		case Message:
			for k, v := range clients {
				if k == clientKey {
					continue
				}
				go func(conn net.Conn, msg []byte) {
					_, err := conn.Write(msg)
					if err != nil {
						log.Println("write message to client failed ", err)
					}
				}(v.conn, clientMsg.message)
			}
		case Disconnected:
			delete(clients, clientKey)
		default:
		}
	}
}
