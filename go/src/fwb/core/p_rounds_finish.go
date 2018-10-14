package core

import (
	"er"
	"math/rand"
)

func prfInit(me *gameImp) *er.Err {
	if me.gd.Round >= me.gd.MinRounds {
		if rand.Intn(6)+1 > 3 {
			return me.gotoPhase(_P_GAME_SETTLEMENT)
		}
	}

	return me.gotoPhase(_P_ROUNDS_START)
}
