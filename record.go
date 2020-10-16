package main

import (
	"strings"
	"time"

	"github.com/Donnie/Airstrip/ptr"
	"github.com/jinzhu/now"
	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) handleRecord(m *tb.Message) {
	userID := m.Sender.ID

	// end last conversation
	st.Orm.Unscoped().Where("user_id = ?", userID).Delete(&Convo{})

	// Init empty record
	item := &Record{
		UserID:  &userID,
		Mandate: ptr.Bool(strings.Contains(m.Text, "recur")),
	}
	st.Orm.Create(&item)

	// Create new conversation with Context
	convo := &Convo{
		ContextID: &item.ID,
		Expect:    ptr.String("amount"),
		UserID:    &userID,
	}
	st.Orm.Create(&convo)

	convo.response = genQues("amount")
	st.Bot.Send(m.Sender, convo.response, &convo.menu)
}

func (rec *Record) isExpense() (out bool) {
	if !*rec.Mandate &&
		!*rec.AccountIn.Self &&
		*rec.AccountOut.Self {
		out = true
	}
	return
}

func (rec *Record) isCharge() (out bool) {
	if *rec.Mandate &&
		!*rec.AccountIn.Self &&
		*rec.AccountOut.Self {
		out = true
	}
	return
}

func (rec *Record) isGain() (out bool) {
	if !*rec.Mandate &&
		*rec.AccountIn.Self &&
		!*rec.AccountOut.Self {
		out = true
	}
	return
}

func (rec *Record) isIncome() (out bool) {
	if *rec.Mandate &&
		*rec.AccountIn.Self &&
		!*rec.AccountOut.Self {
		out = true
	}
	return
}

func (rec *Record) isCurrent() (out bool) {
	if rec.Date != nil &&
		now.BeginningOfMonth().Unix() <= rec.Date.Unix() &&
		now.EndOfMonth().Unix() >= rec.Date.Unix() {
		out = true
	} else if rec.FromDate != nil &&
		now.BeginningOfMonth().Unix() >= rec.FromDate.Unix() &&
		(rec.TillDate == nil ||
			now.BeginningOfMonth().Unix() <= rec.TillDate.Unix()) {
		out = true
	}
	return
}

func (rec *Record) isOfMonth(month time.Time) (out bool) {
	if month.Unix() >= rec.FromDate.Unix() &&
		(rec.TillDate == nil ||
			month.Unix() <= rec.TillDate.Unix()) {
		out = true
	}
	return
}
