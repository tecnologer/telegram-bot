package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tecnologer/telegram-bot-api/telegram"
)

const token = "659416092:AAEMdKzed2MODndq599fcz9a02YhzvQBjoM"

func main() {
	bot := telegram.NewBot(token)

	data, err := bot.GetMe()
	if err != nil {
		logrus.WithError(err).Error("fale ferga la fida")
		return
	}
	fmt.Printf("Mi nombre es %s y soy el bot @%s\n", data.FirstName, data.Username)
	chanMessages, err := bot.GetUpdatesChan(telegram.UpdateConfig{})
	if err != nil {
		logrus.WithError(err).Error("fale ferga la fida")
		return
	}

	for msg := range chanMessages {
		fmt.Println(msg)
	}
}
