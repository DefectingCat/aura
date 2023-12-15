package main

import (
	"log"
)

func HandleServer(clientCh <-chan ClientMessage) {
	for clientMsg := range clientCh {
		log.Println(clientMsg.msgType)
	}
}
