package main

import (
	s2 "github.com/OkyHp/plg_utils/s2sdk"
	"github.com/OkyHp/plg_utils/utils"
)

func (pl *AdvertPlugin) OnTimerAdvert(timer uint32, userData []any) {
	advert := pl.Adverts[pl.CurrentIndex]

	for i := int32(0); i < s2.GetMaxClients()+1; i++ {
		if s2.IsClientInGame(i) && !s2.IsFakeClient(i) { // IsVipClient
			lang := utils.GetClientLanguageEx(i)
			if advert.MsgText[lang] != "" {
				switch advert.MsgType {
				case "CHAT":
					s2.PrintToChat(i, " "+advert.MsgText[lang])
					break
				case "CENTER":
					s2.PrintCenterText(i, advert.MsgText[lang])
					break
				case "ALERT":
					s2.PrintAlertText(i, advert.MsgText[lang])
					break
				case "HTML":
					s2.PrintCentreHtml(i, advert.MsgText[lang], pl.Config.HtmlMsgDuration)
					break
				}
			}
		}
	}

	pl.CurrentIndex = (pl.CurrentIndex + 1) % uint32(len(pl.Adverts))
}
