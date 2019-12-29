package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type tile struct {
	letter byte
	count  int
	value  int
}

var tiles = map[byte]tile{
	' ': tile{letter: ' ', count: 2, value: 0},
	'A': tile{letter: 'A', count: 9, value: 1},
	'B': tile{letter: 'B', count: 2, value: 3},
	'C': tile{letter: 'C', count: 2, value: 3},
	'D': tile{letter: 'D', count: 4, value: 2},
	'E': tile{letter: 'E', count: 12, value: 1},
	'F': tile{letter: 'F', count: 2, value: 4},
	'G': tile{letter: 'G', count: 3, value: 2},
	'H': tile{letter: 'H', count: 2, value: 4},
	'I': tile{letter: 'I', count: 9, value: 1},
	'J': tile{letter: 'J', count: 1, value: 8},
	'K': tile{letter: 'K', count: 1, value: 5},
	'L': tile{letter: 'L', count: 4, value: 1},
	'M': tile{letter: 'M', count: 2, value: 3},
	'N': tile{letter: 'N', count: 6, value: 1},
	'O': tile{letter: 'O', count: 8, value: 1},
	'P': tile{letter: 'P', count: 2, value: 3},
	'Q': tile{letter: 'Q', count: 1, value: 10},
	'R': tile{letter: 'R', count: 6, value: 1},
	'S': tile{letter: 'S', count: 4, value: 1},
	'T': tile{letter: 'T', count: 6, value: 1},
	'U': tile{letter: 'U', count: 4, value: 1},
	'V': tile{letter: 'V', count: 2, value: 4},
	'W': tile{letter: 'W', count: 2, value: 4},
	'X': tile{letter: 'X', count: 1, value: 8},
	'Y': tile{letter: 'Y', count: 2, value: 4},
	'Z': tile{letter: 'Z', count: 1, value: 10},
}

type player struct {
	name  string
	tiles []byte
	score int
}

type tileBag []byte

type scrabbleGame struct {
	board scrabbleBoard
	tileBag
	players []player
}

func createScrabbleGame(playerNames []string) (scrabbleGame, error) {

	game := scrabbleGame{}

	if len(playerNames) < 2 || len(playerNames) > 4 {
		return game, errors.New("Player count must be between 2 and 4")
	}

	// Initialize squares on board
	game.board.initialize()

	// Populate tile bag
	for t := range tiles {
		for i := 0; i < tiles[t].count; i++ {
			game.tileBag = append(game.tileBag, t)
		}
	}

	// Shuffle tile bag
	game.tileBag.shuffle()

	// Create players
	for _, name := range playerNames {
		game.newPlayer(name)
	}

	return game, nil
}

func (sg *scrabbleGame) newPlayer(playerName string) {
	var p player
	p.name = playerName
	_ = dealTiles(&p, &sg.tileBag, 7)
	sg.players = append(sg.players, p)
}

func dealTiles(p *player, tb *tileBag, tileCount int) []byte {
	var tilesDealt []byte
	tilesDealt, *tb = (*tb)[:tileCount], (*tb)[tileCount:]
	p.tiles = append(p.tiles, tilesDealt...)
	return tilesDealt
}

func (tb tileBag) shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	r.Shuffle(len(tb), func(i, j int) {
		tb[i], tb[j] = tb[j], tb[i]
	})
}

func printTiles(tiles []byte) {
	fmt.Println(string(tiles))
}
