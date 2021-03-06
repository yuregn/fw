package sgs

import (
	"er"
)

//App interface for applications
type App interface {
	Init(s Session, clientIDs []int, profile string) *er.Err
	SendCommand(command Command) *er.Err
}

//AppBuildFunc function to build an application
type AppBuildFunc func() App
