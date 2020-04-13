package wordgameserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func TestCreateGameHandler(t *testing.T) {
	var j GeneralGameRequest

	req, err := http.NewRequest("GET", "/game/create", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(createGameHandler)

	h.ServeHTTP(rr, req)

	if c := rr.Code; c != http.StatusCreated {
		t.Fatalf("Returned status code %v, expected %v",
			c, http.StatusCreated)
	}

	err = json.NewDecoder(rr.Body).Decode(&j)
	if err != nil {
		t.Fatalf("Did not return game_id in correct format. Returned: %v",
			rr.Body)
	} else if j.GameID == (uuid.UUID{}) {
		t.Error("Returned empty game_id")
	}

	_, err = getGame(j.GameID, rr)
	if err != nil {
		t.Fatalf("No existing games with ID %v. Map contents: %v",
			j.GameID, server.activeGames)
	}
}

func TestJoinGameHandler(t *testing.T) {

	newGame := createScrabbleGame()
	maxPlayers := 4

	serverMu.Lock()
	server.activeGames[newGame.ID] = newGame
	serverMu.Unlock()

	playerNames := []string{
		"ashley1",
		"ashley2",
		"ashley3",
		"ashley4",
		"ashley5",
		"ashley6",
	}

	joinCh := make(chan *httptest.ResponseRecorder)
	errCh := make(chan error)

	// Test concurrently joining players
	for i := 0; i < maxPlayers; i++ {
		go joinPlayer(newGame.ID, playerNames[i], joinCh, errCh)
	}

	rrCount := 0

joinLoop:
	for {
		select {
		case rr := <-joinCh:
			rrCount++
			if c := rr.Code; c != http.StatusOK {
				t.Fatalf("Returned status code %v, expected %v. Error: %v",
					c, http.StatusOK, rr.Body)
			}

			j := GeneralGameRequest{}
			err := json.NewDecoder(rr.Body).Decode(&j)
			if err != nil {
				t.Fatalf("Did not return game ID and player ID as expected. Returned %v",
					rr.Body)
			}

			if j.PlayerID == nil {
				t.Fatal("Player ID in response is nil")
			}
			if rrCount == 4 {
				break joinLoop
			}
		case err := <-errCh:
			t.Errorf("Player failed to join with message: %v", err)
		}
	}

	for i := 4; i < 6; i++ {
		go joinPlayer(newGame.ID, playerNames[i], joinCh, errCh)
	}

	rrCount = 0

failJoinLoop:
	for {
		select {
		case err := <-errCh:
			t.Fatal(err)
		case rr := <-joinCh:
			rrCount++
			if rr.Code != http.StatusBadRequest {
				t.Error("Should have failed to add player")
			}
			if rrCount == 2 {
				break failJoinLoop
			}
		}
	}

	if len(newGame.Players) != maxPlayers {
		t.Errorf("Game has %v players, expected %v", len(newGame.Players), maxPlayers)
	}
}

func joinPlayer(gameID uuid.UUID, playerName string, joinCh chan *httptest.ResponseRecorder, errCh chan error) {
	j := GeneralGameRequest{
		GameID:     gameID,
		PlayerName: &playerName,
	}

	payload, err := json.Marshal(j)
	if err != nil {
		errCh <- errors.New("Failed to marshal JSON request object")
	}
	req, err := http.NewRequest("POST", "/game/join", bytes.NewBuffer(payload))
	if err != nil {
		errCh <- err
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(joinGameHandler)

	h.ServeHTTP(rr, req)

	joinCh <- rr
}

func TestStartGameHandler(t *testing.T) {
	newGame := createScrabbleGame()

	serverMu.Lock()
	server.activeGames[newGame.ID] = newGame
	serverMu.Unlock()

	playerNames := []string{
		"ashley1",
		"ashley2",
		"ashley3",
		"ashley4",
	}

	newGame.addPlayer(playerNames[0])

	j := GeneralGameRequest{
		GameID: newGame.ID,
	}

	payload, err := json.Marshal(j)
	if err != nil {
		t.Fatal("Failed to marshal JSON request object")
	}

	rr, err := startGame(payload)
	if err != nil {
		t.Fatal(err)
	}

	// Should fail with only one player
	if c := rr.Code; c == http.StatusOK {
		t.Fatal("Game should not start with one player")
	}

	// Add remaining players so game should start
	for i := 1; i < len(playerNames); i++ {
		newGame.addPlayer(playerNames[i])
	}

	rr, err = startGame(payload)
	if err != nil {
		t.Fatal(err)
	}

	// Game should start successfully
	if c := rr.Code; c != http.StatusOK {
		t.Fatalf("Returned status code %v, expected %v",
			c, http.StatusOK)
	}

	rr, err = startGame(payload)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure it is unsuccessful if already started
	if c := rr.Code; c == http.StatusOK {
		t.Fatal("Second start attempt should have failed")
	}
}

func startGame(payload []byte) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest("GET", "/game/start", bytes.NewBuffer(payload))
	if err != nil {
		return nil, errors.New("Failed to generate HTTP request for game start")
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(startGameHandler)

	h.ServeHTTP(rr, req)

	return rr, nil
}

func TestGameStateHandler(t *testing.T) {

	newGame := createScrabbleGame()
	var playerID uuid.UUID
	var err error

	serverMu.Lock()
	server.activeGames[newGame.ID] = newGame
	serverMu.Unlock()

	playerNames := []string{
		"ashley1",
		"ashley2",
		"ashley3",
		"ashley4",
	}

	for i := 0; i < len(playerNames); i++ {
		playerID, err = newGame.addPlayer(playerNames[i])
		if err != nil {
			t.Fatal("Failed to add valid player to game")
		}
	}

	newGame.start()

	j := GeneralGameRequest{
		GameID:   newGame.ID,
		PlayerID: &playerID,
	}

	payload, err := json.Marshal(j)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/game/state", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	h := http.HandlerFunc(gameStateHandler)

	h.ServeHTTP(rr, req)

	if c := rr.Code; c != http.StatusOK {
		t.Fatalf("Returned status code %v, expected %v", c, http.StatusOK)
	}

	var s GameStateResponse

	err = json.NewDecoder(rr.Body).Decode(&s)
	if err != nil {
		t.Fatal("Response was not in correct format")
	}

	if s.Players[3].Name != playerNames[3] {
		t.Fatalf("Expected player name %v, got %v", playerNames[3], s.Players[4].Name)
	} else if s.PlayerTurn != 0 {
		t.Fatalf("Game not in correct turn state")
	} else if len(s.PlayerTiles) != 7 {
		t.Fatal("Incorrect number of tiles for player")
	}
}
