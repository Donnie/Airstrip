package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleCancel(m *tb.Message) {
	var convo Convo
	// get last conversation
	st.Orm.Where("user_id = ?", m.Sender.ID).First(&convo)
	// end last conversation
	st.Orm.Unscoped().Where("user_id = ?", m.Sender.ID).Delete(&Convo{})

	if convo.ContextID != nil {
		// remove related record
		st.Orm.Unscoped().
			Where(*convo.ContextID).
			Delete(&Record{})
	}

	st.Bot.Send(m.Sender, "Canceled.")
}
