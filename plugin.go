package main

import "unsafe"

type ResetScorePlugin struct {
	Config ConfigData

	Adverts      []AdvertsData
	CurrentIndex uint32

	NetworkSystem unsafe.Pointer
}

func NewResetScorePlugin() *ResetScorePlugin {
	return &ResetScorePlugin{}
}
