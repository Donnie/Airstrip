package main

import (
	"github.com/Donnie/Airstrip/ptr"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleRecord(m *tb.Message) {
	userID := int64(m.Sender.ID)

	// end last conversation
	st.Orm.Unscoped().Where("user_id = ?", userID).Delete(&Convo{})

	// Init empty record
	item := &Record{
		UserID: &userID,
	}
	st.Orm.Create(&item)

	// Create new conversation with Context
	convo := &Convo{
		ContextID: &item.ID,
		Expect:    ptr.String("form"),
		UserID:    &userID,
	}
	st.Orm.Create(&convo)

	convo.response = genQues("form")
	convo.menu = tb.ReplyMarkup{}
	convo.menu.Inline(
		convo.menu.Row(
			convo.menu.Data("Expense", "expense"),
			convo.menu.Data("Charge", "charge"),
			convo.menu.Data("Loan", "loan"),
		),
		convo.menu.Row(
			convo.menu.Data("Gain", "gain"),
			convo.menu.Data("Income", "income"),
			convo.menu.Data("Lend", "lend"),
		),
	)
	st.Bot.Send(m.Sender, convo.response, &convo.menu)
}
