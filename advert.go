package main

import s2 "github.com/OkyHp/plg_utils/s2sdk"

func (rs *ResetScorePlugin) OnTimerAdvert(timer uint32, userData []any) {
	advert := rs.Adverts[rs.CurrentIndex]

	for i := int32(1); i < s2.GetMaxClients()+1; i++ {
		if s2.IsClientInGame(i) && !s2.IsFakeClient(i) { // IsVipClient
			lang := GetClientLanguageEx(i)
			if advert.MsgText[lang] != "" {
				switch advert.MsgType {
				case "CHAT":
					s2.PrintToChat(i, advert.MsgText[lang])
					break
				case "CENTER":
					s2.PrintCenterText(i, advert.MsgText[lang])
					break
				case "ALERT":
					s2.PrintAlertText(i, advert.MsgText[lang])
					break
				case "HTML":
					s2.PrintCentreHtml(i, advert.MsgText[lang], rs.Config.HtmlMsgDuration)
					break
				}
			}
		}
	}

	rs.CurrentIndex = (rs.CurrentIndex + 1) % uint32(len(rs.Adverts))
}
