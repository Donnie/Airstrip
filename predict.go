package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handlePredict(m *tb.Message) {
	var layout = "Jan 2006"
	userID := int64(m.Sender.ID)

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
		Where("user_id = ?", userID).
		Find(&recs)

	cashCurr := calcCurr(recs)
	costPlan := plannedExp(recs)
	cashFutr := calcFutr(t, recs)

	output := fmt.Sprintf("*Prediction* for month end %s:\n%d EUR", t.Format(layout), (cashCurr-costPlan+cashFutr)/100)
	st.Bot.Send(m.Sender, output, tb.ModeMarkdown)
}

func calcCurr(recs []Record) (cash int64) {
	for _, rec := range recs {
		if !*rec.Mandate {
			if *rec.AccountIn.Self && !*rec.AccountOut.Self {
				cash += *rec.Amount
				continue
			}
			if !*rec.AccountIn.Self {
				cash -= *rec.Amount
			}
		}
	}
	return
}

func plannedExp(recs []Record) (cash int64) {
	// get current month expenses
	currExp := []Record{}
	for _, rec := range recs {
		if !*rec.AccountIn.Self && !*rec.Mandate &&
			(now.BeginningOfMonth().Unix() <= rec.Date.Unix() &&
				now.EndOfMonth().Unix() >= rec.Date.Unix()) {
			currExp = append(currExp, rec)
		}
	}

	// get current month charges
	currChg := []Record{}
	for _, rec := range recs {
		if !*rec.AccountIn.Self && *rec.Mandate &&
			(now.BeginningOfMonth().Unix() >= rec.FromDate.Unix() &&
				(rec.TillDate == nil ||
					now.BeginningOfMonth().Unix() <= rec.TillDate.Unix())) {
			currChg = append(currChg, rec)
		}
	}

	// find diff
	for _, chg := range currChg {
		tobespent := *chg.Amount
		for _, exp := range currExp {
			if *exp.AccountInID == *chg.AccountInID {
				tobespent -= *exp.Amount
			}
		}
		if tobespent > 0 {
			cash += tobespent
		}
	}
	return
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
		if *rec.Mandate &&
			(month.Unix() >= rec.FromDate.Unix() &&
				(rec.TillDate == nil ||
					month.Unix() <= rec.TillDate.Unix())) {
			if *rec.AccountIn.Self && !*rec.AccountOut.Self {
				cash += *rec.Amount
				continue
			}
			if !*rec.AccountIn.Self {
				cash -= *rec.Amount
			}
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
