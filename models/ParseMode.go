package models

//ParseMode The Bot API supports basic formatting for messages. You can use bold, italic, underlined and strikethrough text, as well as inline links and pre-formatted code in your bots' messages. Telegram clients will render them accordingly. You can use either markdown-style or HTML-style formatting.
type ParseMode string

const (
	//Markdown2 Markdown sintax version 2
	//https://core.telegram.org/bots/api#markdownv2-style
	Markdown2 ParseMode = "MarkdownV2"
	//HTML format
	HTML ParseMode = "HTML"
	//Markdown legacy markdown
	Markdown ParseMode = "Markdown"
)
