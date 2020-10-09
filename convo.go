package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Donnie/Airstrip/ptr"
	"gorm.io/gorm"
)

// Handle handles a conversation by mapping expectors
// with contextual answer
func (convo *Convo) Handle(expect string, handler Expector) {
	convo.handlers[expect] = handler
}

func (convo *Convo) expectNext(db *gorm.DB, expect string) {
	if convo.Expect != nil {
		convo.handlers[*convo.Expect](db, expect)
	}

	if convo.Expect != nil {
		db.Save(&convo)
		convo.response = genQues(*convo.Expect)
		return
	}
	db.Unscoped().Delete(&convo)
	convo.response = "Record stored\\!"
}

func (convo *Convo) expectAccountIn(db *gorm.DB, input string) {
	// find out account
	var account Account
	err := db.Where("name = ?", input).
		Where("user_id = ?", *convo.UserID).
		First(&account).Error

	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		convo.Expect = ptr.String("account que")
		convo.menu.Inline(
			convo.menu.Row(
				convo.menu.Data("Yes", "y-"+input),
				convo.menu.Data("No", "n-"+input),
			),
		)
		return
	}

	db.Model(&Record{}).
		Where("id = ?", *convo.ContextID).
		Update("account_in_id", account.ID)
	convo.Expect = ptr.String("account out")
}

func (convo *Convo) expectAccountOut(db *gorm.DB, input string) {
	// find out account
	var account Account
	err := db.Where("name = ?", input).First(&account).Error

	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		convo.Expect = ptr.String("account out")
		return
	}

	db.Model(&Record{}).
		Where("id = ?", *convo.ContextID).
		Update("account_out_id", account.ID)
	convo.Expect = ptr.String("amount")
}

func (convo *Convo) expectAccountQue(db *gorm.DB, input string) {
	inp := strings.Split(input, "-")
	switch strings.ToLower(inp[0]) {
	case "y":
		var account Account
		account.Name = &inp[1]
		account.UserID = convo.UserID
		db.Create(&account)

		db.Model(&Record{}).
			Where("id = ?", *convo.ContextID).
			Update("account_in_id", account.ID)
		convo.Expect = ptr.String("account out")
	default:
		convo.Expect = ptr.String("account in")
	}
}

func (convo *Convo) expectAmount(db *gorm.DB, input string) {
	flt, err := strconv.ParseFloat(input, 64)
	if err != nil || flt < 0 {
		return
	}
	db.Model(&Record{}).
		Where("id = ?", *convo.ContextID).
		Update("amount", int(math.Round(flt*100)))
	convo.Expect = ptr.String("description")
}

func (convo *Convo) expectDate(db *gorm.DB, input string) {
	dateTime := parseDate(input)
	if dateTime.IsZero() {
		return
	}

	db.Model(&Record{}).Where("id = ?", *convo.ContextID).
		Update("date", dateTime)
	convo.Expect = nil
}

func (convo *Convo) expectDescription(db *gorm.DB, input string) {
	record := &Record{}
	db.First(&record, *convo.ContextID)
	record.Description = &input
	db.Save(&record)
	if *record.Mandate {
		convo.Expect = ptr.String("from date")
		return
	}
	convo.Expect = ptr.String("date")
}

func (convo *Convo) expectFromDate(db *gorm.DB, input string) {
	layout := "Jan 2006"
	dateTime, err := time.Parse(layout, input)
	if err != nil {
		return
	}
	db.Model(&Record{}).Where("id = ?", *convo.ContextID).Update("from_date", dateTime)
	convo.Expect = ptr.String("till date")
}

func (convo *Convo) expectTillDate(db *gorm.DB, input string) {
	layout := "Jan 2006"
	dateTime, err := time.Parse(layout, input)
	if err != nil {
		return
	}
	dateTime = dateTime.AddDate(0, 1, -1)

	record := &Record{}
	db.First(&record, *convo.ContextID)
	record.TillDate = &dateTime
	db.Save(&record)
	convo.Expect = nil
}

func genQues(ask string) (out string) {
	switch ask {
	case "account que":
		out = "No account found by that name\\. Create one?"
	case "account choose in", "account choose out":
		out = "More than one account found\\. Be more specific\\."
	case "account name":
		out = "What is the new account name?"
	case "from date", "till date":
		out = fmt.Sprintf("What is the %s?\n\nSpecify in this format: *_Jan 2006_*", ask)
	default:
		out = fmt.Sprintf("What is the %s?", ask)
	}
	return
}
