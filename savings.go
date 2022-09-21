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
	st.Bot.Send(m.Sender, savings, tb.ModeHTML)
}

// Analyse generates savings analyse for userID for the past
func (st *State) Analyse(past time.Time, userID int64) (out string) {
	savings := st.PastSavings(userID, past)
	cash := st.CashTillNow(userID)
	out += fmt.Sprintf("Assets: %.2f EUR\n\n", float64(cash)/100)

	for _, save := range savings {
		// show only last twelve months
		// because of telegram message size limitation
		dateTime := parseDate(save.Month)
		out += fmt.Sprintf(
			"<strong>%s:</strong> <code>%.2f EUR</code>\n",
			dateTime.Format(monthFormat), (float64(save.Effect) / 100),
		)
	}
	return
}

// PastSavings calculates savings till a past date
func (st *State) PastSavings(userID int64, past time.Time) (savings []Saving) {
	query := fmt.Sprintf(`
		WITH RECURSIVE MonthDates(monthDate) AS (
			SELECT '%s'
			UNION ALL
			SELECT date(monthDate, '+1 month')
			FROM MonthDates
			WHERE date(monthDate, '+1 month') < CURRENT_TIMESTAMP
		)
		SELECT
			endDate as month,
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
					AND (aci.self = 0 OR aco.self = 0)
					AND r.date IS NOT NULL
					AND r.date BETWEEN months.startDate AND months.endDate
					AND r.user_id = %d
			), 0) as effect
		FROM (
			SELECT date(monthDate, 'start of month') as startDate,
				date(monthDate, '+1 month', 'start of month', '-1 day') as endDate
			FROM MonthDates
		) AS months
	`, past.Format("2006-01-02"), userID)

	st.Orm.Raw(query).Scan(&savings)
	return
}
