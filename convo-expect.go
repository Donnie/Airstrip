package main

import (
	"errors"
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

func (convo *Convo) expectAmount(db *gorm.DB, input string) {
	flt, err := strconv.ParseFloat(input, 64)
	if err != nil || flt < 0 {
		return
	}
	db.Model(&Record{}).
		Where("id = ?", *convo.ContextID).
		Update("amount", int(math.Round(flt*100)))

	convo.askAccountIn(db)
}

func (convo *Convo) expectAccountIn(db *gorm.DB, input string) {
	if strings.Contains(input, "more") {
		inp := strings.Split(input, "-")
		start, _ := strconv.Atoi(inp[1])
		convo.getRecentAccountBtns(db, start, "in")
		convo.Expect = ptr.String("account in")
		return
	}

	// find out account
	var account Account
	err := db.Where("name = ?", input).
		Where("user_id = ?", *convo.UserID).
		First(&account).Error

	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		convo.askCreateAccountIn(input)
		return
	}

	db.Model(&Record{}).
		Where("id = ?", *convo.ContextID).
		Update("account_in_id", account.ID)
	convo.askAccountOut(db)
}

func (convo *Convo) expectAccountInQue(db *gorm.DB, input string) {
	inp := strings.Split(input, "-")
	switch strings.ToLower(inp[0]) {
	case "y":
		var account Account
		account.Currency = ptr.String("EUR")
		account.Name = &inp[1]
		account.Self = ptr.Bool(false)
		account.UserID = convo.UserID
		db.Create(&account)

		db.Model(&Record{}).
			Where("id = ?", *convo.ContextID).
			Update("account_in_id", account.ID)

		convo.askCreateAccountSelfIn(inp[1])
	default:
		convo.askAccountIn(db)
	}
}

func (convo *Convo) expectCreateAccountSelfIn(db *gorm.DB, input string) {
	inp := strings.Split(input, "-")
	switch strings.ToLower(inp[0]) {
	case "y":
		db.Model(&Account{}).
			Where("name = ?", inp[1]).
			Where("user_id = ?", *convo.UserID).
			Update("self", true)

		convo.askAccountOut(db)
	default:
		convo.askAccountOut(db)
	}
}

func (convo *Convo) expectAccountOut(db *gorm.DB, input string) {
	if strings.Contains(input, "more") {
		inp := strings.Split(input, "-")
		start, _ := strconv.Atoi(inp[1])
		convo.getRecentAccountBtns(db, start, "out")
		convo.Expect = ptr.String("account out")
		return
	}

	// find out account
	var account Account
	err := db.Where("name = ?", input).
		Where("user_id = ?", *convo.UserID).
		First(&account).Error

	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		convo.askCreateAccountOut(input)
		return
	}

	db.Model(&Record{}).
		Where("id = ?", *convo.ContextID).
		Update("account_out_id", account.ID)
	convo.askDescription()
}

func (convo *Convo) expectAccountOutQue(db *gorm.DB, input string) {
	inp := strings.Split(input, "-")
	switch strings.ToLower(inp[0]) {
	case "y":
		var account Account
		account.Currency = ptr.String("EUR")
		account.Name = &inp[1]
		account.Self = ptr.Bool(false)
		account.UserID = convo.UserID
		db.Create(&account)

		db.Model(&Record{}).
			Where("id = ?", *convo.ContextID).
			Update("account_out_id", account.ID)

		convo.askCreateAccountSelfOut(inp[1])
	default:
		convo.askAccountOut(db)
	}
}

func (convo *Convo) expectCreateAccountSelfOut(db *gorm.DB, input string) {
	inp := strings.Split(input, "-")
	switch strings.ToLower(inp[0]) {
	case "y":
		db.Model(&Account{}).
			Where("name = ?", inp[1]).
			Where("user_id = ?", *convo.UserID).
			Update("self", true)

		convo.askDescription()
	default:
		convo.askDescription()
	}
}

func (convo *Convo) expectDescription(db *gorm.DB, input string) {
	record := &Record{}
	db.First(&record, *convo.ContextID)
	record.Description = &input
	db.Save(&record)
	if *record.Mandate {
		convo.askFromDate()
		return
	}
	convo.askDate()
}

func (convo *Convo) expectDate(db *gorm.DB, input string) {
	dateTime := parseDate(input)
	if dateTime.IsZero() {
		return
	}

	db.Model(&Record{}).
		Where("id = ?", *convo.ContextID).
		Update("date", dateTime)
	convo.end()
}

func (convo *Convo) expectFromDate(db *gorm.DB, input string) {
	dateTime, err := time.Parse(monthFormat, input)
	if err != nil {
		return
	}
	db.Model(&Record{}).
		Where("id = ?", *convo.ContextID).
		Update("from_date", dateTime)
	convo.askTillDate()
}

func (convo *Convo) expectTillDate(db *gorm.DB, input string) {
	dateTime, err := time.Parse(monthFormat, input)
	if err != nil {
		return
	}
	dateTime = dateTime.AddDate(0, 1, -1)

	db.Model(&Record{}).
		Where("id = ?", *convo.ContextID).
		Update("till_date", dateTime)
	convo.end()
}
