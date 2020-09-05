package main

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleText(m *tb.Message) {
	gl.handleContext(m.Sender, m.Text)
}

func (gl *Global) handleCallback(m *tb.Callback) {
	gl.handleContext(m.Sender, strings.TrimSpace(m.Data))
	gl.Bot.Respond(m, &tb.CallbackResponse{
		CallbackID: m.ID,
		Text:       "Cool!",
	})
}

func (gl *Global) handleContext(sender *tb.User, input string) {
	convo := &Convo{}
	res := gl.Orm.
		Where("user_id = ?", sender.ID).
		Last(&convo)

	if res.Error != nil {
		gl.Bot.Send(sender, "Sorry we are out of context! /help")
		return
	}

	convo.handlers = make(map[string]Expector)
	convo.Handle("account", convo.expectAccount)
	convo.Handle("account choose", convo.expectAccount)
	convo.Handle("account que", convo.expectAccountQue)
	convo.Handle("account name", convo.expectAccountName)
	convo.Handle("amount", convo.expectAmount)
	convo.Handle("currency", convo.expectCurrency)
	convo.Handle("description", convo.expectDescription)
	convo.Handle("date", convo.expectDate)
	convo.Handle("from date", convo.expectFromDate)
	convo.Handle("till date", convo.expectTillDate)
	convo.expectNext(gl.Orm, input)

	gl.Bot.Send(sender, convo.response, &convo.menu)
}
