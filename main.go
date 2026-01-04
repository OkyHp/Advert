package main

import (
	"fmt"
	"runtime/debug"

	s2 "github.com/OkyHp/plg_utils/s2sdk"
	"github.com/untrustedmodders/go-plugify"
)

var Plugin *ResetScorePlugin

func init() {
	//utils.CreateManifest("Advert", "1.0.0", "OkyHek", []string{"s2sdk"})

	Plugin = NewResetScorePlugin()

	plugify.OnPluginStart(Plugin.OnPluginStart)
	plugify.OnPluginEnd(Plugin.OnPluginEnd)
	plugify.OnPluginPanic(Plugin.OnPluginPanic)
}

func (rs *ResetScorePlugin) OnPluginStart() {
	var err error

	rs.Config, err = ReadConfig()
	if err != nil {
		s2.PrintToServer(fmt.Sprintf("CONFIG: %s", err))
		return
	}

	err = InitDatabase()
	if err != nil {
		s2.PrintToServer(fmt.Sprintf("DATABASE: %s", err))
		return
	}

	s2.OnServerActivate_Register(rs.OnServerActivate)
}

func (rs *ResetScorePlugin) OnPluginEnd() {
	s2.OnServerActivate_Unregister(rs.OnServerActivate)
}

func (rs *ResetScorePlugin) OnPluginPanic() []byte {
	return debug.Stack() // workaround for could not import runtime/debug inside plugify package
}

func (rs *ResetScorePlugin) OnServerActivate() { // it`s OnMapStart
	Plugin.CurrentIndex = 0
	ReplacePlaceholders()

	if len(Plugin.Adverts) > 0 {
		s2.CreateTimer(Plugin.Config.TimerInterval, rs.OnTimerAdvert, s2.TimerFlag_NoMapChange, []any{})
	}
}

func main() {}
