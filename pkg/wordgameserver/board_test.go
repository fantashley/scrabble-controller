package wordgameserver

import (
	"testing"
)

func TestBoard(t *testing.T) {
	occurrences := map[string]int{
		"star":         1,
		"doubleLetter": 24,
		"doubleWord":   16,
		"tripleLetter": 12,
		"tripleWord":   8,
		"plain":        164,
	}

	for _, row := range initializedBoard {
		for _, square := range row {
			if squareType, ok := occurrences[square.SquareType]; !ok {
				t.Fatalf("Too many of square type %v", squareType)
			}
			occurrences[square.SquareType]--
			if occurrences[square.SquareType] == 0 {
				delete(occurrences, square.SquareType)
			}
		}
	}
	if len(occurrences) != 0 {
		t.Fatalf("Length of occurrences is %v but should be 0", len(occurrences))
	}
}
