package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

// handleCancel cancels the last conversation
func (st *State) handleCancel(m *tb.Message) {
	// Get the last conversation.
	var convo Convo
	st.Orm.Where("user_id = ?", m.Sender.ID).First(&convo)

	// End the last conversation.
	st.Orm.Unscoped().Where("user_id = ?", m.Sender.ID).Delete(&Convo{})

	// Remove the related record.
	if convo.ContextID != nil {
		st.Orm.Unscoped().
			Where(*convo.ContextID).
			Delete(&Record{})
	}

	st.Bot.Send(m.Sender, "Canceled.")
}
