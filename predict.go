package main

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handlePredict(m *tb.Message) {
	var layout = "Jan 2006"
	t, err := time.Parse(layout, m.Payload)
	if err != nil {
		t = time.Now().AddDate(1, 0, 0)
	}

	cashCurr := st.CashTillNow(m.Sender.ID)
	costPlan := st.PlannedExpensesCurrentMonth(m.Sender.ID)
	cashFutr := st.FutureSavings(m.Sender.ID, t)

	output := fmt.Sprintf("*Prediction* for month end %s:\n%d EUR", t.Format(layout), (cashCurr-costPlan+cashFutr)/100)
	st.Bot.Send(m.Sender, output, tb.ModeMarkdown)
}

// CashTillNow calculates Summation of all assets till now
func (st *State) CashTillNow(userID int) int {
	var res struct {
		Sum int
	}

	st.Orm.Raw(`SELECT SUM(
		CASE 
		WHEN ao.self AND ai.self = false THEN amount * -1 
		WHEN ai.self AND ao.self = false THEN amount * 1 
		END
	)
	FROM records AS r 
	JOIN accounts AS ai ON r.account_in_id = ai.id 
	JOIN accounts AS ao ON r.account_out_id = ao.id 
	WHERE r.mandate = false 
	AND r.user_id = ?`, userID).Scan(&res)

	return res.Sum
}

// PlannedExpensesCurrentMonth calculates remaining expenses for current month
func (st *State) PlannedExpensesCurrentMonth(userID int) int {
	var res struct {
		Sum int
	}
	t := time.Now()

	st.Orm.Raw(`SELECT SUM(charge - expense) AS sum FROM (
		SELECT ac.name,
		COALESCE((
			SELECT SUM(r1.amount) 
			FROM records AS r1
			WHERE (
				(
					?::date BETWEEN r1.from_date AND r1.till_date 
					OR
					?::date >= r1.from_date AND r1.till_date IS NULL
				)
				AND r1.account_in_id = ac.id
				AND r1.mandate = true
			)
		), 0) AS charge,
		COALESCE((
			SELECT SUM(r1.amount) 
			FROM records AS r1
			WHERE (
				(
					EXTRACT(MONTH FROM date) = ? 
					AND
					EXTRACT(YEAR FROM date) = ?
				)
				AND r1.account_in_id = ac.id
				AND r1.mandate = false
			)
		), 0) AS expense
		FROM accounts AS ac
		WHERE ac.user_id = ? 
		AND ac.self = false
	) AS planned
	WHERE charge-expense > 0`, t, t, int(t.Month()), t.Year(), userID).Scan(&res)

	return res.Sum
}

// FutureSavings calculates savings till a future date
func (st *State) FutureSavings(userID int, fut time.Time) int {
	var res struct {
		Sum int
	}

	st.Orm.Raw(`SELECT SUM(income - charge) AS sum
	FROM (
		SELECT future.month, COALESCE((
			SELECT SUM(r1.amount) 
			FROM records AS r1
			JOIN accounts AS ai1 ON r1.account_in_id = ai1.id 
			JOIN accounts AS ao1 ON r1.account_out_id = ao1.id 
			WHERE (
				(
					future.month::date BETWEEN r1.from_date AND r1.till_date 
					OR
					future.month::date >= r1.from_date AND r1.till_date IS NULL
				)
				AND r1.mandate
				AND ao1.self = false
				AND ai1.self
				AND r1.user_id = ?
			)
		), 0) AS income, 
		COALESCE((
			SELECT SUM(r1.amount) 
			FROM records AS r1
			JOIN accounts AS ai1 ON r1.account_in_id = ai1.id 
			JOIN accounts AS ao1 ON r1.account_out_id = ao1.id 
			WHERE (
				(
					future.month::date BETWEEN r1.from_date AND r1.till_date 
					OR
					future.month::date >= r1.from_date AND r1.till_date IS NULL
				)
				AND r1.mandate
				AND ao1.self
				AND ai1.self = false
				AND r1.user_id = ?
			)
		), 0) AS charge
		FROM (
			SELECT date_trunc('month', current_date)
			+ (INTERVAL '1 month' * generate_series(1, months::int)) AS month 
			FROM (
				SELECT EXTRACT(year FROM diff) * 12 + EXTRACT(month FROM diff) AS months 
				FROM (
					SELECT age(
						date_trunc('month', CAST(? AS TIMESTAMP)) + INTERVAL '1 month - 1 day',
						current_timestamp
					) AS diff
				) AS fut
			) AS reps
		) AS future
	) AS savings`, userID, userID, fut).Scan(&res)

	return res.Sum
}
