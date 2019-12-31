package main

import "fmt"

type squareCoordinate struct {
	row int
	col int
}

type SquareType struct {
	Name             string             `json:"name"`
	LetterMultiplier int                `json:"letterMultiplier"`
	WordMultiplier   int                `json:"wordMultiplier"`
	Coordinates      []squareCoordinate `json:"-"`
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
		Coordinates: []squareCoordinate{
			squareCoordinate{row: 7, col: 7},
		},
	},
	"doubleLetter": SquareType{
		Name:             "doubleLetter",
		LetterMultiplier: 2,
		WordMultiplier:   1,
		Coordinates: []squareCoordinate{
			squareCoordinate{row: 0, col: 3},
			squareCoordinate{row: 2, col: 6},
			squareCoordinate{row: 3, col: 0},
			squareCoordinate{row: 3, col: 7},
			squareCoordinate{row: 6, col: 2},
			squareCoordinate{row: 6, col: 6},
			squareCoordinate{row: 7, col: 3},
		},
	},
	"doubleWord": SquareType{
		Name:             "doubleWord",
		LetterMultiplier: 1,
		WordMultiplier:   2,
		Coordinates: []squareCoordinate{
			squareCoordinate{row: 1, col: 1},
			squareCoordinate{row: 2, col: 2},
			squareCoordinate{row: 3, col: 3},
			squareCoordinate{row: 4, col: 4},
		},
	},
	"tripleLetter": SquareType{
		Name:             "tripleLetter",
		LetterMultiplier: 3,
		WordMultiplier:   1,
		Coordinates: []squareCoordinate{
			squareCoordinate{row: 1, col: 5},
			squareCoordinate{row: 5, col: 1},
			squareCoordinate{row: 5, col: 5},
		},
	},
	"tripleWord": SquareType{
		Name:             "tripleWord",
		LetterMultiplier: 1,
		WordMultiplier:   3,
		Coordinates: []squareCoordinate{
			squareCoordinate{row: 0, col: 0},
			squareCoordinate{row: 0, col: 7},
			squareCoordinate{row: 7, col: 0},
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
			sb[sc.row][columnCount-1-sc.col] = squ

			// Quadrant 2
			sb[sc.row][sc.col] = squ

			// Quadrant 3
			sb[rowCount-1-sc.row][sc.col] = squ

			// Quadrant 4
			sb[rowCount-1-sc.row][columnCount-1-sc.col] = squ
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
