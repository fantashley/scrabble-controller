package wordgameserver

import (
	"reflect"
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

	count := make(map[string]int)
	for _, row := range initializedBoard {
		for _, square := range row {
			count[square.SquareType]++
		}
	}

	if !reflect.DeepEqual(occurrences, count) {
		t.Error("Scrabble board does not have correct square counts")
	}
}
