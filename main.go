package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	if _, err := os.Stat(".env.local"); os.IsNotExist(err) {
		godotenv.Load(".env")
	} else {
		godotenv.Load(".env.local")
	}
	fmt.Println("Running for " + os.Getenv("ENV"))
}

func main() {
	teleToken, exists := os.LookupEnv("TELEGRAM_TOKEN")
	if !exists {
		fmt.Println("Add TELEGRAM_TOKEN to .env file")
		os.Exit(1)
	}

	filename, exists := os.LookupEnv("DBFILE")
	if !exists {
		fmt.Println("Add DBFILE to .env file")
		os.Exit(1)
	}

	b, err := tb.NewBot(tb.Settings{
		Token:  teleToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	db, err := gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Variable{}, &Fixed{})

	gl := Global{
		Bot: b,
		Orm: db,
	}

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, fmt.Sprintf("Hello %s!", m.Sender.FirstName))
	})

	b.Handle("/charge", gl.handleCharge)
	b.Handle("/expense", gl.handleExpense)
	b.Handle("/gain", gl.handleGain)
	b.Handle("/income", gl.handleIncome)
	b.Handle("/lend", gl.handleLend)
	b.Handle("/loan", gl.handleLoan)
	b.Handle("/predict", gl.handlePredict)
	b.Handle("/show", gl.handleShow)

	b.Start()
}
