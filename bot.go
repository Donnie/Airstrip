package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (st *State) startBot() {
	st.Bot.SetWebhook(&tb.Webhook{
		Listen:   ":" + st.Env.PORT,
		Endpoint: &tb.WebhookEndpoint{PublicURL: st.Env.WEBHOOK + "/hook"},
	})

	st.Bot.Handle("/start", st.handleHelp)
	st.Bot.Handle("/help", st.handleHelp)
	st.Bot.Handle("/charge", st.handleRecord)
	st.Bot.Handle("/expense", st.handleRecord)
	st.Bot.Handle("/delete", st.handleDelete)
	st.Bot.Handle("/gain", st.handleRecord)
	st.Bot.Handle("/income", st.handleRecord)
	st.Bot.Handle("/lend", st.handleRecord)
	st.Bot.Handle("/loan", st.handleRecord)
	st.Bot.Handle("/predict", st.handlePredict)
	st.Bot.Handle("/view", st.handleView)

	st.Bot.Handle(tb.OnText, st.handleText)
	st.Bot.Handle(tb.OnCallback, st.handleCallback)
}

func (st *State) handleHook() {
	http.HandleFunc("/hook", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var inp tb.Update
		err = json.Unmarshal(b, &inp)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		st.Bot.ProcessUpdate(inp)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(":"+st.Env.PORT, nil)
}
