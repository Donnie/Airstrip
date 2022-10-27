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

	query := fmt.Sprintf(`SELECT ai.name, CAST(r.amount AS REAL)/100 AS amount, (
		CASE
		WHEN ai.self = 1 AND ao.self = 1 THEN 'Transfers'
		WHEN r.mandate = 0 AND ai.self = 0 THEN 'Expenses'
		WHEN r.mandate = 0 AND ai.self = 1 THEN 'Gains'
		WHEN r.mandate = 1 AND ai.self = 1 THEN 'Incomes'
		WHEN r.mandate = 1 AND ai.self = 0 THEN 'Charges'
		END
	) as type 
	FROM records AS r 
	JOIN accounts AS ai ON r.account_in_id = ai.id 
	JOIN accounts AS ao ON r.account_out_id = ao.id 
	WHERE (
		(
			(
				"%s" BETWEEN r.from_date AND r.till_date
				OR
				"%s" >= r.from_date AND r.till_date IS NULL
			) OR (
				LTRIM(STRFTIME('%%m', date), "0") = "%d"
				AND
				STRFTIME('%%Y', date) = "%d"
			)
		)
		AND r.user_id = %d 
		AND r.deleted_at IS NULL
	) ORDER BY r.date, r.from_date`,
		t.Format("2006-01-02"), t.Format("2006-01-02"),
		int(t.Month()), t.Year(), m.Sender.ID,
	)
	st.Orm.Raw(query).Scan(&lines)

	output := fmt.Sprintf("*Overview of %s*\n", t.Format(monthFormat))
	output += prepareView(viewLines(lines))

	st.Bot.Send(m.Sender, output, tb.ModeHTML)
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
		out += fmt.Sprintf("\n<strong>%s</strong>: <code>€%.2f<code>\n", view.Type, view.Total)
		for _, line := range view.Lines {
			out += fmt.Sprintf("<code>%s: €%.2f<code>\n", line.Name, line.Amount)
		}
	}
	return
}
