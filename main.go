package main

import (
	"fmt"
	"runtime/debug"
	"time"
	"unsafe"

	s2 "github.com/fr0nch/go-plugify-s2sdk/v2"
	"github.com/untrustedmodders/go-plugify"
)

var Plugin *SPlugin

func init() {
	//utils.CreateManifest("Advert", "1.0.0", "OkyHek", []string{"s2sdk"})
	Plugin = NewPlugin()

	plugify.OnPluginStart(Plugin.OnPluginStart)
	plugify.OnPluginEnd(Plugin.OnPluginEnd)
	plugify.OnPluginPanic(Plugin.OnPluginPanic)
}

func (pl *SPlugin) OnPluginStart() {
	iface := s2.FindInterface("NetworkSystemVersion001")
	if iface == 0 {
		panic("interface nil")
	}
	pl.NetworkSystem = unsafe.Pointer(iface)

	var err error
	Plugin.Config, err = ReadConfig()
	if err != nil {
		fmt.Printf("[Advert] ReadConfig: %s\n", err)
		return
	}
	MSGDebug("Advert ReadConfig: %v", Plugin.Config)

	s2.OnServerActivate_Register(pl.OnServerActivate)
}

func (pl *SPlugin) OnPluginEnd() {
	MSGDebug("Advert OnPluginEnd")

	s2.OnServerActivate_Unregister(pl.OnServerActivate)
}

func (pl *SPlugin) OnPluginPanic() []byte {
	return debug.Stack() // workaround for could not import runtime/debug inside plugify package
}

func (pl *SPlugin) OnServerActivate() { // it`s OnMapStart
	if time.Now().Unix() > pl.MapLoadTime+int64(3) {
		err := LoadAdvert()
		if err != nil {
			fmt.Printf("[Advert] LoadAdvert: %s\n", err)
			return
		}
	}
	pl.MapLoadTime = time.Now().Unix()

	pl.CurrentIndex = 0
	MSGDebug("Advert OnServerActivate. Index: %d | Adverts len %d", pl.CurrentIndex, len(pl.Adverts))

	if len(pl.Adverts) > 0 {
		s2.CreateTimer(Plugin.Config.TimerInterval, pl.OnTimerAdvert, s2.TimerFlag_NoMapChange|s2.TimerFlag_Repeat, []any{})
	}
}

func main() {}
