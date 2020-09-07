package main

import (
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleText(m *tb.Message) {
	st.handleContext(m.Sender, m.Text)
}

func (st *State) handleCallback(m *tb.Callback) {
	st.handleContext(m.Sender, strings.TrimSpace(m.Data))
	st.Bot.Respond(m, &tb.CallbackResponse{
		CallbackID: m.ID,
		Text:       "Cool!",
	})
}

func (st *State) handleContext(sender *tb.User, input string) {
	convo := &Convo{}
	res := st.Orm.
		Where("user_id = ?", sender.ID).
		Last(&convo)

	if res.Error != nil {
		st.Bot.Send(sender, "Sorry we are out of context! /help")
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
	convo.Handle("form", convo.expectForm)
	convo.Handle("from date", convo.expectFromDate)
	convo.Handle("till date", convo.expectTillDate)
	convo.expectNext(st.Orm, input)

	st.Bot.Send(sender, convo.response, &convo.menu)
}
