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
	convo.Handle("amount", convo.expectAmount)
	convo.Handle("account in", convo.expectAccountIn)
	convo.Handle("account out", convo.expectAccountOut)
	convo.Handle("account new in", convo.expectAccountInQue)
	convo.Handle("account new out", convo.expectAccountOutQue)
	convo.Handle("account new self in", convo.expectCreateAccountSelfIn)
	convo.Handle("account new self out", convo.expectCreateAccountSelfOut)
	convo.Handle("description", convo.expectDescription)
	convo.Handle("date", convo.expectDate)
	convo.Handle("from date", convo.expectFromDate)
	convo.Handle("till date", convo.expectTillDate)
	convo.expectNext(st.Orm, input)

	st.Bot.Send(sender, convo.response, &convo.menu, tb.ModeHTML)
}
