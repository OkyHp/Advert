package main

type ResetScorePlugin struct {
	Config       ConfigData
	Adverts      []AdvertsData
	CurrentIndex uint32
}

func NewResetScorePlugin() *ResetScorePlugin {
	return &ResetScorePlugin{}
}
