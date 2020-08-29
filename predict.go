package main

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handlePredict(m *tb.Message) {
	var layout = "Jan 2006"
	userID := int64(m.Sender.ID)

	date := m.Payload
	t, err := time.Parse(layout, date)
	if err != nil {
		t = time.Now().AddDate(1, 0, 0)
	} else {
		t = t.AddDate(0, 1, -1)
	}

	recs := []Record{}
	gl.Orm.
		Where("user_id = ?", userID).
		Find(&recs)

	output := fmt.Sprintf("%d", predictFuture(t, recs)/100)
	gl.Bot.Send(m.Sender, output)
}

func predictFuture(future time.Time, recs []Record) (cash int64) {
	for _, rec := range recs {
		switch *rec.Form {
		case "gain", "lend":
			cash += *rec.Amount
		case "expense", "loan":
			cash -= *rec.Amount
		}
	}

	timeStr := fmt.Sprintf("%d%d", time.Now().Year(), time.Now().Month())
	timeNow, _ := time.Parse("20061", timeStr)
	timeNow = timeNow.AddDate(0, 1, 0)

	reps := monthDiff(timeNow, future)
	for i := 0; i < reps; i++ {
		cash += calcMonthEnd(recs, timeNow.AddDate(0, i, 0))
	}
	return
}

func calcMonthEnd(recs []Record, month time.Time) (cash int64) {
	for _, rec := range recs {
		if month.Unix() >= rec.FromDate.Unix() && month.Unix() <= rec.TillDate.Unix() {
			switch *rec.Form {
			case "income":
				cash += *rec.Amount
			case "charge":
				cash -= *rec.Amount
			}
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

	month = int(m2 - m1 + 1)
	year := int(y2 - y1)

	month += (year * 12)

	return
}
