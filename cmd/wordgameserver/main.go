package main

import (
	"flag"
	"log"

	"github.com/fantashley/wordgame-controller/pkg/wordgameserver"
)

func main() {
	bindAddr := flag.String("bind-addr", ":8080", "Bind address for server (ip:port)")
	wordFile := flag.String("word-file", "/data/words.txt", "File containing word dictionary")
	server, err := wordgameserver.GetWordGameServer(*wordFile)
	if err != nil {
		log.Fatalf("Failed to start server. %v", err.Error())
	}
	log.Fatal(server.Start(*bindAddr))
}
