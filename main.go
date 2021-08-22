package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type wallet map[string]float64 // key = currency code, value = currency amount

var db = map[int]wallet{} // key = userID

type bResponce struct { // API responce
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

const natCurName = "RUB" // national currency name (may be "USDT",..)
const natCurSymbol = "₽" // national currency symbol (may be "$",..)

func getToken() string {
	return "1951162952:AAGZGFKqvl0g46PA85wID5jpxTsBvmHaKYQ"
}

func getRate(cur string) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s%s", cur, natCurName)
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	var bRes bResponce
	err = json.NewDecoder(res.Body).Decode(&bRes)
	if err != nil {
		return 0, err
	}
	if bRes.Symbol == "" {
		return 0, errors.New("currency code is wrong")
	}
	return bRes.Price, nil
}

func main() {
	bot, err := tgbotapi.NewBotAPI(getToken())
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u) // new income messages and errors
	var cur string                        // currency code
	var sval string                       // amount in string format
	var fval float64                      // amount in numeric format
	for update := range updates {         // new income message
		if update.Message == nil {
			continue
		}
		command := strings.Split(update.Message.Text, " ")
		commandCode := strings.ToUpper(command[0])
		userId := update.Message.From.ID
		switch commandCode {
		case "ADD", "SUB":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "command format with codes ADD and SUB is: ADD/SUB <currencyCode> <amount>"))
				continue
			}
			cur = strings.ToUpper(command[1]) // currency code
			sval = command[2]                 // amount (string format)
			_, err := getRate(cur)            // get currency rate
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "currency code is wrong"))
				continue
			}
			fval, err = strconv.ParseFloat(sval, 64) // amount (numeric format)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
				continue
			}

			if _, ok := db[userId]; !ok {
				db[userId] = make(wallet) // 1st command from the new user is "ADD"
			}
			switch commandCode {
			case "ADD":
				db[userId][cur] += fval
			case "SUB":
				var newBalance float64
				if newBalance = db[userId][cur] - fval; newBalance < 0 {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("balance %.02f is insufficient", db[userId][cur])))
					continue
				}
				db[userId][cur] = newBalance
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "done"))
		case "DEL":
			if len(command) != 2 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "command format with code DEL is: DEL <currencyCode>"))
				continue
			}
			if _, ok := db[userId]; !ok {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You have no amout"))
				continue
			}
			cur = strings.ToUpper(command[1])
			if db[userId][cur] > 0 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("can't delete currency: balance %.02f > 0", db[userId][cur])))
				continue
			}
			delete(db[userId], cur)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "done"))
		case "SHOW":
			if _, ok := db[userId]; !ok {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You have no amout"))
				continue
			}
			//res = ""
			for cur, fval := range db[userId] {
				rate, err := getRate(cur)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
					continue
				}
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s: %.2f %s%.02f (rate %.02f)\n", cur, fval, natCurSymbol, fval*rate, rate)))
			}
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "command code: "+command[0]+" is wrong"))
			continue
		}
	}
}
