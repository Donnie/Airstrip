package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var layout = "Jan 2006"

func (glob *Global) handleHook(c *gin.Context) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	str := buf.String()

	dat, err := ioutil.ReadFile(glob.File)
	check(err)
	account := &Account{}
	json.Unmarshal(dat, account)

	var input Input

	err = json.Unmarshal([]byte(str), &input)
	check(err)

	date := *input.Message.Text
	matched, _ := regexp.Match(`^(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) \d{4}`, []byte(date))
	if !matched {
		date = "Dec 2030"
	}
	t, _ := time.Parse(layout, date)

	output := strconv.Itoa(int(predictFuture(t, *account, false)))

	glob.sendMessage(*input.Message.Chat.ID, output, input.Message.MessageID)

	c.JSON(200, nil)
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

	cash = cash + (int32(reps) * carry)
	return
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

func (glob *Global) sendMessage(chatID int64, text string, messageID *int64) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true

	if messageID != nil {
		msg.ReplyToMessageID = int(*messageID)
	}
	glob.Bot.Send(msg)
}
