package main

import (
	"fmt"
	"strconv"
	"strings"

	s2 "github.com/fr0nch/go-plugify-s2sdk/v2"
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
		buff = strings.ReplaceAll(buff, "{Map}", s2.GetCurrentMap())

		advert.MsgText[lang] = buff
	}
}

func GetServerIP() string {
	if Plugin.Config.ServerIp != "" {
		return Plugin.Config.ServerIp
	}

	return s2.GetPublicAddress(true)
}

func GetServerPort() string {
	cvar := s2.FindConVar("hostport")
	if cvar != 0 {
		return strconv.Itoa(int(s2.GetConVarInt32(cvar)))
	}

	return ""
}

func MSGDebug(message string, args ...any) {
	if Plugin.Config.Debug == true {
		fmt.Printf("[DEBUG] "+message+"\n", args...)
	}
}
