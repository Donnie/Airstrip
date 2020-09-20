package main

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleView(m *tb.Message) {
	var layout = "Jan 2006"
	t, err := time.Parse(layout, m.Payload)
	if err != nil {
		t = time.Now()
	}

	recs := []Record{}
	st.Orm.Preload("AccountIn").
		Where(
			st.Orm.Where("(?::date BETWEEN from_date AND till_date)", t).
				Or("(?::date >= from_date AND till_date IS NULL)", t),
		).
		Or(
			st.Orm.Where("EXTRACT(MONTH FROM date) = ?", int(t.Month())).
				Where("EXTRACT(YEAR FROM date) = ?", t.Year()),
		).
		Find(&recs)

	output := fmt.Sprintf("*Overview of %s*\n", t.Format(layout))
	output += prepareView(recs)

	st.Bot.Send(m.Sender, output, tb.ModeMarkdown)
}

func prepareView(recs []Record) (output string) {
	var expenses, gains, incomes, charges []string

	for _, rec := range recs {
		switch *rec.AccountIn.Self {
		case true:
			switch *rec.Mandate {
			case true:
				incomes = append(incomes, fmt.Sprintf("`%s %d %s`\n", *rec.AccountIn.Name, *rec.Amount/100, *rec.AccountIn.Currency))
			case false:
				gains = append(gains, fmt.Sprintf("`%s %d %s`\n", *rec.AccountIn.Name, *rec.Amount/100, *rec.AccountIn.Currency))
			}
		case false:
			switch *rec.Mandate {
			case true:
				charges = append(charges, fmt.Sprintf("`%s %d %s`\n", *rec.AccountIn.Name, *rec.Amount/100, *rec.AccountIn.Currency))
			case false:
				expenses = append(expenses, fmt.Sprintf("`%s %d %s`\n", *rec.AccountIn.Name, *rec.Amount/100, *rec.AccountIn.Currency))
			}
		}
	}

	output += "\n*Expenses*\n"
	for _, item := range expenses {
		output += item
	}

	output += "\n*Gains*\n"
	for _, item := range gains {
		output += item
	}

	output += "\n*Incomes*\n"
	for _, item := range incomes {
		output += item
	}

	output += "\n*Charges*\n"
	for _, item := range charges {
		output += item
	}

	return
}
