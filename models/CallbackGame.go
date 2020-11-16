package models

//CallbackGame A placeholder, currently holds no information. Use BotFather to set up your game.
type CallbackGame struct {
	UserID             int    `json:"user_id"`
	Score              int    `json:"score"`
	Force              bool   `json:"force"`
	DisableEditMessage bool   `json:"disable_edit_message"`
	ChatID             int    `json:"chat_id"`
	MessageID          int    `json:"message_id"`
	InlineMessageID    string `json:"inline_message_id"`
}
