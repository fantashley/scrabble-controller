package main

type tile struct {
	letter byte
	count  int
	value  int
}

type squareType int

type squareCoordinate struct {
	row int
	col int
}

const (
	plain        squareType = 0
	star         squareType = 1
	doubleLetter squareType = 2
	doubleWord   squareType = 3
	tripleLetter squareType = 4
	tripleWord   squareType = 5
)

var squareLayout = map[squareType][]squareCoordinate{
	star: []squareCoordinate{
		squareCoordinate{row: 7, col: 7},
	},
	doubleLetter: []squareCoordinate{
		squareCoordinate{row: 0, col: 3},
		squareCoordinate{row: 2, col: 6},
		squareCoordinate{row: 3, col: 0},
		squareCoordinate{row: 3, col: 7},
		squareCoordinate{row: 6, col: 2},
		squareCoordinate{row: 6, col: 6},
		squareCoordinate{row: 7, col: 3},
	},
	doubleWord: []squareCoordinate{
		squareCoordinate{row: 1, col: 1},
		squareCoordinate{row: 2, col: 2},
		squareCoordinate{row: 3, col: 3},
		squareCoordinate{row: 4, col: 4},
	},
	tripleLetter: []squareCoordinate{
		squareCoordinate{row: 1, col: 5},
		squareCoordinate{row: 5, col: 1},
		squareCoordinate{row: 5, col: 5},
	},
	tripleWord: []squareCoordinate{
		squareCoordinate{row: 0, col: 0},
		squareCoordinate{row: 0, col: 7},
		squareCoordinate{row: 7, col: 0},
	},
}
