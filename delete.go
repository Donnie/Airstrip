package main

import (
	"fmt"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleDelete(m *tb.Message) {
	// end last conversation
	st.Orm.Unscoped().Where("user_id = ?", m.Sender.ID).Delete(&Convo{})

	// find if record ID is provided
	recordID, err := strconv.ParseInt(m.Payload, 10, 64)
	if err != nil {
		records := []Record{}
		st.Orm.Preload("AccountIn").
			Limit(3).Order("id desc").
			Where("user_id = ?", m.Sender.ID).
			Find(&records)

		output := "You can choose from last three records:\n"
		for _, rec := range records {
			output += fmt.Sprintf("`ID: %d\t%s: %d %s`\n", rec.ID, *rec.AccountIn.Name, *rec.Amount/100, *rec.AccountIn.Currency)
		}
		output += "\nReply with the ID for e.g.: `/delete 24`"

		st.Bot.Send(m.Sender, output, tb.ModeMarkdown)
		return
	}

	st.Orm.Where(recordID).
		Where("user_id = ?", m.Sender.ID).
		Delete(&Record{})
	st.Bot.Send(m.Sender, "Record Deleted.")
}
