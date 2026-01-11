package main

import "unsafe"

type SPlugin struct {
	Config ConfigData

	Adverts      []AdvertsData
	CurrentIndex uint32

	DatabaseInit  bool
	MapLoadTime   int64
	NetworkSystem unsafe.Pointer
}

func NewPlugin() *SPlugin {
	return &SPlugin{}
}
