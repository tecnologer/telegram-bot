package telegram

import (
	"context"
	"strings"

	"github.com/tecnologer/telegram-bot-api/models"
)

func init() {

}

var commands map[string]func(context.Context, *models.Update)

func validateCmd(ctx context.Context, update *models.Update) {
	cmd := getCmdFromMsg(update.Message)
	if cmd == "" {
		return
	}
	action, exists := commands[cmd]
	if exists {
		action(ctx, update)
	}
	// if action, exists := commands[cmd]; exists && action != nil {
	// 	action(ctx, update)
	// }
}

func getCmdFromMsg(msg *models.Message) string {
	msgParts := strings.Split(msg.Text, " ")
	if len(msgParts) == 0 {
		return ""
	}

	return msgParts[0]
}
