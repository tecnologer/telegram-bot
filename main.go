package main

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tecnologer/go-secrets"
	"github.com/tecnologer/go-secrets/config"
	"github.com/tecnologer/telegram-bot-api/models"
	"github.com/tecnologer/telegram-bot-api/telegram"
)

var bot *telegram.Bot

func main() {
	secrets.InitWithConfig(&config.Config{})
	secTelegram, err := secrets.GetGroup("telegram")
	if err != nil {
		logrus.WithError(err).Error("token not configured in secrets")
	}

	bot = telegram.NewBot(secTelegram.GetString("token"))

	data, err := bot.GetMe()
	if err != nil {
		logrus.WithError(err).Error("fale ferga la fida")
		return
	}
	fmt.Printf("Mi nombre es %s y soy el bot @%s\n", data.FirstName, data.Username)

	bot.SetCommand("hola", sayHi)
	// bot.AllMessage(catchAll)

	// bot.Start(context.Background())
	bot.StartWithWebhook("https://cc1a2dc0536a.ngrok.io", nil)
	logrus.Info("program end")
}

func sayHi(ctx context.Context, update *models.Update) {
	bot.SendTextMessage(update.Message.Chat.ID, "Hola", 0, models.Markdown2)
}

func catchAll(ctx context.Context, update *models.Update) {
	fmt.Println(update.Message.Text)
}
