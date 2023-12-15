package main

import (
	"io"
	"log"
	"net"
)

type MessageType int

const (
	Connected MessageType = iota
	Message
	Disconnected
)

type Client struct {
	addr net.Addr
	conn net.Conn
}
type ClientMessage struct {
	msgType MessageType
	client  Client
	message []byte
}

func HandleClient(conn net.Conn, clientCh chan<- ClientMessage) {
	addr := conn.RemoteAddr()
	log.Println("accpeted client from ", addr)
	defer conn.Close()

	// send a welcome message
	conn.Write([]byte("Genshin impact, Launch!\n"))

	client := Client{
		addr,
		conn,
	}
	clientCh <- ClientMessage{
		msgType: Connected,
		client:  client,
		message: []byte{},
	}

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("client %s disconnected", addr)
				clientCh <- ClientMessage{
					msgType: Disconnected,
					client:  client,
					message: []byte{},
				}
				return
			} else {
				log.Printf("receive data failed %s from %s\n", err, addr)
				return
			}
		}
		if n < 3 {
			continue
		} else {
			message := buffer[:n]
			log.Printf("[%s]: %s", addr, string(message))
			clientCh <- ClientMessage{
				msgType: Message,
				client:  client,
				message: message,
			}
		}
	}
}
