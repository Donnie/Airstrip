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

func (convo *Convo) askCreateAccountIn(input string) {
	convo.Expect = ptr.String("account new in")
	convo.menu.Inline(
		convo.menu.Row(
			convo.menu.Data("Yes", "y-"+input),
			convo.menu.Data("No", "n-"+input),
		),
	)
}

func (convo *Convo) askCreateAccountSelfIn(input string) {
	convo.Expect = ptr.String("account new self in")
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

func (convo *Convo) askCreateAccountOut(input string) {
	convo.Expect = ptr.String("account new out")
	convo.menu.Inline(
		convo.menu.Row(
			convo.menu.Data("Yes", "y-"+input),
			convo.menu.Data("No", "n-"+input),
		),
	)
}

func (convo *Convo) askCreateAccountSelfOut(input string) {
	convo.Expect = ptr.String("account new self out")
	convo.menu.Inline(
		convo.menu.Row(
			convo.menu.Data("Yes", "y-"+input),
			convo.menu.Data("No", "n-"+input),
		),
	)
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
	numberRows := 3
	btnsPerRow := 5

	db.Where(`records.mandate = "f"`).
		Where("records.user_id = ?", *convo.UserID).
		Joins(fmt.Sprintf("JOIN records ON records.account_%s_id = accounts.id", inOut)).
		// get latest accounts and not most popular accounts
		// MAX because you need an aggregate function #postgres
		Group("accounts.id").Order("MAX(records.date) desc, accounts.id").
		Limit(numberRows * btnsPerRow).Offset(start).Find(&accounts)

	var rows [][]tb.Btn
	for i := 0; i < numberRows; i++ {
		var btns []tb.Btn
		for j := i * btnsPerRow; j < getMin(getMin((i+1)*btnsPerRow, (numberRows*btnsPerRow-1)), len(accounts)); j++ {
			fmt.Println(getMin((i+1)*btnsPerRow, (numberRows*btnsPerRow - 1)))
			btns = append(btns, convo.menu.Data(*accounts[j].Name, *accounts[j].Name))
		}
		if i == (numberRows-1) && len(accounts) > (numberRows*btnsPerRow-1) {
			btns = append(btns, convo.menu.Data("More...", fmt.Sprintf("more-%d", start+(numberRows*btnsPerRow-1))))
		}
		rows = append(rows, convo.menu.Row(btns...))
	}
	convo.menu.Inline(rows[0], rows[1], rows[2])
}

func genQues(ask string) (out string) {
	switch ask {
	case "account in":
		out = "Which account to be credited?"
	case "account out":
		out = "Which account to be debited?"
	case "account new in", "account new out":
		out = "No account found by that name\\. Create one?"
	case "account new self in", "account new self out":
		out = "Is this your own account?"
	case "from date", "till date":
		out = fmt.Sprintf("What is the %s?\n\nSpecify in this format: *_Jan 2006_*", ask)
	default:
		out = fmt.Sprintf("What is the %s?", ask)
	}
	return
}
