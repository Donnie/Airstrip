package main

import (
	"fmt"
	"time"

	"github.com/Donnie/Airstrip/ptr"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleSavings(m *tb.Message) {
	past, err := time.Parse(monthFormat, m.Payload)
	if err != nil {
		past = time.Now().AddDate(-1, 0, 0)
	}

	savings := st.Analyse(past, m.Sender.ID)
	st.Bot.Send(m.Sender, savings, tb.ModeHTML)
}

// Analyse generates savings analysis for userID for the past
func (st *State) Analyse(past time.Time, userID int64) (out string) {
	savings := st.PastSavings(userID, past, nil)
	cash := st.CashTillNow(userID)
	out += fmt.Sprintf("Assets: €%.2f\n\n", float64(cash)/100)

	for _, save := range savings {
		// show only last twelve months
		// because of telegram message size limitation
		dateTime := parseDate(save.Month)
		out += fmt.Sprintf(
			"<strong>%s:</strong> <code>€%.2f</code>\n",
			dateTime.Format(monthFormat), (float64(save.Effect) / 100),
		)
	}
	return
}

// CashTillNow calculates Summation of all assets till now
func (st *State) CashTillNow(userID int64) int {
	var res struct{ Sum int }

	st.Orm.Raw(`SELECT SUM(
		CASE 
		WHEN ao.self = 1 AND ai.self = 0 THEN amount * -1 
		WHEN ai.self = 1 AND ao.self = 0 THEN amount * 1 
		END
	) as sum
	FROM records AS r 
	JOIN accounts AS ai ON r.account_in_id = ai.id 
	JOIN accounts AS ao ON r.account_out_id = ao.id 
	WHERE r.mandate = 0 
	AND r.deleted_at IS NULL
	AND r.user_id = ?`, userID).Scan(&res)

	return res.Sum
}

// PastSavings calculates savings till a past date
func (st *State) PastSavings(userID int64, start time.Time, end *time.Time) (savings []Saving) {
	if end == nil {
		end = ptr.Time(time.Now())
	}

	query := fmt.Sprintf(`
		WITH RECURSIVE MonthDates(monthDate) AS (
			SELECT DATETIME('%s 00:00:00')
			UNION ALL
			SELECT DATETIME(monthDate, '+1 month')
			FROM MonthDates
			WHERE DATETIME(monthDate, '+1 month') < DATETIME('%s 23:59:59')
		)
		SELECT
			startDate as month,
			COALESCE((
				SELECT
					SUM(
						CASE
							WHEN aci.self = 1 AND aco.self = 0 THEN amount
							WHEN aci.self = 0 AND aco.self = 1 THEN amount * -1
						END
					) as delta
				FROM records AS r
				LEFT JOIN accounts AS aci ON r.account_in_id = aci.id
				LEFT JOIN accounts AS aco ON r.account_out_id = aco.id
				WHERE
					r.deleted_at IS NULL
					AND r.date BETWEEN months.startDate AND months.endDate
					AND r.user_id = %d
			), 0) as effect
		FROM (
			SELECT DATETIME(monthDate, 'start of month') as startDate,
				DATETIME(monthDate, '+1 month', 'start of month') as endDate
			FROM MonthDates
		) AS months
	`, start.Format("2006-01-02"), end.Format("2006-01-02"), userID)

	st.Orm.Raw(query).Scan(&savings)
	return
}
