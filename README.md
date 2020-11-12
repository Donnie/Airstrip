# Airstrip
[![Go Report Card](https://goreportcard.com/badge/github.com/Donnie/Airstrip)](https://goreportcard.com/report/github.com/Donnie/Airstrip) [![Build Status](https://api.travis-ci.org/Donnie/Airstrip.svg?branch=master&status=passed)](https://travis-ci.org/github/Donnie/Airstrip) [![Maintainability](https://api.codeclimate.com/v1/badges/80f939bc59e3affb38ff/maintainability)](https://codeclimate.com/github/Donnie/Airstrip/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/80f939bc59e3affb38ff/test_coverage)](https://codeclimate.com/github/Donnie/Airstrip/test_coverage)

Personal finance management bot on Telegram

## Dev Setup

### Start Project
Add your Telegram bot token to the .env.local file and then

```make builddev```

```make dev```

## Prod Setup

### Configure your bot .env file
Add your Telegram bot token and Postgres details to the .env file and then

### Make Live
Update Make file `live` command with your server details and do

```make live```

## DB
### Migrate
Put your SQL in `airstrip.sql` and then

```make migrate```

### Dump
Get a SQL dump in `airstrip.sql` file

```make dump```

### Postgres Terminal
Tinker with the database

```make sql```

## Features
/record - Record an expense or gain

/recur - Record an income or a charge

/cancel - Cancel an ongoing record or recur process

/delete - Delete any record or recur

/predict Jan 2025 - Get a prediction of your financial standing

/view Jan 2025 - Get a list of records pertaining to the month

/stand Account - Get a current standing of any account

/stand Account Mar 2022 - Get a month wise total effect on any account

/help - To see this list
