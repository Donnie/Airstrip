package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

func (st *State) handleStand(m *tb.Message) {
	payload := strings.Fields(m.Payload)
	var acc string
	var mon time.Time

	if len(payload) > 0 {
		acc = payload[0]
	}

	if len(payload) > 1 {
		mon, _ = time.Parse(layout, payload[1]+payload[2])
	}

	var account Account
	err := st.Orm.Where("name = ?", acc).
		Where("user_id = ?", m.Sender.ID).
		First(&account).Error

	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		st.Bot.Send(m.Sender, "No account found by that name.")
		return
	}

	amount := getStand(st.Orm, account.ID, mon)
	output := fmt.Sprintf("*Standing*: `%.2f EUR`", amount)
	st.Bot.Send(m.Sender, output, tb.ModeMarkdown)
}

func getStand(db *gorm.DB, acc uint, mon time.Time) float64 {
	var res struct {
		Sum float64
	}

	query := fmt.Sprintf(`SELECT SUM(
		CASE 
		WHEN account_in_id = %d THEN amount * 1 
		WHEN account_out_id = %d THEN amount * -1 
		END
	)/100 as sum
	FROM records
	WHERE mandate = false`, acc, acc)

	if !mon.IsZero() {
		query += fmt.Sprintf(` AND EXTRACT(MONTH FROM date) = %d`, int(mon.Month()))
	}

	db.Raw(query).Scan(&res)
	return res.Sum
}
