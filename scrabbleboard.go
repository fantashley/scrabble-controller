package main

import "fmt"

type SquareCoordinate struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type SquareType struct {
	Name             string             `json:"name"`
	LetterMultiplier int                `json:"letterMultiplier"`
	WordMultiplier   int                `json:"wordMultiplier"`
	Coordinates      []SquareCoordinate `json:"-"`
}

type Square struct {
	SquareType string `json:"type"`
	Tile       `json:"tile,omitempty"`
}

const rowCount int = 15
const columnCount int = 15

type ScrabbleBoard [rowCount][columnCount]Square

var initializedBoard = initializeScrabbleBoard()

var squareTypes = map[string]SquareType{
	"plain": SquareType{
		Name:             "plain",
		LetterMultiplier: 1,
		WordMultiplier:   1,
	},
	"star": SquareType{
		Name:             "star",
		LetterMultiplier: 1,
		WordMultiplier:   1,
		Coordinates: []SquareCoordinate{
			SquareCoordinate{Row: 7, Col: 7},
		},
	},
	"doubleLetter": SquareType{
		Name:             "doubleLetter",
		LetterMultiplier: 2,
		WordMultiplier:   1,
		Coordinates: []SquareCoordinate{
			SquareCoordinate{Row: 0, Col: 3},
			SquareCoordinate{Row: 2, Col: 6},
			SquareCoordinate{Row: 3, Col: 0},
			SquareCoordinate{Row: 3, Col: 7},
			SquareCoordinate{Row: 6, Col: 2},
			SquareCoordinate{Row: 6, Col: 6},
			SquareCoordinate{Row: 7, Col: 3},
		},
	},
	"doubleWord": SquareType{
		Name:             "doubleWord",
		LetterMultiplier: 1,
		WordMultiplier:   2,
		Coordinates: []SquareCoordinate{
			SquareCoordinate{Row: 1, Col: 1},
			SquareCoordinate{Row: 2, Col: 2},
			SquareCoordinate{Row: 3, Col: 3},
			SquareCoordinate{Row: 4, Col: 4},
		},
	},
	"tripleLetter": SquareType{
		Name:             "tripleLetter",
		LetterMultiplier: 3,
		WordMultiplier:   1,
		Coordinates: []SquareCoordinate{
			SquareCoordinate{Row: 1, Col: 5},
			SquareCoordinate{Row: 5, Col: 1},
			SquareCoordinate{Row: 5, Col: 5},
		},
	},
	"tripleWord": SquareType{
		Name:             "tripleWord",
		LetterMultiplier: 1,
		WordMultiplier:   3,
		Coordinates: []SquareCoordinate{
			SquareCoordinate{Row: 0, Col: 0},
			SquareCoordinate{Row: 0, Col: 7},
			SquareCoordinate{Row: 7, Col: 0},
		},
	},
}

func initializeScrabbleBoard() ScrabbleBoard {

	sb := ScrabbleBoard{}

	// Initialize board with plain squares
	for i, row := range sb {
		for j := range row {
			sb[i][j] = Square{
				SquareType: "plain",
			}
		}
	}

	// Place remaining squares based on coordinates
	for s := range squareTypes {
		st := squareTypes[s]
		for _, sc := range st.Coordinates {

			squ := Square{
				SquareType: st.Name,
			}

			// Quadrant 1
			sb[sc.Row][columnCount-1-sc.Col] = squ

			// Quadrant 2
			sb[sc.Row][sc.Col] = squ

			// Quadrant 3
			sb[rowCount-1-sc.Row][sc.Col] = squ

			// Quadrant 4
			sb[rowCount-1-sc.Row][columnCount-1-sc.Col] = squ
		}
	}

	return sb
}

func (sb ScrabbleBoard) print() {
	for _, row := range sb {
		for _, col := range row {
			fmt.Print(col.SquareType + " ")
		}
		fmt.Println()
	}
}
