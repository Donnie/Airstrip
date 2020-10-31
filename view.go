package main

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleView(m *tb.Message) {
	t, err := time.Parse(monthFormat, m.Payload)
	if err != nil {
		t = time.Now()
	}

	var lines []Line
	st.Orm.Raw(`SELECT ai.name, CAST(r.amount AS DOUBLE PRECISION)/100 AS amount, (
		CASE
		WHEN ai.self AND ao.self THEN 'Transfers' 
		WHEN r.mandate = false AND ai.self = false THEN 'Expenses' 
		WHEN r.mandate = false AND ai.self THEN 'Gains' 
		WHEN r.mandate AND ai.self THEN 'Incomes' 
		WHEN r.mandate AND ai.self = false THEN 'Charges'
		END
	) as type 
	FROM records AS r 
	JOIN accounts AS ai ON r.account_in_id = ai.id 
	JOIN accounts AS ao ON r.account_out_id = ao.id 
	WHERE (
		(
			(
				?::date BETWEEN r.from_date AND r.till_date 
				OR
				?::date >= r.from_date AND r.till_date IS NULL
			) OR (
				EXTRACT(MONTH FROM date) = ? 
				AND
				EXTRACT(YEAR FROM date) = ?
			)
		)
		AND r.user_id = ?
	)`, t, t, int(t.Month()), t.Year(), m.Sender.ID).Scan(&lines)

	output := fmt.Sprintf("*Overview of %s*\n", t.Format(monthFormat))
	output += prepareView(viewLines(lines))

	st.Bot.Send(m.Sender, output, tb.ModeMarkdown)
}

func viewLines(lines []Line) (views [4]View) {
	views[0].Type = "Expenses"
	views[1].Type = "Gains"
	views[2].Type = "Incomes"
	views[3].Type = "Charges"
	for _, rec := range lines {
		for i := range views {
			if rec.Type == views[i].Type {
				views[i].Lines = append(views[i].Lines, rec)
				views[i].Total += rec.Amount
			}
		}
	}
	return
}

func prepareView(views [4]View) (out string) {
	for _, view := range views {
		out += fmt.Sprintf("\n*%s*: `%.2f` EUR\n", view.Type, view.Total)
		for _, line := range view.Lines {
			out += fmt.Sprintf("`%s: `%.2f` EUR`\n", line.Name, line.Amount)
		}
	}
	return
}
