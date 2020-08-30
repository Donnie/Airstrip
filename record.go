package main

import (
	"encoding/json"
	"regexp"

	"github.com/Donnie/Airstrip/ptr"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleRecord(m *tb.Message) {
	userID := int64(m.Sender.ID)

	// end last conversation
	gl.Orm.Where("user_id = ?", userID).Delete(&Convo{})

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
	gl.Orm.Create(&item)

	// Create new conversation with Context
	context, _ := json.Marshal(item)
	gl.Orm.Create(&Convo{
		Context: ptr.String(string(context)),
		Expect:  ptr.String("account"),
		UserID:  &userID,
	})

	question := genQues("account", form)
	gl.Bot.Send(m.Sender, question)
}
