package main

import "fmt"

type squareCoordinate struct {
	row int
	col int
}

type squareType struct {
	name             string
	letterMultiplier int
	wordMultiplier   int
	coordinates      []squareCoordinate
}

type square struct {
	squareType string
	tile
}

const rowCount int = 15
const columnCount int = 15

type scrabbleBoard [rowCount][columnCount]square

var squareTypes = map[string]squareType{
	"plain": squareType{
		name:             "plain",
		letterMultiplier: 1,
		wordMultiplier:   1,
	},
	"star": squareType{
		name:             "star",
		letterMultiplier: 1,
		wordMultiplier:   1,
		coordinates: []squareCoordinate{
			squareCoordinate{row: 7, col: 7},
		},
	},
	"doubleLetter": squareType{
		name:             "doubleLetter",
		letterMultiplier: 2,
		wordMultiplier:   1,
		coordinates: []squareCoordinate{
			squareCoordinate{row: 0, col: 3},
			squareCoordinate{row: 2, col: 6},
			squareCoordinate{row: 3, col: 0},
			squareCoordinate{row: 3, col: 7},
			squareCoordinate{row: 6, col: 2},
			squareCoordinate{row: 6, col: 6},
			squareCoordinate{row: 7, col: 3},
		},
	},
	"doubleWord": squareType{
		name:             "doubleWord",
		letterMultiplier: 1,
		wordMultiplier:   2,
		coordinates: []squareCoordinate{
			squareCoordinate{row: 1, col: 1},
			squareCoordinate{row: 2, col: 2},
			squareCoordinate{row: 3, col: 3},
			squareCoordinate{row: 4, col: 4},
		},
	},
	"tripleLetter": squareType{
		name:             "tripleLetter",
		letterMultiplier: 3,
		wordMultiplier:   1,
		coordinates: []squareCoordinate{
			squareCoordinate{row: 1, col: 5},
			squareCoordinate{row: 5, col: 1},
			squareCoordinate{row: 5, col: 5},
		},
	},
	"tripleWord": squareType{
		name:             "tripleWord",
		letterMultiplier: 1,
		wordMultiplier:   3,
		coordinates: []squareCoordinate{
			squareCoordinate{row: 0, col: 0},
			squareCoordinate{row: 0, col: 7},
			squareCoordinate{row: 7, col: 0},
		},
	},
}

func (sb *scrabbleBoard) initialize() {

	// Initialize board with plain squares
	for i, row := range sb {
		for j := range row {
			sb[i][j] = square{
				squareType: "plain",
			}
		}
	}

	// Place remaining squares based on coordinates
	for s := range squareTypes {
		st := squareTypes[s]
		for _, sc := range st.coordinates {

			squ := square{
				squareType: st.name,
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
}

func (sb scrabbleBoard) print() {
	for _, row := range sb {
		for _, col := range row {
			fmt.Print(col.squareType + " ")
		}
		fmt.Println()
	}
}
