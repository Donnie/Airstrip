package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleHelp(m *tb.Message) {
	output := fmt.Sprintf("Hello %s!\n\n", m.Sender.FirstName)
	output += "I understand these commands:\n" +
		"\n/record - Record an expense or gain" +
		"\n/recur - Record an income or a charge" +
		"\n/cancel - Cancel an ongoing record or recur process" +
		"\n/delete - Delete any record or recur" +
		"\n/predict - Get a prediction of your financial standing for next 12 months" +
		"\n/predict Jan 2025 - Get a prediction of your financial standing" +
		"\n/view - Get a list of records pertaining to the current month" +
		"\n/view Jan 2025 - Get a list of records pertaining to the month" +
		"\n/savings - Get savings per month from last 12 months" +
		"\n/savings Jan 2025 - Get savings per month beginning from a certain month" +
		"\n/stand - Get a current status of all accounts" +
		"\n/stand Account - Get a current standing of any account" +
		"\n/stand Account Mar 2022 - Get a month wise total effect on any account" +
		"\n/help - To see this list again"

	st.Bot.Send(m.Sender, output)
}
