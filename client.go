package main

import (
	"log"
	"net"
)

func HandleClient(conn net.Conn) {
	log.Println("accpeted client from ", conn.RemoteAddr())
}
