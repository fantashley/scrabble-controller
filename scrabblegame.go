package main

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
	tiles []tile
	score int
}

type scrabbleGame struct {
	board   scrabbleBoard
	tileBag []tile
	players []player
}

func createScrabbleGame() scrabbleGame {
	game := scrabbleGame{}
	game.board.initializeSquares()

	return game
}
