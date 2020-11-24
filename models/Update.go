package models

//Update This object represents an incoming update.
type Update struct {
	UpdateID           int                 `json:"update_id"`
	Message            *Message            `json:"message"`
	EditedMessage      *Message            `json:"edited_message"`
	ChannelPost        *Message            `json:"channel_post"`
	EditedChannelPost  *Message            `json:"edited_channel_post"`
	InlineQuery        *InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *CallbackQuery      `json:"callback_query"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query"`
	Poll               *Poll               `json:"poll"`
	PollAnswer         *PollAnswer         `json:"poll_answer"`
}

//GetChatID returns the chat from come the update
func (u *Update) GetChatID() int {
	msg := u.GetMessage()
	if msg == nil {
		return 0
	}

	return msg.Chat.ID
}

//GetMessage returns the message or edited message
func (u *Update) GetMessage() *Message {
	if u.EditedMessage != nil {
		return u.EditedMessage
	}

	return u.Message
}

//GetMessageID returns the id of message or edited message
func (u *Update) GetMessageID() int {
	if u.EditedMessage != nil {
		return u.EditedMessage.MessageID
	}

	return u.Message.MessageID
}

//IsEdited returns the if the update is a edited message
func (u *Update) IsEdited() bool {
	return u.EditedMessage != nil
}
