package main

import (
	"log"
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
				_, err := v.conn.Write(clientMsg.message)
				if err != nil {
					log.Println("write message to client failed ", err)
				}
			}
		case Disconnected:
			delete(clients, clientKey)
		default:
		}
	}
}
