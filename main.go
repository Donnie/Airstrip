package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/postgres"
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

	bot, err := tb.NewBot(tb.Settings{
		Token:  teleToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	dsn := "user=airstrip password=postgres dbname=airstrip port=5432 sslmode=disable host=postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// db.AutoMigrate(&Account{}, &Convo{}, &Record{})

	gl := Global{
		Bot: bot,
		Orm: db,
	}

	bot.Handle("/start", gl.handleHelp)
	bot.Handle("/help", gl.handleHelp)
	bot.Handle("/charge", gl.handleRecord)
	bot.Handle("/expense", gl.handleRecord)
	bot.Handle("/delete", gl.handleDelete)
	bot.Handle("/gain", gl.handleRecord)
	bot.Handle("/income", gl.handleRecord)
	bot.Handle("/lend", gl.handleRecord)
	bot.Handle("/loan", gl.handleRecord)
	bot.Handle("/predict", gl.handlePredict)
	bot.Handle("/view", gl.handleView)

	bot.Handle(tb.OnText, gl.handleContext)
	bot.Handle(tb.OnCallback, gl.handleCallback)

	bot.Start()
}
