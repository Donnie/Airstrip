package main

import (
	"fmt"

	"github.com/Donnie/Airstrip/ptr"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

func (convo *Convo) askAmount() {
	convo.Expect = ptr.String("amount")
}

func (convo *Convo) askAccountIn(db *gorm.DB) {
	convo.getRecentAccountBtns(db, 0, "in")
	convo.Expect = ptr.String("account in")
}

func (convo *Convo) askCreateAccount(input string) {
	convo.Expect = ptr.String("account new")
	convo.menu.Inline(
		convo.menu.Row(
			convo.menu.Data("Yes", "y-"+input),
			convo.menu.Data("No", "n-"+input),
		),
	)
}

func (convo *Convo) askAccountOut(db *gorm.DB) {
	convo.getRecentAccountBtns(db, 0, "out")
	convo.Expect = ptr.String("account out")
}

func (convo *Convo) askDescription() {
	convo.Expect = ptr.String("description")
}

func (convo *Convo) askDate() {
	convo.Expect = ptr.String("date")
	convo.menu.Inline(
		convo.menu.Row(
			convo.menu.Data("Now", "now"),
			convo.menu.Data("Today", "today"),
		),
		convo.menu.Row(
			convo.menu.Data("Yesterday", "yesterday"),
			convo.menu.Data("Tomorrow", "tomorrow"),
		),
	)
}

func (convo *Convo) askFromDate() {
	convo.Expect = ptr.String("from date")
}

func (convo *Convo) askTillDate() {
	convo.Expect = ptr.String("till date")
}

func (convo *Convo) end() {
	convo.Expect = nil
}

func (convo *Convo) getRecentAccountBtns(db *gorm.DB, start int, inOut string) {
	var accounts []Account
	db.Where("records.mandate = false").
		Joins(fmt.Sprintf("JOIN records ON records.account_%s_id = accounts.id", inOut)).
		Group("accounts.id").Order("MAX(records.date) desc, accounts.id").
		Limit(8).Offset(start).Find(&accounts)

	var btns []tb.Btn
	for i := 0; i < getMin(4, len(accounts)); i++ {
		btns = append(btns, convo.menu.Data(*accounts[i].Name, *accounts[i].Name))
	}
	rowOne := convo.menu.Row(btns...)

	btns = []tb.Btn{}
	for i := 4; i < getMin(7, len(accounts)); i++ {
		btns = append(btns, convo.menu.Data(*accounts[i].Name, *accounts[i].Name))
	}

	if len(accounts) > 7 {
		btns = append(btns, convo.menu.Data("More...", fmt.Sprintf("more-%d", start+7)))
	}

	rowTwo := convo.menu.Row(btns...)
	convo.menu.Inline(rowOne, rowTwo)
}

func genQues(ask string) (out string) {
	switch ask {
	case "account new":
		out = "No account found by that name\\. Create one?"
	case "account name":
		out = "What is the new account name?"
	case "from date", "till date":
		out = fmt.Sprintf("What is the %s?\n\nSpecify in this format: *_Jan 2006_*", ask)
	default:
		out = fmt.Sprintf("What is the %s?", ask)
	}
	return
}
