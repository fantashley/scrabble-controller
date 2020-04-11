package main

import (
	"log"

	"github.com/fantashley/scrabble-controller/pkg/scrabbleserver"
)

func main() {

	log.Fatal(scrabbleserver.StartScrabbleServer(":8080"))
}
