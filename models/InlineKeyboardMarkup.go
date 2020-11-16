package models

//InlineKeyboardMarkup This object represents an inline keyboard that appears right next to the message it belongs to.
type InlineKeyboardMarkup struct {
	InlineKeyboard []*InlineKeyboardButton `json:"inline_keyboard"`
}
