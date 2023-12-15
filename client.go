package main

import (
	"log"
	"net"
)

type ConnectType int

const (
	Connected ConnectType = iota
	Disconnected
)

type Client struct {
	addr net.Addr
	conn net.Conn
}
type ClientMessage struct {
	msgType ConnectType
	client  Client
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
	}

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("receive data failed %s from %s\n", err, addr)
			return
		}
		if n < 3 {
			continue
		} else {
			log.Printf("[%s]: %s", addr, string(buffer[:n]))
			clear(buffer)
		}
	}
}
