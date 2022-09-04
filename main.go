package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

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

	orm := initSQLDB()

	bot, err := tb.NewBot(tb.Settings{
		Token:  teleToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	st := State{
		Bot: bot,
		Orm: orm,
		Env: &Env{
			TELETOKEN: teleToken,
		},
	}

	st.startBot()
}
