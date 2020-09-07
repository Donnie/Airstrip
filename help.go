package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleHelp(m *tb.Message) {
	output := fmt.Sprintf("Hello %s!\n\n", m.Sender.FirstName)
	output += "I understand these commands:\n"
	output += "/expense Record an expense\n/gain Record any receipt\n/charge Record fixed costs like rent, etc.\n/income Record an income source like Salary\n/lend Lend money to someone\n/loan Take a loan from someone\n/predict Jan 2025 - Get a prediction of your financial standing\n/view Jan 2025 - Get a list of records pertaining to the month\n"

	st.Bot.Send(m.Sender, output)
}
