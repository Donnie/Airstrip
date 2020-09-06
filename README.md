# Airstrip
[![Go Report Card](https://goreportcard.com/badge/github.com/Donnie/Airstrip)](https://goreportcard.com/report/github.com/Donnie/Airstrip) [![Build Status](https://api.travis-ci.org/Donnie/Airstrip.svg?branch=master&status=passed)](https://travis-ci.org/github/Donnie/Airstrip) [![Maintainability](https://api.codeclimate.com/v1/badges/80f939bc59e3affb38ff/maintainability)](https://codeclimate.com/github/Donnie/Airstrip/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/80f939bc59e3affb38ff/test_coverage)](https://codeclimate.com/github/Donnie/Airstrip/test_coverage)

Personal finance management bot on Telegram

## Dev Setup
### Webhook
The app needs to set up a webhook for Telegram to relay updates.

You can set up your webhook locally on your port using ngrok like so:

```./ngrok http 8080```

Copy the forwarding URL from ngrok to the .env.local file

### Start Project
Add your Telegram bot token to the .env.local file and then

```make builddev```

```make dev```

### Migrate
Copy `airstrip-sample.sql` to `airstrip.sql` and then

```make migrate```

## Prod Setup
Add your Telegram bot token, port and webhook to the .env file and then

```make build```

```make up```

```make migrate```

## Features
/expense Record an expense

/gain Record any receipt

/charge Record fixed costs like rent, etc.

/income Record an income source like Salary

/lend Lend money to someone

/loan Take a loan from someone

`/predict Jan 2025` Get a prediction for your financial standing

`/view Jan 2025` Get a list of records pertaining to the month
