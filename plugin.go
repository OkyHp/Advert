package main

import "unsafe"

type AdvertPlugin struct {
	Config ConfigData

	Adverts      []AdvertsData
	CurrentIndex uint32

	DatabaseInit  bool
	MapLoadTime   int64
	NetworkSystem unsafe.Pointer
}

func NewAdvertPlugin() *AdvertPlugin {
	return &AdvertPlugin{}
}
