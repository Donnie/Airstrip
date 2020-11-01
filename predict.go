package main

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handlePredict(m *tb.Message) {
	fut, err := time.Parse(monthFormat, m.Payload)
	if err != nil {
		fut = time.Now().AddDate(1, 0, 0)
	}

	predictions := st.Predict(fut, m.Sender.ID)
	st.Bot.Send(m.Sender, predictions, tb.ModeMarkdownV2)
}

// CashTillNow calculates Summation of all assets till now
func (st *State) CashTillNow(userID int) int {
	var res struct{ Sum int }

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

// PlannedCurrentMonth calculates remaining costs and incomes for current month
func (st *State) PlannedCurrentMonth(userID int, cost bool) int {
	var res struct{ Sum int }
	dir := "out"
	if cost {
		dir = "in"
	}

	query := fmt.Sprintf(
		`SELECT COALESCE(SUM(target - current), 0) AS sum 
		FROM (SELECT ac.name, 
				SUM(CASE WHEN r.mandate THEN r.amount ELSE 0 END) AS target, 
				SUM(CASE WHEN r.mandate THEN 0 ELSE r.amount END) AS current
			FROM accounts AS ac 
			JOIN records AS r 
				ON r.account_%s_id = ac.id 
				AND (
					(
						CURRENT_TIMESTAMP BETWEEN r.from_date AND r.till_date
						OR CURRENT_TIMESTAMP >= r.from_date AND r.till_date IS NULL
					) OR (
						DATE_PART('month', date) = DATE_PART('month', CURRENT_TIMESTAMP)
						AND DATE_PART('year', date) = DATE_PART('month', CURRENT_TIMESTAMP)
					)
				) 
			WHERE ac.self = false
			AND ac.user_id = %d
			GROUP BY ac.id
		) AS planned
		WHERE target > current`,
		dir, userID,
	)
	st.Orm.Raw(query).Scan(&res)
	return res.Sum
}

// FutureSavings calculates savings till a future date
func (st *State) FutureSavings(userID int, fut time.Time) (savings []Saving) {
	st.Orm.Raw(`SELECT month, income, charge, (income - charge) AS effect, SUM(income-charge) OVER (ORDER BY month) AS net_effect
	FROM (
		SELECT month, COALESCE((
			SELECT SUM(r1.amount) 
			FROM records AS r1
			JOIN accounts AS ai1 ON r1.account_in_id = ai1.id 
			JOIN accounts AS ao1 ON r1.account_out_id = ao1.id 
			WHERE (
				(
					month::date BETWEEN r1.from_date AND r1.till_date 
					OR
					month::date >= r1.from_date AND r1.till_date IS NULL
				)
				AND r1.mandate
				AND ao1.self = false AND ai1.self
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
					month::date BETWEEN r1.from_date AND r1.till_date 
					OR
					month::date >= r1.from_date AND r1.till_date IS NULL
				)
				AND r1.mandate
				AND ao1.self AND ai1.self = false
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
						date_trunc('month', CAST(? AS TIMESTAMP)),
						date_trunc('month', current_timestamp)
					) AS diff
				) AS fut
			) AS reps
		) AS future
	) AS savings`, userID, userID, fut).Scan(&savings)
	return
}

// Predict generates prediction for userID till fut
func (st *State) Predict(fut time.Time, userID int) (out string) {
	cash := st.CashTillNow(userID)
	cost := st.PlannedCurrentMonth(userID, true)
	income := st.PlannedCurrentMonth(userID, false)
	monthEnd := cash - cost + income
	savings := st.FutureSavings(userID, fut)

	out += fmt.Sprintf("*Prediction till EOM %s*\n\n", fut.Format(monthFormat))
	out += fmt.Sprintf("*%s:* %d EUR\n", time.Now().Format(monthFormat), monthEnd/100)
	out += fmt.Sprintf("Planned Expenses: %d EUR\n", cost/100)
	out += fmt.Sprintf("Receivables: %d EUR\n", income/100)
	out += fmt.Sprintf("Assets: %d EUR", cash/100)

	for i, save := range savings {
		if len(savings)-i <= 12 {
			// show only last twelve months
			// because of telegram message size limitation
			out += fmt.Sprintf(
				"\n\n*%s:* %d EUR\nCharge: %d EUR\nIncome: %d EUR\nEffect: %d EUR",
				save.Month.Format(monthFormat), (monthEnd+save.NetEffect)/100, save.Charge/100, save.Income/100, save.Effect/100,
			)
		}
	}
	return
}
