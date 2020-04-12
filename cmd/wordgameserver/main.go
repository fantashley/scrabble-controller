package main

import (
	"log"

	"github.com/fantashley/wordgame-controller/pkg/wordgameserver"
)

func main() {

	log.Fatal(wordgameserver.StartWordGameServer(":8080"))
}
