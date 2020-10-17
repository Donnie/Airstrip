package main

import (
	"strings"

	"github.com/Donnie/Airstrip/ptr"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleRecord(m *tb.Message) {
	userID := m.Sender.ID

	// end last conversation
	st.Orm.Unscoped().Where("user_id = ?", userID).Delete(&Convo{})

	// Init empty record
	item := &Record{
		UserID:  &userID,
		Mandate: ptr.Bool(strings.Contains(m.Text, "recur")),
	}
	st.Orm.Create(&item)

	// Create new conversation with Context
	convo := &Convo{
		ContextID: &item.ID,
		UserID:    &userID,
	}
	convo.askAmount()
	st.Orm.Create(&convo)

	convo.response = genQues("amount")
	st.Bot.Send(m.Sender, convo.response, &convo.menu)
}
