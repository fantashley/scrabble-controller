package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type scrabbleServer struct {
	activeGames   map[uuid.UUID]*ScrabbleGame
	activePlayers map[uuid.UUID]*Player
}

type JoinGameRequest struct {
	GameID   uuid.UUID `json:"game_id"`
	PlayerID uuid.UUID `json:"player_id,omitempty"`
	Player   `json:"player"`
}

var (
	serverMu sync.Mutex
	server   scrabbleServer
)

func startScrabbleServer(bindAddr string) error {
	server.activeGames = make(map[uuid.UUID]*ScrabbleGame)
	server.activePlayers = make(map[uuid.UUID]*Player)

	r := mux.NewRouter()
	r.HandleFunc("/games/create", createGameHandler)
	r.HandleFunc("/games/join", joinGameHandler)

	return http.ListenAndServe(bindAddr, r)
}

func createGameHandler(w http.ResponseWriter, r *http.Request) {
	newGame := createScrabbleGame()

	serverMu.Lock()
	server.activeGames[newGame.ID] = &newGame
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

func joinGameHandler(w http.ResponseWriter, r *http.Request) {
	var j JoinGameRequest
	var g *ScrabbleGame
	var ok bool

	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	j.Player.ID = uuid.New()
	j.PlayerID = j.Player.ID

	serverMu.Lock()
	if g, ok = server.activeGames[j.GameID]; !ok {
		serverMu.Unlock()
		http.Error(w, "No active game with that ID", http.StatusBadRequest)
		return
	}
	server.activePlayers[j.Player.ID] = &(j.Player)
	serverMu.Unlock()

	g.PlayerMu.Lock()
	g.Players = append(g.Players, &(j.Player))
	g.PlayerMu.Unlock()

	resp, err := json.Marshal(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
