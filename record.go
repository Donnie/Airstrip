package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/Donnie/Airstrip/ptr"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleRecord(m *tb.Message) {
	userID := int64(m.Sender.ID)

	// end last conversation
	st.Orm.Unscoped().Where("user_id = ?", userID).Delete(&Convo{})

	// find what is this conversation about
	r, _ := regexp.Compile(`\/(\w+)`)
	form := r.FindString(m.Text)[1:]
	typ := "variable"
	if form == "income" || form == "charge" {
		typ = "fixed"
	}

	// Init empty record
	item := &Record{
		Form:   &form,
		UserID: &userID,
		Type:   &typ,
	}
	st.Orm.Create(&item)

	// Create new conversation with Context
	st.Orm.Create(&Convo{
		ContextID: &item.ID,
		Expect:    ptr.String("account"),
		UserID:    &userID,
	})

	question := genQues("account")
	st.Bot.Send(m.Sender, question)
}

func (st *State) handleDelete(m *tb.Message) {
	// end last conversation
	st.Orm.Unscoped().Where("user_id = ?", m.Sender.ID).Delete(&Convo{})

	// find if record ID is provided
	recordID, err := strconv.ParseInt(m.Payload, 10, 64)
	if err != nil {
		records := []Record{}
		st.Orm.Preload("Account").Limit(3).Order("id desc").Where("user_id = ?", m.Sender.ID).Find(&records)
		output := "You can choose from last three records:\n"
		for _, rec := range records {
			output += fmt.Sprintf("`ID: %d\t%s: %d %s`\n", rec.ID, *rec.Account.Name, *rec.Amount/100, *rec.Currency)
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
