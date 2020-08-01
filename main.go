package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// Trans defines a money transaction
type Trans struct {
	Name   string `json:"name"`
	Amount int32  `json:"amount"`
}

// Account represents costs and earnings
type Account struct {
	Costs    []Trans `json:"costs"`
	Earnings []Trans `json:"earnings"`
	Savings  []Trans `json:"savings"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	dat, err := ioutil.ReadFile("db/db.json")
	check(err)

	account := &Account{}
	json.Unmarshal(dat, account)

	layout := "Jan 2006"
	str := "Sep 2020"
	t, _ := time.Parse(layout, str)

	fmt.Println(predictFuture(t, *account, false))
}

func predictFuture(future time.Time, account Account, start bool) (cash int32) {
	for _, trans := range account.Savings {
		cash = cash + trans.Amount
	}

	reps := monthDiff(time.Now(), future)
	if start {
		reps--
	}
	carry := calcMonthEnd(account)

	fmt.Println(reps)

	cash = cash + (int32(reps) * carry)

	return cash
}

func calcMonthEnd(account Account) (cash int32) {
	for _, trans := range account.Earnings {
		cash = cash + trans.Amount
	}
	for _, trans := range account.Costs {
		cash = cash - trans.Amount
	}

	return
}

func monthDiff(a, b time.Time) (month int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, m1, _ := a.Date()
	y2, m2, _ := b.Date()

	month = int(m2 - m1)
	year := int(y2 - y1)

	month = month + (year * 12)

	return
}
