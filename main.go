package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tecnologer/telegram-bot-api/models"
	"github.com/tecnologer/telegram-bot-api/telegram"
)

const token = "<bot-token>"

var bot *telegram.Bot

func main() {

	bot = telegram.NewBot(token)

	data, err := bot.GetMe()
	if err != nil {
		logrus.WithError(err).Error("fale ferga la fida")
		return
	}
	fmt.Printf("Mi nombre es %s y soy el bot @%s\n", data.FirstName, data.Username)

	bot.SetCommand("hola", sayHi)
	// bot.AllMessage(catchAll)

	bot.Start(nil)
	logrus.Info("program end")
}

func sayHi(ctx context.Context, update *models.Update) {
	bot.SendTextMessage(update.Message.Chat.ID, "Hola", update.Message.MessageID, models.Markdown2)
}

func catchAll(ctx context.Context, update *models.Update) {
	fmt.Println(update.Message.Text)
}
