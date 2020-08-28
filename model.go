package main

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
)

// Global holds fundamental items
type Global struct {
	Bot *tb.Bot
	Orm *gorm.DB
}

// Variable represents one variable cost
type Variable struct {
	ID        *int64     `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`

	Amount      *int64     `json:"amount"`
	Currency    *string    `json:"currency"`
	Date        *time.Time `json:"date"`
	Description *string    `json:"description"`
	Form        *string    `json:"form"`
	Name        *string    `json:"name"`
}

// Fixed represents fixed costs
type Fixed struct {
	ID        *int64     `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`

	Amount      *int64     `json:"amount"`
	Currency    *string    `json:"currency"`
	Description *string    `json:"description"`
	Form        *string    `json:"form"`
	FromDate    *time.Time `json:"fromdate"`
	Name        *string    `json:"name"`
	TillDate    *time.Time `json:"tilldate"`
}

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
