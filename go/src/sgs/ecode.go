package sgs

import "er"

const (
	_E_SGS_RUNNER = 0x1000

	_E_LOAD_CONF_FAIL = _E_SGS_RUNNER | er.IMPT_RECOVERABLE | er.ET_SERVICE | 0x1
)

const (
	_E_SGS_SSVR = 0x2000

	//E_JOIN_SESSION_INVALID_CLIENT send join session request without valid client
	_E_JOIN_SESSION_INVALID_CLIENT = _E_SGS_SSVR | er.IMPT_REMARKABLE | er.ET_INTERACTION | er.EI_INVALID_REQUEST | 0x1

	_E_INVALID_SERVER_PARAM = _E_SGS_SSVR | er.IMPT_UNRECOVERABLE | er.ET_INTERNAL | er.EI_ILLEGAL_PARAMETER | 0x2

	_E_BIND_CONN_INVALID_CLIENT = _E_SGS_SSVR | er.IMPT_REMARKABLE | er.ET_INTERACTION | er.EI_INVALID_REQUEST | 0x3

	_E_SESSION_INVALID_COMMAND = _E_SGS_SSVR | er.IMPT_REMARKABLE | er.ET_INTERNAL | er.EI_ILLEGAL_PARAMETER | 0x4

	_E_SESSION_INVALID_CLIENT = _E_SGS_SSVR | er.IMPT_REMARKABLE | er.ET_INTERNAL | er.EI_ILLEGAL_PARAMETER | 0x5

	_E_CLIENT_CONNECTION_FAIL = _E_SGS_SSVR | er.IMPT_REMARKABLE | er.ET_INTERACTION | er.EI_NETWORK | 0x6

	_E_NO_DUAL_CONNECTION_SUPPORT = _E_SGS_SSVR | er.IMPT_DEGRADE | er.ET_INTERACTION | er.EI_NETWORK | 0x7

	_E_CLIENT_ALREADY_JOIN_SESSION = _E_SGS_SSVR | er.IMPT_REMARKABLE | er.ET_INTERACTION | er.EI_INVALID_REQUEST | 0x8

	_E_QUIT_SESSION_QUEUE_INVALID = _E_SGS_SSVR | er.IMPT_ACCEPTIBLE | er.ET_INTERACTION | er.EI_INVALID_REQUEST | 0x9

	_E_CLIENT_FAILE_RECONNECT = _E_SGS_SSVR | er.IMPT_ACCEPTIBLE | er.ET_INTERACTION | er.EI_INVALID_REQUEST | 0x10
)
