package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "4000"
	}
	log.Println(port)
	/* if err != nil {
	        log.Fatalln("error loading .env file: ", err)
		} */

	/* listener := net.Listen("tcp", address string) */
}
