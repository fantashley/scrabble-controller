package wordgameserver

import "fmt"

// SquareCoordinate represents a coordinate of a Scrabble board
type SquareCoordinate struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// SquareType represents the underlying types of squares on a Scrabble board
type SquareType struct {
	Name             string             `json:"name"`             // type of square, such as plain or tripleWord
	LetterMultiplier int                `json:"letterMultiplier"` // multiplier for letters on square
	WordMultiplier   int                `json:"wordMultiplier"`   // multiplier for words on square
	Coordinates      []SquareCoordinate `json:"-"`                // second quadrant symmetrical coordinates for square type
}

// Square represents the squares on a Scrabble Board
type Square struct {
	SquareType string `json:"type"`
	Tile       `json:"tile,omitempty"`
}

const rowCount int = 15
const columnCount int = 15

// ScrabbleBoard represents the board containing a grid of Squares
type ScrabbleBoard [rowCount][columnCount]Square

var initializedBoard = initializeScrabbleBoard()

// squareTypes is a definition of the possible square types and the values they
// hold
var squareTypes = map[string]SquareType{
	"plain": {
		Name:             "plain",
		LetterMultiplier: 1,
		WordMultiplier:   1,
	},
	"star": {
		Name:             "star",
		LetterMultiplier: 1,
		WordMultiplier:   1,
		Coordinates: []SquareCoordinate{
			{Row: 7, Col: 7},
		},
	},
	"doubleLetter": {
		Name:             "doubleLetter",
		LetterMultiplier: 2,
		WordMultiplier:   1,
		Coordinates: []SquareCoordinate{
			{Row: 0, Col: 3},
			{Row: 2, Col: 6},
			{Row: 3, Col: 0},
			{Row: 3, Col: 7},
			{Row: 6, Col: 2},
			{Row: 6, Col: 6},
			{Row: 7, Col: 3},
		},
	},
	"doubleWord": {
		Name:             "doubleWord",
		LetterMultiplier: 1,
		WordMultiplier:   2,
		Coordinates: []SquareCoordinate{
			{Row: 1, Col: 1},
			{Row: 2, Col: 2},
			{Row: 3, Col: 3},
			{Row: 4, Col: 4},
		},
	},
	"tripleLetter": {
		Name:             "tripleLetter",
		LetterMultiplier: 3,
		WordMultiplier:   1,
		Coordinates: []SquareCoordinate{
			{Row: 1, Col: 5},
			{Row: 5, Col: 1},
			{Row: 5, Col: 5},
		},
	},
	"tripleWord": {
		Name:             "tripleWord",
		LetterMultiplier: 1,
		WordMultiplier:   3,
		Coordinates: []SquareCoordinate{
			{Row: 0, Col: 0},
			{Row: 0, Col: 7},
			{Row: 7, Col: 0},
		},
	},
}

// initializeScrabbleBoard places the squares on the board
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

// print outputs the scrabble board square contents
func (sb ScrabbleBoard) print() {
	for _, row := range sb {
		for _, col := range row {
			fmt.Print(col.SquareType + " ")
		}
		fmt.Println()
	}
}
