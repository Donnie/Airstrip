package main

import (
	"encoding/json"
	"regexp"

	"github.com/Donnie/Airstrip/ptr"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleRecord(m *tb.Message) {
	var typ string
	userID := int64(m.Sender.ID)
	r, _ := regexp.Compile(`\/(\w+)`)

	form := r.FindString(m.Text)[1:]
	if form == "income" || form == "charge" {
		typ = "fixed"
	} else {
		typ = "variable"
	}

	question := genQues("account", form)

	item := &Record{
		Form:   &form,
		UserID: &userID,
		Type:   &typ,
	}
	gl.Orm.Create(&item)

	gl.Orm.Where("user_id = ?", userID).Delete(&Convo{})

	cont, _ := json.Marshal(item)
	gl.Orm.Create(&Convo{
		Context: ptr.String(string(cont)),
		UserID:  &userID,
	})
	gl.Bot.Send(m.Sender, question)
}
