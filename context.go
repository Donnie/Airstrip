package main

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleContext(m *tb.Message) {
	userID := int64(m.Sender.ID)

	convo := &Convo{}
	res := gl.Orm.
		Where("user_id = ?", userID).
		Last(convo)

	if res.Error != nil {
		gl.Bot.Send(m.Sender, "Sorry we are out of context! /help")
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

	gl.Bot.Send(m.Sender, convo.expectNext(gl.Orm, m.Text), &convo.menu)
}

func (gl *Global) handleCallback(m *tb.Callback) {
	userID := int64(m.Sender.ID)
	data := strings.TrimSpace(m.Data)

	convo := &Convo{}
	res := gl.Orm.
		Where("user_id = ?", userID).
		Last(convo)

	if res.Error != nil {
		gl.Bot.Send(m.Sender, "Sorry we are out of context! /help")
		return
	}

	gl.Bot.Respond(m, &tb.CallbackResponse{
		CallbackID: m.ID,
		Text:       "Cool!",
	})

	convo.handlers = make(map[string]Expector)
	convo.Handle("account que", convo.expectAccountQue)
	gl.Bot.Send(m.Sender, convo.expectNext(gl.Orm, data), &convo.menu)
}
