package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Tile struct {
	Letter byte `json:"letter"`
	Count  int  `json:"-"`
	Value  int  `json:"value"`
}

var tiles = map[byte]Tile{
	' ': Tile{Letter: ' ', Count: 2, Value: 0},
	'A': Tile{Letter: 'A', Count: 9, Value: 1},
	'B': Tile{Letter: 'B', Count: 2, Value: 3},
	'C': Tile{Letter: 'C', Count: 2, Value: 3},
	'D': Tile{Letter: 'D', Count: 4, Value: 2},
	'E': Tile{Letter: 'E', Count: 12, Value: 1},
	'F': Tile{Letter: 'F', Count: 2, Value: 4},
	'G': Tile{Letter: 'G', Count: 3, Value: 2},
	'H': Tile{Letter: 'H', Count: 2, Value: 4},
	'I': Tile{Letter: 'I', Count: 9, Value: 1},
	'J': Tile{Letter: 'J', Count: 1, Value: 8},
	'K': Tile{Letter: 'K', Count: 1, Value: 5},
	'L': Tile{Letter: 'L', Count: 4, Value: 1},
	'M': Tile{Letter: 'M', Count: 2, Value: 3},
	'N': Tile{Letter: 'N', Count: 6, Value: 1},
	'O': Tile{Letter: 'O', Count: 8, Value: 1},
	'P': Tile{Letter: 'P', Count: 2, Value: 3},
	'Q': Tile{Letter: 'Q', Count: 1, Value: 10},
	'R': Tile{Letter: 'R', Count: 6, Value: 1},
	'S': Tile{Letter: 'S', Count: 4, Value: 1},
	'T': Tile{Letter: 'T', Count: 6, Value: 1},
	'U': Tile{Letter: 'U', Count: 4, Value: 1},
	'V': Tile{Letter: 'V', Count: 2, Value: 4},
	'W': Tile{Letter: 'W', Count: 2, Value: 4},
	'X': Tile{Letter: 'X', Count: 1, Value: 8},
	'Y': Tile{Letter: 'Y', Count: 2, Value: 4},
	'Z': Tile{Letter: 'Z', Count: 1, Value: 10},
}

type Player struct {
	ID     uuid.UUID              `json:"-"`
	Name   string                 `json:"name"`
	Number int                    `json:"number"`
	Tiles  []byte                 `json:"-"`
	Score  int                    `json:"score"`
	State  chan GameStateResponse `json:"-"`
}

type TileBag []byte

var initializedTileBag = initializeTileBag()

type ScrabbleGame struct {
	sync.Mutex
	ID         uuid.UUID
	Active     bool
	Action     chan GamePlayRequest
	TurnCount  int
	Board      ScrabbleBoard
	TileBag    TileBag
	Players    map[uuid.UUID]*Player
	PlayerList []*Player
}

func createScrabbleGame() ScrabbleGame {

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

	return game
}

func dealTiles(p *Player, tb *TileBag, tileCount int) {
	var tilesDealt []byte
	tilesDealt, *tb = (*tb)[:tileCount], (*tb)[tileCount:]
	p.Tiles = append(p.Tiles, tilesDealt...)
}

func initializeTileBag() TileBag {
	bag := TileBag{}
	for t := range tiles {
		for i := 0; i < tiles[t].Count; i++ {
			bag = append(bag, t)
		}
	}
	return bag
}

func (tb TileBag) shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	r.Shuffle(len(tb), func(i, j int) {
		tb[i], tb[j] = tb[j], tb[i]
	})
}

func (sg *ScrabbleGame) stateController() {

	// Deal tiles to players
	for p := range sg.Players {
		dealTiles(sg.Players[p], &sg.TileBag, 7)
	}

	sg.PlayerList = sg.playerList()

	numPlayers := len(sg.Players)

	for request := range sg.Action {
		switch request.Play {
		case false: // Return the game state
			sg.Players[request.PlayerID].State <- GameStateResponse{
				GameID:      sg.ID,
				PlayerID:    request.PlayerID,
				Players:     sg.PlayerList,
				Board:       sg.Board,
				PlayerTurn:  sg.TurnCount % numPlayers,
				PlayerTiles: sg.Players[request.PlayerID].Tiles,
			}
		}
	}
}

func (sg *ScrabbleGame) playerList() []*Player {
	p := make([]*Player, len(sg.Players))
	for player := range sg.Players {
		p[sg.Players[player].Number] = sg.Players[player]
	}
	return p
}

func printTiles(tiles []byte) {
	fmt.Println(string(tiles))
}
