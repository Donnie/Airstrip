package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
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

	gl := Global{
		File: filename,
		Bot:  b,
	}

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, fmt.Sprintf("Hello %s!", m.Sender.FirstName))
	})

	b.Handle("/predict", gl.handlePredict)

	b.Start()
}
