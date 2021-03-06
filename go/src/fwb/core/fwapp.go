package core

import (
	"er"
	"fwb"
	"hlf"
	"sgs"
	"strconv"
)

var _debugGlobalApp *fwAppImp

//FwAppBuildFunc hook up with the SGS server
func FwAppBuildFunc() sgs.App {
	return &fwAppImp{}
}

var _execMap = map[int]func(*fwAppImp, sgs.Command) *er.Err{
	sgs.CMD_TICK:           onTick,
	sgs.CMD_APP_RUN:        forwardToGame,
	sgs.CMD_FORWARD_TO_APP: forwardToPlayer,
}

type fwAppImp struct {
	g  fwb.Game
	pm map[int]fwb.PlayerAgent
	mp mockPlayers
	s  sgs.Session
	lg hlf.Logger
}

func (me *fwAppImp) Init(s sgs.Session, clients []int, profile string) *er.Err {
	_debugGlobalApp = me

	me.s = s

	me.lg = s.GetLogger()

	me.pm = make(map[int]fwb.PlayerAgent)

	for _, c := range clients {
		me.lg.Dbg("Try to make player: clientID 0x%x, name %v", c, me.s.GetClientName(c))
		me.pm[c] = makePlayer(me, c, me.s.GetClientName(c))
	}

	game, err := makeGame(me, profile)
	if err.Importance() >= er.IMPT_DEGRADE {
		return err
	}

	me.g = game

	me.mp = mockPlayers{}

	me.mp.init(game)

	return err
}

func (me *fwAppImp) GetSession() sgs.Session {
	return me.s
}

func (me *fwAppImp) SendCommand(command sgs.Command) *er.Err {
	exec, found := _execMap[command.ID]

	if found {
		return exec(me, command)
	}

	if command.InCategory(sgs.CMD_C_CLIENT_TO_APP) {
		return forwardToGame(me, command)
	}

	return er.Throw(fwb.E_INVALID_CMD, er.EInfo{
		"details": "command is not supposed to be received, no clue about the execution",
		"command": command.HexID(),
	}).To(me.lg)
}

func (me *fwAppImp) GetLogger() hlf.Logger {
	return me.lg
}

func onTick(app *fwAppImp, command sgs.Command) *er.Err {
	err := app.mp.execOne()
	if err.Importance() >= er.IMPT_DEGRADE {
		return err
	}

	return forwardToGame(app, command)
}

func forwardToGame(app *fwAppImp, command sgs.Command) *er.Err {
	return app.g.SendCommand(command)
}

func forwardToPlayer(app *fwAppImp, command sgs.Command) *er.Err {
	p, found := app.pm[command.Who]
	if !found {
		return er.Throw(fwb.E_CMD_INVALID_CLIENT, er.EInfo{
			"details": "cannot forward the command to a player, invalid client ID " + strconv.Itoa(command.Who) + ", command " + command.HexID(),
		}).To(app.lg)
	}

	return p.SendCommand(command)
}

func (me *fwAppImp) SendAllPlayers(command sgs.Command) *er.Err {
	var err *er.Err

	for _, p := range me.pm {
		err = err.Push(me.SendToPlayer(p.ID(), command))
		if err.Importance() >= er.IMPT_DEGRADE {
			return err
		}
	}

	return err
}

func (me *fwAppImp) SendToPlayer(playerID int, command sgs.Command) *er.Err {
	player, ok := me.pm[playerID]
	if !ok {
		return er.Throw(fwb.E_CMD_INVALID_CLIENT, er.EInfo{
			"details": "send command to player with invalid client ID",
			"ID":      playerID,
		}).To(me.lg)
	}

	err := player.SendCommand(command)
	if err.Importance() > er.IMPT_RECOVERABLE && err.Importance() < er.IMPT_DEGRADE {
		me.SendToMockPlayer(playerID, command)
	}
	return err
}

func (me *fwAppImp) SendToGame(command sgs.Command) *er.Err {
	return me.g.SendCommand(command)
}

func (me *fwAppImp) SendToSession(command sgs.Command) *er.Err {
	me.s.CmdChan() <- command
	return nil
}

func (me *fwAppImp) GetPlayers() []fwb.PlayerAgent {
	players := make([]fwb.PlayerAgent, 0, len(me.pm))
	for _, p := range me.pm {
		players = append(players, p)
	}
	return players
}

func (me *fwAppImp) SendToMockPlayer(playerID int, command sgs.Command) {
	me.mp.sendCommand(playerID, command)
}

func (me *fwAppImp) GetPlayer(playerID int) fwb.PlayerAgent {
	for _, player := range me.pm {
		if player.ID() == playerID {
			return player
		}
	}
	return nil
}
