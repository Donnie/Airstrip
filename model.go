package main

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

var layout = "Jan 2006"

// State holds fundamental items
type State struct {
	Bot *tb.Bot
	Orm *gorm.DB
	Env *Env
}

// Env to hold env vars
type Env struct {
	PORT      string
	TELETOKEN string
	WEBHOOK   string
}

// Record represents one Record
type Record struct {
	gorm.Model

	AccountIn    *Account   `json:"account_in"`
	AccountInID  *uint      `json:"account_in_id"`
	AccountOut   *Account   `json:"account_out"`
	AccountOutID *uint      `json:"account_out_id"`
	Amount       *int64     `json:"amount"`
	Date         *time.Time `json:"date"`
	Description  *string    `json:"description"`
	FromDate     *time.Time `json:"from_date"`
	Mandate      *bool      `json:"mandate"`
	TillDate     *time.Time `json:"till_date"`
	UserID       *int       `json:"user_id"`
}

// Expector is a function which expects a contextual response
type Expector func(*gorm.DB, string)

// Convo stores a conversation context
// the non-db fields are used to manipulate the bot
type Convo struct {
	gorm.Model

	UserID    *int    `json:"user_id"`
	Expect    *string `json:"expect"`
	ContextID *uint   `json:"context_id"`

	// non db fields
	handlers map[string]Expector `gorm:"-"`
	response string              `gorm:"-"`
	menu     tb.ReplyMarkup      `gorm:"-"`
}

// Account represents a partaker in a transaction
type Account struct {
	gorm.Model

	Currency *string `json:"currency"`
	Name     *string `json:"name"`
	Self     *bool   `json:"self"`
	UserID   *int    `json:"user_id"`
}

// Line is a transaction description
type Line struct {
	Amount float64
	Name   string
	Type   string
}

// View consists of lines
type View struct {
	Lines []Line
	Total float64
	Type  string
}

// Saving is per month net effect
type Saving struct {
	Month  string
	Income int
	Charge int
	Effect int
}
