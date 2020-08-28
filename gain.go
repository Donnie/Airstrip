package main

import (
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleGain(m *tb.Message) {
	layout := "2006-01-02 15:04"
	form := "gain"
	var name, amount, currency, description, date string
	payload := strings.Split(m.Payload, "  ")
	switch len(payload) {
	case 2:
		name = payload[0]
		amount = payload[1]
	case 3:
		name = payload[0]
		amount = payload[1]
		currency = payload[2]
	case 4:
		name = payload[0]
		amount = payload[1]
		currency = payload[2]
		description = payload[3]
	case 5:
		name = payload[0]
		amount = payload[1]
		currency = payload[2]
		description = payload[3]
		date = payload[4]
	}

	if name == "" || amount == "" {
		panic("Name or amount is empty")
	}

	amountFlt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic("amount could not be parsed")
	}
	amountInt := int64(amountFlt * 100)

	var dateTime time.Time
	if date != "" {
		dateTime, err = time.Parse(layout, date)
		if err != nil {
			panic("date could not be parsed")
		}
	} else {
		dateTime = time.Now()
	}

	if currency == "" {
		currency = "EUR"
	}

	item := &Variable{
		Amount:      &amountInt,
		Currency:    &currency,
		Date:        &dateTime,
		Description: &description,
		Form:        &form,
		Name:        &name,
	}

	gl.Orm.Create(item)

	output := "Gain added."
	gl.Bot.Send(m.Sender, output)
}
