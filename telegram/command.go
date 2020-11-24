package telegram

import (
	"strings"

	"github.com/tecnologer/telegram-bot-api/models"
)

func (b *Bot) validateCmd(update *models.Update) {
	cmd := getCmdFromMsg(update.Message)
	if cmd == "" {
		return
	}

	if action, exists := b.commands[cmd]; exists {
		action(update)
	}
}

func getCmdFromMsg(msg *models.Message) string {
	if msg == nil {
		return ""
	}

	msgParts := strings.Split(msg.Text, " ")
	if len(msgParts) == 0 {
		return ""
	}

	return msgParts[0]
}
