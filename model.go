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

// Record represents one Record
type Record struct {
	ID        *int64     `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`

	AccountID   *int64     `json:"account_id"`
	Account     *Account   `json:"account"`
	Amount      *int64     `json:"amount"`
	Currency    *string    `json:"currency"`
	Date        *time.Time `json:"date"`
	Description *string    `json:"description"`
	Form        *string    `json:"form"`
	FromDate    *time.Time `json:"from_date"`
	TillDate    *time.Time `json:"till_date"`
	Type        *string    `json:"type"`
	UserID      *int64     `json:"user_id"`
}

// Expector is a function which expects a contextual response
type Expector func(*Global, string)

// Convo stores a conversation context
type Convo struct {
	ID        *int64     `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`

	UserID    *int64              `json:"user_id"`
	Expect    *string             `json:"expect"`
	ContextID *int64              `json:"context_id"`
	handlers  map[string]Expector `gorm:"-"`
}

// Account represents a partaker in a transaction
type Account struct {
	ID        *int64     `json:"id" gorm:"primaryKey"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`

	Name *string `json:"name"`
}
