package main

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var layout = "Jan 2006"

func (gl *Global) handlePredict(m *tb.Message) {
	dat, err := ioutil.ReadFile(gl.File)
	check(err)
	account := &Account{}
	json.Unmarshal(dat, account)

	date := m.Payload
	matched, _ := regexp.Match(`^(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) \d{4}`, []byte(date))
	if !matched {
		date = "Dec 2030"
	}
	t, _ := time.Parse(layout, date)

	output := strconv.Itoa(int(predictFuture(t, *account, false)))
	gl.Bot.Send(m.Sender, output)
}

func predictFuture(future time.Time, account Account, start bool) (cash int32) {
	for _, trans := range account.Savings {
		cash = cash + trans.Amount
	}

	reps := monthDiff(time.Now(), future)
	if start {
		reps--
	}
	carry := calcMonthEnd(account)

	cash = cash + (int32(reps) * carry)
	return
}

func calcMonthEnd(account Account) (cash int32) {
	for _, trans := range account.Earnings {
		cash = cash + trans.Amount
	}
	for _, trans := range account.Costs {
		cash = cash - trans.Amount
	}

	return
}

func monthDiff(a, b time.Time) (month int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, m1, _ := a.Date()
	y2, m2, _ := b.Date()

	month = int(m2 - m1)
	year := int(y2 - y1)

	month = month + (year * 12)

	return
}
