package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Donnie/Airstrip/ptr"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleContext(m *tb.Message) {
	userID := int64(m.Sender.ID)
	var ask string

	convo := &Convo{}
	res := gl.Orm.
		Where("user_id = ?", userID).
		Last(convo)

	if res.Error != nil {
		gl.Bot.Send(m.Sender, res.Error)
		return
	}

	context := &Record{}
	json.Unmarshal([]byte(*convo.Context), &context)

	if context.Account == nil {
		context.Account = &m.Text
		ask = "amount"
		gl.Orm.Save(&context)
		cont, _ := json.Marshal(context)
		convo.Context = ptr.String(string(cont))
		gl.Orm.Save(&convo)
	} else if context.Amount == nil {
		amountFlt, err := strconv.ParseFloat(m.Text, 64)
		if err != nil {
			gl.Bot.Send(m.Sender, "amount could not be parsed")
			return
		}
		context.Amount = ptr.Int64(int64(amountFlt * 100))
		ask = "currency"
		gl.Orm.Save(&context)
		cont, _ := json.Marshal(context)
		convo.Context = ptr.String(string(cont))
		gl.Orm.Save(&convo)
	} else if context.Currency == nil {
		currency := m.Text
		if currency == "" || currency == "." {
			currency = "EUR"
		}
		context.Currency = &currency
		ask = "description"
		gl.Orm.Save(&context)
		cont, _ := json.Marshal(context)
		convo.Context = ptr.String(string(cont))
		gl.Orm.Save(&convo)
	} else if context.Description == nil {
		context.Description = &m.Text
		if *context.Type == "variable" {
			ask = "date"
		} else {
			ask = "from date"
		}
		gl.Orm.Save(&context)
		cont, _ := json.Marshal(context)
		convo.Context = ptr.String(string(cont))
		gl.Orm.Save(&convo)
	} else if context.Date == nil && *context.Type == "variable" {
		date := m.Text
		layout := "2006-01-02 15:04"
		dateTime, err := time.Parse(layout, date)
		if err != nil {
			dateTime = time.Now()
		}
		context.Date = &dateTime
		gl.Orm.Save(&context)
		gl.Orm.Delete(&convo)
	} else if context.FromDate == nil && *context.Type == "fixed" {
		date := m.Text
		layout := "Jan 2006"
		dateTime, _ := time.Parse(layout, date)
		context.FromDate = &dateTime
		ask = "till date"
		gl.Orm.Save(&context)
		cont, _ := json.Marshal(context)
		convo.Context = ptr.String(string(cont))
		gl.Orm.Save(&convo)
	} else if context.TillDate == nil && *context.Type == "fixed" {
		date := m.Text
		layout := "Jan 2006"
		dateTime, err := time.Parse(layout, date)
		if err == nil {
			dateTime = dateTime.AddDate(0, 1, -1)
			context.TillDate = &dateTime
			gl.Orm.Save(&context)
		}
		if *context.Form == "lend" {
			context.Form = ptr.String("expense")
			gl.Orm.Create(&context)
		} else if *context.Form == "loan" {
			context.Form = ptr.String("gain")
			gl.Orm.Create(&context)
		}
		gl.Orm.Delete(&convo)
	}

	if ask != "" {
		question := genQues(ask, *context.Form)
		gl.Bot.Send(m.Sender, question)
	} else {
		gl.Bot.Send(m.Sender, "Record stored!")
	}
}

func genQues(ask, form string) (out string) {
	out = fmt.Sprintf("What is the %s of the %s?", ask, form)
	return
}
