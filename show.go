package main

import (
	"encoding/json"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleShow(m *tb.Message) {
	items := []Item{}
	gl.Orm.Find(&items)

	j, _ := json.Marshal(items)
	gl.Bot.Send(m.Sender, "`"+string(j)+"`", tb.ModeMarkdown)
}
