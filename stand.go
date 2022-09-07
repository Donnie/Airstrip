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
	var output string
	totalAmount := 0.0
	totalLiquid := 0.0

	if len(payload) > 0 {
		acc = payload[0]
	}

	if acc == "" {
		stands := getStandAll(st.Orm, m.Sender.ID)
		for _, stand := range stands {
			if stand.Stand == 0 {
				continue
			}
			totalAmount += stand.Stand
			if stand.Liquid {
				totalLiquid += stand.Stand
			}
			output += fmt.Sprintf("*%s*: `%.2f EUR`\n", stand.Name, stand.Stand)
		}
		output += fmt.Sprintf("\n*Liquid*: `%.2f EUR`", totalLiquid)
		output += fmt.Sprintf("\n*Total*: `%.2f EUR`", totalAmount)
	} else {
		if len(payload) > 1 {
			mon, _ = time.Parse(monthFormat, payload[1]+" "+payload[2])
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
		output = fmt.Sprintf("*Standing*: `%.2f EUR`", amount)
	}
	st.Bot.Send(m.Sender, output, tb.ModeMarkdown)
}

func getStand(db *gorm.DB, acc uint, mon time.Time) float64 {
	var res struct {
		Sum float64
	}

	query := fmt.Sprintf(`SELECT CAST(SUM(
		CASE 
		WHEN account_in_id = %d THEN amount * 1 
		WHEN account_out_id = %d THEN amount * -1 
		END
	) as REAL)/100 as sum
	FROM records
	WHERE mandate = 0
	AND deleted_at IS NULL`, acc, acc)

	if !mon.IsZero() {
		query += fmt.Sprintf(` AND LTRIM(STRFTIME('%%m', date), "0") = "%d"`, int(mon.Month()))
		query += fmt.Sprintf(` AND STRFTIME('%%Y', date) = "%d"`, int(mon.Year()))
	}

	db.Raw(query).Scan(&res)
	return res.Sum
}

func getStandAll(db *gorm.DB, userID int64) (res []Stand) {
	db.Raw(`SELECT name, liquid, CAST((total_in-total_out) as REAL)/100 AS stand
	FROM (
		SELECT a.name, a.liquid,
		(SELECT COALESCE(SUM(amount), 0) FROM records WHERE account_in_id = a.id AND mandate = 0 AND deleted_at IS NULL) AS total_in, 
		(SELECT COALESCE(SUM(amount), 0) FROM records WHERE account_out_id = a.id AND mandate = 0 AND deleted_at IS NULL) AS total_out
		FROM accounts AS a
		WHERE a.self = 1
		AND a.user_id = ?
	) AS total
	ORDER BY name asc`, userID).Scan(&res)
	return
}
