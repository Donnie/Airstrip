package main

import (
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handlePredict(m *tb.Message) {
	var layout = "Jan 2006"

	date := m.Payload
	t, err := time.Parse(layout, date)
	if err != nil {
		t = time.Now().AddDate(1, 0, 0)
	} else {
		// end of the month
		t = t.AddDate(0, 1, -1)
	}

	recs := []Record{}
	st.Orm.Preload("AccountIn").Preload("AccountOut").
		Where("user_id = ?", m.Sender.ID).
		Find(&recs)

	cashCurr := st.CashTillNow(m.Sender.ID)
	costPlan := st.PlannedExpensesCurrentMonth(m.Sender.ID)
	cashFutr := calcFutr(t, recs)

	output := fmt.Sprintf("*Prediction* for month end %s:\n%d EUR", t.Format(layout), (cashCurr-costPlan+cashFutr)/100)
	st.Bot.Send(m.Sender, output, tb.ModeMarkdown)
}

// CashTillNow calculates Summation of all assets till now
func (st *State) CashTillNow(userID int) int64 {
	var res struct {
		Sum int64
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
func (st *State) PlannedExpensesCurrentMonth(userID int) int64 {
	var res struct {
		Sum int64
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

func calcFutr(future time.Time, recs []Record) (cash int64) {
	timeNow := time.Now().AddDate(0, 1, 0)
	reps := monthDiff(timeNow, future)
	if reps > 0 {
		for i := 0; i < reps; i++ {
			cash += calcMonthEnd(recs, timeNow.AddDate(0, i, 0))
		}
	}
	return
}

func calcMonthEnd(recs []Record, month time.Time) (cash int64) {
	for _, rec := range recs {
		if rec.isIncome() && rec.isOfMonth(month) {
			cash += *rec.Amount
			continue
		}
		if rec.isCharge() && rec.isOfMonth(month) {
			cash -= *rec.Amount
		}
	}
	return
}

func monthDiff(a, b time.Time) (month int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	y1, m1, _ := a.Date()
	y2, m2, _ := b.Date()

	month = int(m2 - m1 + 1)
	year := int(y2 - y1)

	month += (year * 12)

	return
}
