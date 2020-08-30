package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleContext(m *tb.Message) {
	userID := int64(m.Sender.ID)

	convo := &Convo{}
	res := gl.Orm.
		Where("user_id = ?", userID).
		Last(convo)

	if res.Error != nil {
		gl.Bot.Send(m.Sender, "Sorry didn't get you! /help")
		return
	}

	que := convo.expectWhat(gl, m.Text)
	gl.Bot.Send(m.Sender, que)
}

func genQues(ask, form string) string {
	return fmt.Sprintf("What is the %s of the %s?", ask, form)
}
