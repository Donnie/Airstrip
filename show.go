package main

import (
	"encoding/json"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleShow(m *tb.Message) {
	recs := []Record{}
	gl.Orm.Find(&recs)

	r, _ := json.Marshal(recs)
	gl.Bot.Send(m.Sender, "`"+string(r)+"`", tb.ModeMarkdown)
}
