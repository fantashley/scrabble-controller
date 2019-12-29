package main

import (
	"fmt"
	"os"
)

func main() {

	playerNames := []string{
		"Ashley",
		"Emily",
		"Kelsey",
		"Michelle",
	}

	game, err := createScrabbleGame(playerNames)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	printTiles(game.tileBag)
	fmt.Println("Length of tile bag:", len(game.tileBag))

	for _, player := range game.players {
		fmt.Print(player.name, ": ")
		printTiles(player.tiles)
	}
}
