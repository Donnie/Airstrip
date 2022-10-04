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
	st.Bot.Send(m.Sender, predictions, tb.ModeHTML)
}

// PlannedCurrentMonth calculates remaining costs and incomes for current month
func (st *State) PlannedCurrentMonth(userID int64, cost bool) int {
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
				AND r.deleted_at IS NULL
				AND (
					(
						CURRENT_TIMESTAMP BETWEEN r.from_date AND r.till_date
						OR CURRENT_TIMESTAMP >= r.from_date AND r.till_date IS NULL
					) OR (
						DATE_PART('month', date) = DATE_PART('month', CURRENT_TIMESTAMP)
						AND DATE_PART('year', date) = DATE_PART('year', CURRENT_TIMESTAMP)
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
func (st *State) FutureSavings(userID int64, fut time.Time) (savings []Saving) {
	st.Orm.Raw(`SELECT month, income, charge, (income - charge) AS effect, SUM(income-charge) OVER (ORDER BY month) AS net_effect
	FROM (
		SELECT month, COALESCE((
			SELECT SUM(r.amount) 
			FROM records AS r
			JOIN accounts AS ai1 ON r.account_in_id = ai1.id 
			JOIN accounts AS ao1 ON r.account_out_id = ao1.id 
			WHERE (
				(
					month::date BETWEEN r.from_date AND r.till_date 
					OR
					month::date >= r.from_date AND r.till_date IS NULL
				)
				AND r.mandate
				AND ao1.self = false AND ai1.self
				AND r.user_id = ?
				AND r.deleted_at IS NULL
			)
		), 0) AS income, 
		COALESCE((
			SELECT SUM(r.amount) 
			FROM records AS r
			JOIN accounts AS ai1 ON r.account_in_id = ai1.id 
			JOIN accounts AS ao1 ON r.account_out_id = ao1.id 
			WHERE (
				(
					month::date BETWEEN r.from_date AND r.till_date 
					OR
					month::date >= r.from_date AND r.till_date IS NULL
				)
				AND r.mandate
				AND ao1.self AND ai1.self = false
				AND r.user_id = ?
				AND r.deleted_at IS NULL
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
func (st *State) Predict(fut time.Time, userID int64) (out string) {
	cash := st.CashTillNow(userID)
	cost := st.PlannedCurrentMonth(userID, true)
	income := st.PlannedCurrentMonth(userID, false)
	monthEnd := cash - cost + income
	savings := st.FutureSavings(userID, fut)

	out += fmt.Sprintf("<strong>Prediction till EOM %s</strong><br /><br />", fut.Format(monthFormat))
	out += fmt.Sprintf("<strong>%s:</strong> `%d EUR`<br />", time.Now().Format(monthFormat), monthEnd/100)
	out += fmt.Sprintf("Planned Expenses: `%d EUR`<br />", cost/100)
	out += fmt.Sprintf("Receivables: `%d EUR`<br />", income/100)
	out += fmt.Sprintf("Assets: `%d EUR`", cash/100)

	for i, save := range savings {
		if len(savings)-i <= 12 {
			// show only last twelve months
			// because of telegram message size limitation
			out += fmt.Sprintf(
				"<br /><br /><strong>%s:</strong> `%d EUR`<br />Charge: `%d EUR`<br />Income: `%d EUR`<br />Effect: `%d EUR`",
				save.Month, (monthEnd+save.NetEffect)/100, save.Charge/100, save.Income/100, save.Effect/100,
			)
		}
	}
	return
}
