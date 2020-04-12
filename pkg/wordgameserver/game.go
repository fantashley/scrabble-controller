package wordgameserver

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Tile represents a Scrabble tile that would be played on a board
type Tile struct {
	Letter byte `json:"letter"` // the character written on the tile
	Count  int  `json:"-"`      // the number of tiles with the character
	Value  int  `json:"value"`  // the point value of playing the tile
}

var tiles = map[byte]Tile{
	' ': {Letter: ' ', Count: 2, Value: 0},
	'A': {Letter: 'A', Count: 9, Value: 1},
	'B': {Letter: 'B', Count: 2, Value: 3},
	'C': {Letter: 'C', Count: 2, Value: 3},
	'D': {Letter: 'D', Count: 4, Value: 2},
	'E': {Letter: 'E', Count: 12, Value: 1},
	'F': {Letter: 'F', Count: 2, Value: 4},
	'G': {Letter: 'G', Count: 3, Value: 2},
	'H': {Letter: 'H', Count: 2, Value: 4},
	'I': {Letter: 'I', Count: 9, Value: 1},
	'J': {Letter: 'J', Count: 1, Value: 8},
	'K': {Letter: 'K', Count: 1, Value: 5},
	'L': {Letter: 'L', Count: 4, Value: 1},
	'M': {Letter: 'M', Count: 2, Value: 3},
	'N': {Letter: 'N', Count: 6, Value: 1},
	'O': {Letter: 'O', Count: 8, Value: 1},
	'P': {Letter: 'P', Count: 2, Value: 3},
	'Q': {Letter: 'Q', Count: 1, Value: 10},
	'R': {Letter: 'R', Count: 6, Value: 1},
	'S': {Letter: 'S', Count: 4, Value: 1},
	'T': {Letter: 'T', Count: 6, Value: 1},
	'U': {Letter: 'U', Count: 4, Value: 1},
	'V': {Letter: 'V', Count: 2, Value: 4},
	'W': {Letter: 'W', Count: 2, Value: 4},
	'X': {Letter: 'X', Count: 1, Value: 8},
	'Y': {Letter: 'Y', Count: 2, Value: 4},
	'Z': {Letter: 'Z', Count: 1, Value: 10},
}

// Player represents an instance of a player and stores their current state
type Player struct {
	ID     uuid.UUID              `json:"-"`      // unique identifier
	Name   string                 `json:"name"`   // player's chosen display name
	Number int                    `json:"number"` // number that dictates their turn
	Tiles  []byte                 `json:"-"`      // tiles currenty in possession
	Score  int                    `json:"score"`  // current score in the game
	State  chan GameStateResponse `json:"-"`      // channel on which to send state responses
	Play   chan GameStateResponse `json:"-"`      // channel on which to send play responses
}

// TileBag represents the bag of undistributed tiles in a game
type TileBag []byte

var initializedTileBag = initializeTileBag()

// ScrabbleGame represents the state of an active game instance
type ScrabbleGame struct {
	sync.Mutex
	ID        uuid.UUID             // unique identifier
	Active    bool                  // true if the game has started
	Action    chan GamePlayRequest  // channel for receiving player's turns
	TurnCount int                   // counter that increments for each turn played
	Board     ScrabbleBoard         // board representation with current tiles
	TileBag   TileBag               // bag of tiles not yet distributed
	Players   map[uuid.UUID]*Player // players indexed by UUID
}

// createScrabbleGame initializes a game instance
func createScrabbleGame() *ScrabbleGame {

	game := ScrabbleGame{}

	game.ID = uuid.New()

	game.Action = make(chan GamePlayRequest)

	// Initialize squares on board
	game.Board = initializedBoard

	// Populate tile bag
	game.TileBag = make(TileBag, len(initializedTileBag))
	copy(game.TileBag, initializedTileBag)

	// Shuffle tile bag
	game.TileBag.shuffle()

	game.Players = make(map[uuid.UUID]*Player)

	return &game
}

// dealTiles disperses tiles from the tile bag to players so they always have 7
// tiles in their hand
func dealTiles(p *Player, tb *TileBag, tileCount int) {
	var tilesDealt []byte
	tilesDealt, *tb = (*tb)[:tileCount], (*tb)[tileCount:]
	p.Tiles = append(p.Tiles, tilesDealt...)
}

// initializeTileBag fills the tile bag with tiles before the game begins
func initializeTileBag() TileBag {
	bag := TileBag{}
	for t := range tiles {
		for i := 0; i < tiles[t].Count; i++ {
			bag = append(bag, t)
		}
	}
	return bag
}

// shuffle make sure the tiles are in random order in the tile bag
func (tb TileBag) shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	r.Shuffle(len(tb), func(i, j int) {
		tb[i], tb[j] = tb[j], tb[i]
	})
}

func (sg *ScrabbleGame) start() error {

	if sg.Active {
		return errors.New("Game has already started")
	} else if len(sg.Players) < 2 {
		return errors.New("At least two players needed to start game")
	}

	sg.Active = true

	go sg.stateController()

	return nil
}

// stateController is the main goroutine for the game that handles state
// requests and play requests
func (sg *ScrabbleGame) stateController() {

	// Deal tiles to players
	for p := range sg.Players {
		dealTiles(sg.Players[p], &sg.TileBag, 7)
	}

	// Get ordered list of players to send to clients
	playerList := sg.playerList()

	numPlayers := len(sg.Players)

	// Loop on requests in queue
	for request := range sg.Action {
		switch request.Play {
		case false: // Return the game state
			sg.Players[request.PlayerID].State <- GameStateResponse{
				GameID:      sg.ID,
				PlayerID:    request.PlayerID,
				Players:     playerList,
				Board:       sg.Board,
				PlayerTurn:  sg.TurnCount % numPlayers,
				PlayerTiles: sg.Players[request.PlayerID].Tiles,
			}
		}
	}
}

func (sg *ScrabbleGame) request(r GamePlayRequest) (GameStateResponse, error) {

	// Send request to game controller
	sg.Action <- r

	switch r.Play {
	case false:
		return <-sg.Players[r.PlayerID].State, nil
	}

	return GameStateResponse{}, errors.New("Play functionality not yet implemented")
}

// playerList generates an ordered list of players for consistency across all
// clients
func (sg *ScrabbleGame) playerList() []*Player {
	p := make([]*Player, len(sg.Players))
	for player := range sg.Players {
		p[sg.Players[player].Number] = sg.Players[player]
	}
	return p
}

// addPlayer checks that a new player can be added to the game, and adds the
// player if so
func (sg *ScrabbleGame) addPlayer(name string) (uuid.UUID, error) {

	// Create player to be added to game
	p := Player{
		ID:    uuid.New(),
		Name:  name,
		Tiles: make([]byte, 0),
		State: make(chan GameStateResponse),
	}

	playerCount := len(sg.Players)

	// Check that game is valid to join
	if playerCount == 4 {
		return p.ID, errors.New("Maximum players reached for game")
	} else if sg.Active {
		return p.ID, errors.New("Game has already started")
	}

	// Assign player their number based on when they joined
	p.Number = playerCount
	// Add player to game
	sg.Players[p.ID] = &p

	return p.ID, nil
}

// printTiles prints the tiles in a player's hand
func printTiles(tiles []byte) {
	fmt.Println(string(tiles))
}
