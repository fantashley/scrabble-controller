package wordgameserver

import "errors"

func (wg *WordGame) executePlay(j GamePlayRequest) error {
	playerTurn := wg.TurnCount % len(wg.Players)
	if playerTurn != wg.Players[j.PlayerID].Number {
		return errors.New("Playing out of turn. Expected Player " + string(playerTurn))
	} else if len(j.Tiles) > 7 {
		return errors.New("Cannot play more than 7 tiles")
	}

	if j.Swap {
		return wg.swapTiles(j)
	}

	return nil
}

func (wg *WordGame) swapTiles(j GamePlayRequest) error {
	if len(j.Tiles) > len(wg.TileBag) {
		return errors.New("Not enough tiles available for swap")
	}

	// Remove tiles from player's hand
	cp := wg.Players[j.PlayerID]
	err := removeTiles(cp, j.Tiles)
	if err != nil {
		return err
	}

	// Deal new tiles to player
	dealTiles(cp, &wg.TileBag, len(j.Tiles))

	// Add swapped tiles to bag and shuffle
	wg.TileBag = append(wg.TileBag, j.Tiles...)
	wg.TileBag.shuffle()

	return nil
}
