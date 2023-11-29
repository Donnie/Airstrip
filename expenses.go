package main

import (
	"fmt"
	"time"

	"github.com/Donnie/Airstrip/ptr"
	tb "gopkg.in/tucnak/telebot.v2"
)

// handleExpenses handles the expenses page
func (st *State) handleExpenses(m *tb.Message) {
	// get date from input
	past, err := time.Parse(monthFormat, m.Payload)

	// if date is not provided, use current month
	if err != nil {
		past = time.Now().AddDate(-1, 0, 0)
	}

	expenses := st.ExpensesAnalyse(past, m.Sender.ID)
	st.Bot.Send(m.Sender, expenses, tb.ModeHTML)
}

// ExpensesAnalyse returns the expenses till the given month
func (st *State) ExpensesAnalyse(past time.Time, userID int64) (out string) {
	expenses := st.PastExpenses(userID, past, nil)

	for _, save := range expenses {
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

// PastExpenses returns the expenses between the given months
func (st *State) PastExpenses(userID int64, start time.Time, end *time.Time) (savings []Saving) {
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
							WHEN aci.self = 0 AND aco.self = 1 THEN amount
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
			SELECT DATETIME(monthDate, 'start of month', '-1 second') as startDate,
				DATETIME(monthDate, 'start of month', '+1 month', '-1 second') as endDate
			FROM MonthDates
		) AS months
	`, start.Format("2006-01-02"), end.Format("2006-01-02"), userID)

	st.Orm.Raw(query).Scan(&savings)
	return
}
