package main

import (
	"github.com/Donnie/Airstrip/ptr"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleRecord(m *tb.Message) {
	userID := m.Sender.ID

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
		Expect:    ptr.String("account in"),
		UserID:    &userID,
	}
	st.Orm.Create(&convo)

	convo.response = genQues("account in")
	st.Bot.Send(m.Sender, convo.response, &convo.menu)
}
