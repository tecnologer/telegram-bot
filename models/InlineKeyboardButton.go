package models

//InlineKeyboardButton This object represents one button of an inline keyboard. You must use exactly one of the optional fields.
type InlineKeyboardButton struct {
	Text                         string        `json:"text"`
	URL                          string        `json:"url"`
	LoginURL                     *LoginUrl     `json:"login_url"`
	CallbackData                 string        `json:"callback_data"`
	SwitchInlineQuery            string        `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat string        `json:"switch_inline_query_current_chat"`
	CallbackGame                 *CallbackGame `json:"callback_game"`
	Pay                          bool          `json:"pay"`
}
