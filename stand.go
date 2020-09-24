package main

import (
	"errors"
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

func (st *State) handleStand(m *tb.Message) {
	// find out account
	var account Account
	err := st.Orm.Where("name = ?", m.Payload).
		Where("user_id = ?", m.Sender.ID).
		First(&account).Error

	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		st.Bot.Send(m.Sender, "No account found by that name.")
		return
	}

	amount := getStand(st.Orm, account.ID)
	output := fmt.Sprintf("*Current Standing*: `%.2f EUR`", amount)
	st.Bot.Send(m.Sender, output, tb.ModeMarkdown)
}

func getStand(db *gorm.DB, acc uint) float64 {
	var res struct {
		Sum float64
	}
	db.Raw(`SELECT SUM(
		CASE 
		WHEN account_in_id = ? THEN amount * 1 
		WHEN account_out_id = ? THEN amount * -1 
		END
	)/100 as sum FROM records WHERE mandate = false`, acc, acc).Scan(&res)
	return res.Sum
}
