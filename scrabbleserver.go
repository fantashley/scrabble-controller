package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type scrabbleServer struct {
	activeGames map[uuid.UUID]ScrabbleGame
}

var (
	serverMu sync.Mutex
	server   scrabbleServer
)

func startScrabbleServer(bindAddr string) error {
	server.activeGames = make(map[uuid.UUID]ScrabbleGame)

	r := mux.NewRouter()
	r.HandleFunc("/games/create", createGameHandler)

	return http.ListenAndServe(bindAddr, r)
}

func createGameHandler(w http.ResponseWriter, r *http.Request) {
	newGame := createScrabbleGame()

	serverMu.Lock()
	server.activeGames[newGame.ID] = newGame
	serverMu.Unlock()

	gameData, err := json.Marshal(newGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(gameData)
}
