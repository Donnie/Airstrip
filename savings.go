package main

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleSavings(m *tb.Message) {
	past, err := time.Parse(monthFormat, m.Payload)
	if err != nil {
		past = time.Now().AddDate(-1, 0, 0)
	}

	savings := st.Analyse(past, m.Sender.ID)
	st.Bot.Send(m.Sender, savings, tb.ModeMarkdownV2)
}

// Analyse generates savings analyse for userID for the past
func (st *State) Analyse(past time.Time, userID int) (out string) {
	savings := st.PastSavings(userID, past)
	cash := st.CashTillNow(userID)
	out += fmt.Sprintf("Assets: %d EUR\n\n", cash/100)

	for _, save := range savings {
		// show only last twelve months
		// because of telegram message size limitation
		out += fmt.Sprintf(
			"*%s:* `%.2f EUR`\n",
			save.Month.Format(monthFormat), (float64(save.Effect) / 100),
		)
	}
	return
}

// PastSavings calculates savings till a past date
func (st *State) PastSavings(userID int, past time.Time) (savings []Saving) {
	st.Orm.Raw(`SELECT
		month,
		COALESCE((
			SELECT
				SUM(
					CASE
						WHEN aci.self AND aco.self = false THEN amount
						WHEN aci.self = false AND aco.self THEN amount * -1
					END
				) as delta
			FROM records AS r
			LEFT JOIN accounts AS aci ON r.account_in_id = aci.id
			LEFT JOIN accounts AS aco ON r.account_out_id = aco.id
			WHERE 
				r.deleted_at IS NULL
				AND (aci.self = false OR aco.self = false)
				AND r.date IS NOT NULL
				AND DATE_TRUNC('month',r.date) = month
				AND r.user_id = ?
		), 0) as effect
	FROM (
		SELECT
			DATE_TRUNC('month', current_date) - (INTERVAL '1 month' * GENERATE_SERIES(0, months::int)) AS month
		FROM 
			(SELECT EXTRACT(year FROM diff) * 12 + EXTRACT(month FROM diff) AS months
			FROM 
				(SELECT age(
						DATE_TRUNC('month', CURRENT_TIMESTAMP),
						DATE_TRUNC('month', CAST(? AS TIMESTAMP)) 
				) AS diff
			) AS fut
		) AS reps
	) AS months`, userID, past).Scan(&savings)
	return
}
