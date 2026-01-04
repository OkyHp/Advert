package main

import (
	"fmt"
	"strings"

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

		buff = strings.ReplaceAll(buff, "{IP}", GetServerIP())
		buff = strings.ReplaceAll(buff, "{PORT}", GetServerPort())

		advert.MsgText[lang] = buff
	}
}

func ReplacePlaceholders() {
	for index, _ := range Plugin.Adverts {
		for lang, _ := range Plugin.Adverts[index].MsgText {
			Plugin.Adverts[index].MsgText[lang] =
				strings.ReplaceAll(Plugin.Adverts[index].MsgText[lang], "{MAP}", s2.GetCurrentMap())
		}
	}
}

func GetServerIP() string {
	if Plugin.Config.ServerIp == "" {
		cvar := s2.FindConVar("hostip")
		if cvar != 0 {
			return Uint32ToIPv4(s2.GetConVarUInt32(cvar))
		}
	}

	return Plugin.Config.ServerIp
}

func GetServerPort() string {
	cvar := s2.FindConVar("hostport")
	if cvar != 0 {
		return s2.GetConVarString(cvar)
	}

	return ""
}

func Uint32ToIPv4(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24),
		byte(ip>>16),
		byte(ip>>8),
		byte(ip),
	)
}

func GetClientLanguageEx(playerSlot int32) string {
	lang := "en"

	buff := s2.GetClientLanguage(playerSlot)
	if len(buff) >= 2 {
		lang = buff[:2]
	}

	return lang
}
