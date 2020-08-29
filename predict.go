package main

import (
	"fmt"
	"regexp"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handlePredict(m *tb.Message) {
	var layout = "Jan 2006"
	userID := int64(m.Sender.ID)

	date := m.Payload
	matched, _ := regexp.Match(`^(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) \d{4}`, []byte(date))
	if !matched {
		date = "Dec 2030"
	}
	t, _ := time.Parse(layout, date)

	recs := []Record{}
	gl.Orm.
		Where("user_id = ?", userID).
		Find(&recs)

	output := fmt.Sprintf("%d", predictFuture(t, recs, false)/100)
	gl.Bot.Send(m.Sender, output)
}

func predictFuture(future time.Time, recs []Record, start bool) (cash int64) {
	for _, trans := range recs {
		if *trans.Form == "gain" {
			cash += *trans.Amount
		}
	}

	reps := monthDiff(time.Now(), future)
	if start {
		reps--
	}
	carry := calcMonthEnd(recs)

	cash += (int64(reps) * carry)
	return
}

func calcMonthEnd(recs []Record) (cash int64) {
	for _, trans := range recs {
		if *trans.Form == "income" {
			cash += *trans.Amount
		}
		if *trans.Form == "charge" {
			cash -= *trans.Amount
		}
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

	month += (year * 12)

	return
}
