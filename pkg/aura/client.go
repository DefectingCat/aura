package aura

import (
	"errors"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
)

// Command enums
const (
	Nickname = "/nickname"
)

// Ctrl c signal in uinx system telnet client
var CTRL_C = []byte{255, 244, 255, 253, 6}

type MessageType int

// Client message enum
const (
	Connected MessageType = iota
	Message
	Disconnected
)

// The client
type Client struct {
	// Client address
	addr net.Addr
	// Client connection
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
			msg := string(message)
			if isCommand(&msg) {
				if err := parseCommand(&msg); err != nil {
					log.Println(err)
				}
				continue
			}
			log.Printf("[%s]: %s", addr, msg)
			clientCh <- ClientMessage{
				msgType: Message,
				client:  client,
				message: message,
			}
		}
	}
}

func isCommand(msg *string) bool {
	return strings.HasPrefix(*msg, "/")
}

func parseCommand(msg *string) error {
	commands := strings.Fields(*msg)
	if len(commands) <= 1 {
		return errors.New("comand: empty command " + commands[0])
	}
	command := commands[0]
	arg := commands[1]
	log.Println("get command ", command, arg)
	return nil
}
