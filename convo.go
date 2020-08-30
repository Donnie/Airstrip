package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Donnie/Airstrip/ptr"
)

func (convo *Convo) expectWhat(gl *Global, expect string) (out string) {
	context := &Record{}
	json.Unmarshal([]byte(*convo.Context), &context)

	switch *convo.Expect {
	case "account":
		context.Account = &expect
		gl.Orm.Save(&context)

		cont, _ := json.Marshal(context)
		convo.Context = ptr.String(string(cont))
		convo.Expect = ptr.String("amount")
		gl.Orm.Save(&convo)
	case "amount":
		amountFlt, err := strconv.ParseFloat(expect, 64)
		if err != nil {
			break
		}
		context.Amount = ptr.Int64(int64(amountFlt * 100))
		gl.Orm.Save(&context)
		cont, _ := json.Marshal(context)
		convo.Context = ptr.String(string(cont))
		convo.Expect = ptr.String("currency")
		gl.Orm.Save(&convo)
	case "currency":
		currency := expect
		if len(currency) != 3 {
			currency = "EUR"
		}
		context.Currency = &currency
		gl.Orm.Save(&context)
		cont, _ := json.Marshal(context)
		convo.Context = ptr.String(string(cont))
		convo.Expect = ptr.String("description")
		gl.Orm.Save(&convo)
	case "description":
		context.Description = &expect
		gl.Orm.Save(&context)
		cont, _ := json.Marshal(context)
		convo.Context = ptr.String(string(cont))
		if *context.Type == "variable" {
			convo.Expect = ptr.String("date")
		} else {
			convo.Expect = ptr.String("from date")
		}
		gl.Orm.Save(&convo)
	case "date":
		date := expect
		layout := "2006-01-02 15:04"
		dateTime, err := time.Parse(layout, date)
		if err != nil {
			dateTime = time.Now()
		}
		context.Date = &dateTime
		gl.Orm.Save(&context)
		convo.Expect = nil
		gl.Orm.Delete(&convo)
	case "from date":
		date := expect
		layout := "Jan 2006"
		dateTime, _ := time.Parse(layout, date)
		context.FromDate = &dateTime
		gl.Orm.Save(&context)
		cont, _ := json.Marshal(context)
		convo.Context = ptr.String(string(cont))
		convo.Expect = ptr.String("till date")
		gl.Orm.Save(&convo)
	case "till date":
		date := expect
		layout := "Jan 2006"
		dateTime, err := time.Parse(layout, date)
		if err == nil {
			dateTime = dateTime.AddDate(0, 1, -1)
			context.TillDate = &dateTime
			gl.Orm.Save(&context)
		}
		if *context.Form == "lend" {
			context.Form = ptr.String("expense")
			gl.Orm.Create(&context)
		} else if *context.Form == "loan" {
			context.Form = ptr.String("gain")
			gl.Orm.Create(&context)
		}
		convo.Expect = nil
		gl.Orm.Delete(&convo)
	}

	if convo.Expect != nil {
		out = genQues(*convo.Expect, *context.Form)
	} else {
		out = "Record stored!"
	}
	return
}
