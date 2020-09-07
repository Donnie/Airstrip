package main

import (
	"fmt"
	"os"

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

	port, exists := os.LookupEnv("PORT")
	if !exists {
		fmt.Println("Add PORT to .env file")
		os.Exit(1)
	}

	webhook, exists := os.LookupEnv("WEBHOOK")
	if !exists {
		fmt.Println("Add WEBHOOK to .env file")
		os.Exit(1)
	}

	dsn := "user=airstrip password=postgres dbname=airstrip port=5432 sslmode=disable host=postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database")
		os.Exit(1)
	}
	// db.AutoMigrate(&Account{}, &Convo{}, &Record{})

	bot, err := tb.NewBot(tb.Settings{Token: teleToken, Synchronous: true})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	st := State{
		Bot: bot,
		Orm: db,
		Env: &Env{
			PORT:      port,
			TELETOKEN: teleToken,
			WEBHOOK:   webhook,
		},
	}

	st.startBot()
	st.handleHook()
}
