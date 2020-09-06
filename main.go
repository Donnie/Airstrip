package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

	bot, err := tb.NewBot(tb.Settings{Token: teleToken, Synchronous: true})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dsn := "user=airstrip password=postgres dbname=airstrip port=5432 sslmode=disable host=postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database")
		os.Exit(1)
	}
	// db.AutoMigrate(&Account{}, &Convo{}, &Record{})

	gl := Global{Bot: bot, Orm: db}

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

	bot.Handle(tb.OnText, gl.handleText)
	bot.Handle(tb.OnCallback, gl.handleCallback)

	http.HandleFunc("/hook", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var inp tb.Update
		err = json.Unmarshal(b, &inp)
		fmt.Println(string(b))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		bot.ProcessUpdate(inp)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(":"+port, nil)
}
