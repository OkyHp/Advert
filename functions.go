package main

/*
#include <stdint.h>

typedef void* (*UpdatePublicIPFn)(void*);

static inline void* Call_UpdatePublicIP(void* fn, void* thisptr) {
    return ((UpdatePublicIPFn)fn)(thisptr);
}
*/
import "C"

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	s2 "github.com/OkyHp/plg_utils/s2sdk"
)

var placeholders = map[string]string{
	"Default":     "\x01",
	"White":       "\x01",
	"DarkRed":     "\x02",
	"Green":       "\x04",
	"LightYellow": "\x09",
	"LightBlue":   "\x0B",
	"Olive":       "\x05",
	"Lime":        "\x06",
	"Red":         "\x07",
	"LightPurple": "\x03",
	"Purple":      "\x0E",
	"Grey":        "\x08",
	"Yellow":      "\x09",
	"Gold":        "\x10",
	"Silver":      "\x0A",
	"Blue":        "\x0B",
	"DarkBlue":    "\x0C",
	"BlueGrey":    "\x0A",
	"Magenta":     "\x0E",
	"LightRed":    "\x0F",
	"Orange":      "\x10",
	"Darkred":     "\x02", // Obsolete, но оставляем для совместимости
	"NewLine":     "\u2029",
}

func ReplaceStaticPlaceholders(advert *AdvertsData) {
	for lang, _ := range advert.MsgText {
		buff := advert.MsgText[lang]

		for name, value := range placeholders {
			tag := fmt.Sprintf("{%s}", name)
			buff = strings.ReplaceAll(buff, tag, value)
		}

		buff = strings.ReplaceAll(buff, "{Ip}", GetServerIP())
		buff = strings.ReplaceAll(buff, "{Port}", GetServerPort())

		advert.MsgText[lang] = buff
	}
}

func ReplacePlaceholders(adverts []AdvertsData) []AdvertsData {
	buff := adverts

	for index, _ := range buff {
		for lang, _ := range buff[index].MsgText {
			buff[index].MsgText[lang] =
				strings.ReplaceAll(buff[index].MsgText[lang], "{Map}", s2.GetCurrentMap())
		}
	}

	return buff
}

func GetServerIP() string {
	if Plugin.Config.ServerIp != "" {
		return Plugin.Config.ServerIp
	}

	fnPtr := getVFunc(Plugin.NetworkSystem, 32)
	if fnPtr == nil {
		panic("fnPtr nil")
	}

	netadr := unsafe.Pointer(C.Call_UpdatePublicIP(
		unsafe.Pointer(fnPtr),
		Plugin.NetworkSystem,
	))
	if netadr == nil {
		panic("netadr nil")
	}

	ip := (*[4]byte)(unsafe.Add(netadr, 4))

	return fmt.Sprintf("Public IP: %d.%d.%d.%d\n",
		ip[0], ip[1], ip[2], ip[3],
	)
}

func GetServerPort() string {
	cvar := s2.FindConVar("hostport")
	if cvar != 0 {
		return strconv.Itoa(int(s2.GetConVarInt32(cvar)))
	}

	return ""
}

func getVFunc(obj unsafe.Pointer, index int) unsafe.Pointer {
	vtbl := *(*unsafe.Pointer)(obj)
	return *(*unsafe.Pointer)(unsafe.Add(vtbl, uintptr(index)*unsafe.Sizeof(uintptr(0))))
}

func MSGDebug(text string, args ...interface{}) {
	if Plugin.Config.Debug == true {
		s2.PrintToServer(fmt.Sprintf("[DEBUG] "+text, args...))
	}
}
