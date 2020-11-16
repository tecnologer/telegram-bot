package models

// LoginUrl This object represents a parameter of the inline keyboard button used to automatically authorize a user.
// Serves as a great replacement for the Telegram Login Widget when the user is coming from Telegram.
// All the user needs to do is tap/click a button and confirm that they want to log in:
// Telegram apps support these buttons as of version 5.7.
// https://core.telegram.org/file/811140015/1734/8VZFkwWXalM.97872/6127fa62d8a0bf2b3c
// Sample bot: @discussbot
type LoginUrl struct {
	URL         string `json:"url"`
	ForwardText string `json:"forward_text"`
	BotUsername string `json:"bot_username"`
}
