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

	b, err := tb.NewBot(tb.Settings{
		Token:  teleToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	dsn := "user=airstrip password=postgres dbname=airstrip port=5432 sslmode=disable host=airstrip-db"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Record{}, &Convo{})

	gl := Global{
		Bot: b,
		Orm: db,
	}

	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Sender, fmt.Sprintf("Hello %s!", m.Sender.FirstName))
	})

	b.Handle("/charge", gl.handleRecord)
	b.Handle("/expense", gl.handleRecord)
	b.Handle("/gain", gl.handleRecord)
	b.Handle("/income", gl.handleRecord)
	b.Handle("/lend", gl.handleRecord)
	b.Handle("/loan", gl.handleRecord)
	b.Handle("/predict", gl.handlePredict)
	b.Handle("/show", gl.handleShow)

	b.Handle(tb.OnText, gl.handleContext)

	b.Start()
}
