package main

import (
	"io"
	"log"
	"net"
	"reflect"
)

type MessageType int

var CTRL_C = []byte{255, 244, 255, 253, 6}

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

// HandleClient handle user tcp connection.
// Each connection will be sent to server goroutine.
// Will receive data from client and send to server goroutine.
func HandleClient(conn net.Conn, clientCh chan<- ClientMessage) {
	addr := conn.RemoteAddr()
	log.Println("accepted client from ", addr)
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Printf("close client connection failed %s", err)
		}
	}(conn)

	// send a welcome message
	_, err := conn.Write([]byte("Genshin impact, Launch!\n"))
	if err != nil {
		log.Printf("write welcome message to client failed %s", err)
		return
	}

	client := Client{
		addr,
		conn,
	}
	// register client to server
	clientCh <- ClientMessage{
		msgType: Connected,
		client:  client,
		message: []byte{},
	}

	// receive data from client
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
			isCtrlc := reflect.DeepEqual(CTRL_C, message)
			if isCtrlc {
				continue
			}
			log.Printf("[%s]: %s", addr, string(message))
			clientCh <- ClientMessage{
				msgType: Message,
				client:  client,
				message: message,
			}
		}
	}
}
