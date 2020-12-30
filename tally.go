package main

import (
	"fmt"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

func (st *State) handleTally(m *tb.Message) {
	payload := strings.Fields(m.Payload)
	var mon time.Time
	var output string

	if payload[0] == "" {
		output = "Please specify tag."
	} else {
		if len(payload) > 2 {
			mon, _ = time.Parse(monthFormat, payload[1]+" "+payload[2])
		}

		amount := getTally(st.Orm, payload[0], mon)
		output = fmt.Sprintf("*Standing*: `%.2f EUR`", amount)
	}
	st.Bot.Send(m.Sender, output, tb.ModeMarkdown)
}

func getTally(db *gorm.DB, tag string, mon time.Time) float64 {
	var res struct {
		Sum float64
	}

	query := fmt.Sprintf(`SELECT SUM(
		CASE 
		WHEN '%s' = any(ai.tags) THEN amount * 1 
		WHEN '%s' = any(ao.tags) THEN amount * -1 
		END
	)/100 as sum
	FROM records
	JOIN accounts AS ai ON ai.id = records.account_in_id
	JOIN accounts AS ao ON ao.id = records.account_out_id
	WHERE deleted_at = NULL
	AND mandate = false`, tag, tag)

	if !mon.IsZero() {
		query += fmt.Sprintf(` AND EXTRACT(MONTH FROM date) = %d`, int(mon.Month()))
	}

	db.Raw(query).Scan(&res)
	return res.Sum
}
