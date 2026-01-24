package main

type SPlugin struct {
	Config ConfigData

	Adverts      []AdvertsData
	CurrentIndex uint32

	DatabaseInit bool
}

func NewPlugin() *SPlugin {
	return &SPlugin{}
}
