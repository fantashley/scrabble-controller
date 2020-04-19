package wordgameserver

import "errors"

func (sg *ScrabbleGame) executePlay(j GamePlayRequest) error {
	playerTurn := sg.TurnCount % len(sg.Players)
	if playerTurn != sg.Players[j.PlayerID].Number {
		return errors.New("Playing out of turn. Expected Player " + string(playerTurn))
	} else if len(j.Tiles) > 7 {
		return errors.New("Cannot play more than 7 tiles")
	}

	if j.Swap {
		return sg.swapTiles(j)
	}

	return nil
}

func (sg *ScrabbleGame) swapTiles(j GamePlayRequest) error {
	if len(j.Tiles) > len(sg.TileBag) {
		return errors.New("Not enough tiles available for swap")
	}

	// Remove tiles from player's hand
	cp := sg.Players[j.PlayerID]
	err := removeTiles(cp, j.Tiles)
	if err != nil {
		return err
	}

	// Deal new tiles to player
	dealTiles(cp, &sg.TileBag, len(j.Tiles))

	// Add swapped tiles to bag and shuffle
	sg.TileBag = append(sg.TileBag, j.Tiles...)
	sg.TileBag.shuffle()

	return nil
}
