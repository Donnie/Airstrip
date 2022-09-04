package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) startBot() {
	st.Bot.Handle("/start", st.handleHelp)
	st.Bot.Handle("/cancel", st.handleCancel)
	st.Bot.Handle("/delete", st.handleDelete)
	st.Bot.Handle("/help", st.handleHelp)
	st.Bot.Handle("/predict", st.handlePredict)
	st.Bot.Handle("/record", st.handleRecord)
	st.Bot.Handle("/recur", st.handleRecord)
	st.Bot.Handle("/savings", st.handleSavings)
	st.Bot.Handle("/stand", st.handleStand)
	st.Bot.Handle("/view", st.handleView)

	st.Bot.Handle(tb.OnText, st.handleText)
	st.Bot.Handle(tb.OnCallback, st.handleCallback)
	st.Bot.Start()
}
