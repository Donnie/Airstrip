package main

import (
	"encoding/json"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleShow(m *tb.Message) {
	var layout = "Jan 2006"
	t, err := time.Parse(layout, m.Payload)
	if err != nil {
		t = time.Now()
	}

	recs := []Record{}
	gl.Orm.
		Where(
			gl.Orm.Where("(?::date BETWEEN from_date AND till_date)", time.Now()).
				Or("(?::date >= from_date AND till_date IS NULL)", time.Now()),
		).
		Or(
			gl.Orm.Where("EXTRACT(MONTH FROM date) = ?", int(t.Month())).
				Where("EXTRACT(YEAR FROM date) = ?", t.Year()),
		).
		Find(&recs)

	r, _ := json.Marshal(recs)
	gl.Bot.Send(m.Sender, "`"+string(r)+"`", tb.ModeMarkdown)
}
