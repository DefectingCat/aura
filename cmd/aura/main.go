package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"

	"git.rua.plus/xfy/pkg/aura"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("not found .env file")
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "4000"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listening on %s failed: %s", listener.Addr(), err)
	}
	log.Println("listening on ", listener.Addr())

	clientCh := make(chan aura.ClientMessage)
	go aura.HandleServer(clientCh)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept connection failed", err)
			continue
		}
		go aura.HandleClient(conn, clientCh)
	}
}
