package models

//SendMessage struct to send text messages
type SendMessage struct {
	ChatID                   int              `json:"chat_id"`
	Text                     string           `json:"text"`
	ParseMode                ParseMode        `json:"parse_mode"`
	Entities                 []*MessageEntity `json:"entities"`
	DisableWebPagePreview    bool             `json:"disable_web_page_preview"`
	DisableNotification      bool             `json:"disable_notification"`
	ReplyToMessageID         int              `json:"reply_to_message_id"`
	AllowSendingWithoutReply bool             `json:"allow_sending_without_reply"`
	ReplyMarkup              interface{}      `json:"reply_markup"`
}
