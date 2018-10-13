package actn

import (
	"encoding/json"
	"er"
	"fmt"
	"fwb"
	"sgs"
)

type actnWeave struct {
	amount   int
	playerID int
}

func actnWeaveParser(command sgs.Command) fwb.Action {
	payload, err := json.Marshal(command.Payload)

	if err != nil {
		return nil
	}

	var weave struct {
		ActionID int
		Payload  int
	}

	err = json.Unmarshal(payload, &weave)

	if err != nil {
		return nil
	}

	return &actnWeave{
		amount:   weave.Payload,
		playerID: command.Source,
	}
}

func (me *actnWeave) String() string {
	return fmt.Sprintf("[Action %v from Player %v, Amount %v]", _actionNames[ACTN_WEAVE], me.playerID, me.amount)
}

func (me *actnWeave) ID() int {
	return ACTN_WEAVE
}

func (me *actnWeave) getCost() fwb.PlayerData {
	cost := make(fwb.PlayerData, fwb.PD_MAX)
	cost[fwb.PD_PAWNS] = -1
	cost[fwb.PD_PT_WOOL] = me.amount * -2
	return cost
}

func (me *actnWeave) ValidateAgainst(gd *fwb.GameData) bool {
	if !hasCardSlots(gd, ACTN_WEAVE) {
		return false
	}

	i := gd.GetPDIndex(me.playerID)

	if i < 0 {
		return false
	}

	p := gd.PData[i]

	cost := me.getCost()

	res := fwb.PDAdd(p, cost)

	return res.AllAboveZero()
}

func (me *actnWeave) Do(gd *fwb.GameData) *er.Err {
	i := gd.GetPDIndex(me.playerID)
	if i < 0 {
		return er.Throw(fwb.E_DOACTION_INVALID_CLIENTID, er.EInfo{
			"details":  "invalid player ID for do action",
			"playerID": me.playerID,
		})
	}

	cost := me.getCost()
	gain := make(fwb.PlayerData, fwb.PD_MAX)
	gain[fwb.PD_PT_SWEATER] = me.amount

	p := gd.PData[me.playerID]
	p = fwb.PDAdd(cost, p)
	p = fwb.PDAdd(gain, p)

	gd.PData[me.playerID] = p
	return checkCard(gd, ACTN_TRAIN, me.playerID, 1)
}