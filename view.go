package main

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (gl *Global) handleView(m *tb.Message) {
	var layout = "Jan 2006"
	t, err := time.Parse(layout, m.Payload)
	if err != nil {
		t = time.Now()
	}

	recs := []Record{}
	gl.Orm.
		Where(
			gl.Orm.Where("(?::date BETWEEN from_date AND till_date)", time.Now()).
				Or("(?::date >= from_date AND till_date IS NULL)", time.Now()),
		).
		Or(
			gl.Orm.Where("EXTRACT(MONTH FROM date) = ?", int(t.Month())).
				Where("EXTRACT(YEAR FROM date) = ?", t.Year()),
		).
		Find(&recs)

	expenses := []string{}
	gains := []string{}
	loans := []string{}
	lends := []string{}
	incomes := []string{}
	charges := []string{}
	output := fmt.Sprintf("*Overview of %s*\n", t.Format(layout))

	for _, rec := range recs {
		switch *rec.Form {
		case "expense":
			expenses = append(expenses, fmt.Sprintf("`%s %d %s`\n", *rec.Account, *rec.Amount/100, *rec.Currency))
		case "gain":
			gains = append(gains, fmt.Sprintf("`%s %d %s`\n", *rec.Account, *rec.Amount/100, *rec.Currency))
		case "loan":
			loans = append(loans, fmt.Sprintf("`%s %d %s`\n", *rec.Account, *rec.Amount/100, *rec.Currency))
		case "lend":
			lends = append(lends, fmt.Sprintf("`%s %d %s`\n", *rec.Account, *rec.Amount/100, *rec.Currency))
		case "income":
			incomes = append(incomes, fmt.Sprintf("`%s %d %s`\n", *rec.Account, *rec.Amount/100, *rec.Currency))
		case "charge":
			charges = append(charges, fmt.Sprintf("`%s %d %s`\n", *rec.Account, *rec.Amount/100, *rec.Currency))
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

	output += "\n*Loans*\n"
	for _, item := range loans {
		output += item
	}

	output += "\n*Lends*\n"
	for _, item := range lends {
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

	gl.Bot.Send(m.Sender, output, tb.ModeMarkdown)
}
