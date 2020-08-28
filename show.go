package main

import (
	"encoding/json"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleShow(m *tb.Message) {
	variables := []Variable{}
	gl.Orm.Find(&variables)
	fixeds := []Fixed{}
	gl.Orm.Find(&fixeds)

	v, _ := json.Marshal(variables)
	f, _ := json.Marshal(fixeds)
	gl.Bot.Send(m.Sender, "`"+string(v)+string(f)+"`", tb.ModeMarkdown)
}
