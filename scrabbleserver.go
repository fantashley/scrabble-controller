package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type scrabbleServer struct {
	activeGames map[uuid.UUID]*ScrabbleGame
}

type GeneralGameRequest struct {
	GameID     uuid.UUID  `json:"game_id"`
	PlayerID   *uuid.UUID `json:"player_id,omitempty"`
	PlayerName *string    `json:"player_name,omitempty"`
}

type GameStateResponse struct {
	GameID      uuid.UUID     `json:"game_id"`
	PlayerID    uuid.UUID     `json:"-"`
	Players     []*Player     `json:"players"`
	Board       ScrabbleBoard `json:"board"`
	PlayerTurn  int           `json:"turn"`
	PlayerTiles []byte        `json:"tiles"`
}

type GamePlayRequest struct {
	GameID   uuid.UUID        `json:"game_id"`
	PlayerID uuid.UUID        `json:"player_id"`
	StartPos SquareCoordinate `json:"start_pos"`
	EndPos   SquareCoordinate `json:"end_pos"`
	Tiles    []byte           `json:"tiles"`
	Swap     bool             `json:"swap"`
	Play     bool             `json:"-"`
}

var (
	serverMu sync.Mutex
	server   scrabbleServer
)

func startScrabbleServer(bindAddr string) error {
	server.activeGames = make(map[uuid.UUID]*ScrabbleGame)

	r := mux.NewRouter()
	r.HandleFunc("/game/create", createGameHandler)
	r.HandleFunc("/game/join", joinGameHandler)
	r.HandleFunc("/game/start", startGameHandler)
	r.HandleFunc("/game/state", gameStateHandler)

	return http.ListenAndServe(bindAddr, r)
}

func createGameHandler(w http.ResponseWriter, r *http.Request) {
	newGame := createScrabbleGame()

	resp := GeneralGameRequest{
		GameID: newGame.ID,
	}

	serverMu.Lock()
	server.activeGames[newGame.ID] = &newGame
	serverMu.Unlock()

	gameData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(gameData)
}

func joinGameHandler(w http.ResponseWriter, r *http.Request) {
	var j GeneralGameRequest
	var g *ScrabbleGame

	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p := Player{
		ID:    uuid.New(),
		Name:  *j.PlayerName,
		Tiles: make([]byte, 0),
		State: make(chan GameStateResponse),
	}

	j.PlayerID = &p.ID

	g, err = getGame(j.GameID, &w)
	if err != nil {
		return
	}

	g.Lock()
	playerCount := len(g.Players)
	if playerCount == 4 {
		g.Unlock()
		http.Error(w, "Maximum players reached for game", http.StatusBadRequest)
		return
	} else if g.Active {
		g.Unlock()
		http.Error(w, "Game has already started", http.StatusBadRequest)
		return
	}
	p.Number = playerCount
	g.Players[p.ID] = &p
	g.Unlock()

	resp, err := json.Marshal(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	var j GeneralGameRequest
	var g *ScrabbleGame

	// Decode Game ID
	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve game instance
	g, err = getGame(j.GameID, &w)
	if err != nil {
		return
	}

	// Set game to active
	g.Lock()
	defer g.Unlock()
	if g.Active {
		http.Error(w, "Game has already started", http.StatusBadRequest)
		return
	}
	if len(g.Players) < 2 {
		http.Error(w, "At least two players needed to start game", http.StatusBadRequest)
		return
	}
	g.Active = true

	go g.stateController()

	w.WriteHeader(http.StatusOK)
}

func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	var j GeneralGameRequest

	err := json.NewDecoder(r.Body).Decode(&j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g, err := getGame(j.GameID, &w)
	if err != nil {
		return
	}

	g.Action <- GamePlayRequest{
		GameID:   j.GameID,
		PlayerID: *j.PlayerID,
	}

	resp, err := json.Marshal(<-g.Players[*j.PlayerID].State)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func getGame(gameID uuid.UUID, w *(http.ResponseWriter)) (*ScrabbleGame, error) {
	var g *ScrabbleGame
	var ok bool
	serverMu.Lock()
	defer serverMu.Unlock()
	if g, ok = server.activeGames[gameID]; !ok {
		http.Error(*w, "No existing game with that ID", http.StatusBadRequest)
		return nil, errors.New("Game does not exist")
	}
	return g, nil
}
