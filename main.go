package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "4000"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listening on %s failed: %s", listener.Addr(), err)
	}
	log.Println("listening on ", listener.Addr())

	clientCh := make(chan ClientMessage)
	go HandleServer(clientCh)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept connection failed", err)
			continue
		}
		go HandleClient(conn, clientCh)
	}
}
