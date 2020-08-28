package main

import (
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleCharge(m *tb.Message) {
	var layout = "Jan 2006"
	form := "charge"
	var name, amount, currency, description, from, till string
	payload := strings.Split(m.Payload, ". ")
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
		from = payload[3]
	case 5:
		name = payload[0]
		amount = payload[1]
		currency = payload[2]
		from = payload[3]
		till = payload[4]
	}

	if name == "" || amount == "" {
		panic("Name or amount is empty")
	}

	amountFlt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic("amount could not be parsed")
	}
	amountInt := int64(amountFlt * 100)

	var fromDate time.Time
	if from != "" {
		fromDate, err = time.Parse(layout, from)
		if err != nil {
			panic("date could not be parsed")
		}
	} else {
		fromDate = time.Now()
	}

	var tillDate *time.Time
	if till != "" {
		tyme, err := time.Parse(layout, till)
		if err != nil {
			panic("till date could not be parsed")
		}
		tillDate = &tyme
	}

	if currency == "" {
		currency = "EUR"
	}

	item := &Fixed{
		Amount:      &amountInt,
		Currency:    &currency,
		Form:        &form,
		FromDate:    &fromDate,
		TillDate:    tillDate,
		Description: &description,
		Name:        &name,
	}

	gl.Orm.Create(item)

	output := "Charge added."
	gl.Bot.Send(m.Sender, output)
}
