package main

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tecnologer/go-secrets"
	"github.com/tecnologer/go-secrets/config"
	"github.com/tecnologer/telegram-bot-api/models"
	"github.com/tecnologer/telegram-bot-api/telegram"
)

var (
	bot  *telegram.Bot
	bot2 *telegram.Bot
)

var actions map[*regexp.Regexp]func(*models.Update)

func main() {
	secrets.InitWithConfig(&config.Config{})
	secTelegram, err := secrets.GetGroup("telegram")
	if err != nil {
		logrus.WithError(err).Error("token not configured in secrets")
	}

	bot = telegram.NewBot(secTelegram.GetString("token"))
	bot2 = telegram.NewBot(secTelegram.GetString("token2"))

	// actions = make(map[*regexp.Regexp]func(*models.Update))
	// addAction("^como te llamas\\?$", sayName)S
	bot.SetCommand("hola", sayHi)
	bot2.SetCommand("hola", sayHi)
	bot.AllMessage(catchAll)

	bot.StartWithWebhook("https://dc88553f0f65.ngrok.io", 8088)
	bot2.StartWithWebhook("https://dc88553f0f65.ngrok.io", 8089)
	logrus.Info("program end")
}

func sayHi(update *models.Update) {
	bot.SendTextMessage(update.Message.Chat.ID, "Hola", 0)
}

func catchAll(update *models.Update) {
	msg := update.GetMessage()

	fmt.Println(msg.Text)
	if msg.Text == "" {
		return
	}

	for rgx, fn := range actions {
		if rgx.Match([]byte(msg.Text)) {
			fn(update)
			break
		}
	}

	if update.IsEdited() {
		bot.SendTextMessage(update.GetChatID(), "Aunque lo modifiques, ya lo lei", 0)
	}
}

func addAction(expr string, fn func(*models.Update)) error {
	rgx, err := regexp.Compile(expr)
	if err != nil {
		return errors.Wrap(err, "add action")
	}

	actions[rgx] = fn
	return nil
}

func sayName(update *models.Update) {
	botInstance, err := bot.GetMe()
	if err != nil {
		logrus.WithError(err).Error("say name error")
		return
	}

	bot.SendTextMessage(update.GetChatID(), fmt.Sprintf("Hola, mi nombre es %s", botInstance.FirstName), update.GetMessageID())
}
