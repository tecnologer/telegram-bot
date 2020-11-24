# Telegram Bot API

Connect with Telegram bot API

## Create bot

`bot := telegram.NewBot("<token>")`

## Configure command

`bot.SetCommand("[/]<cmd>", <commandAction>)`

## Catch all messages

`bot.AllMessage(<action>)`

## Use webhook

`bot.StartWithWebhook("<url>", <port>)`

## or Use long polling

`bot.Start()`

## Send Message

`bot.SendTextMessage(<chat-id>, <message>, <0|message-id>)`

## Example

```golang
func main(){
    // Create bot
    bot := telegram.NewBot("<token>")

    //Configure command
    bot.SetCommand("start", start)

    //All Messages
    bot.AllMessage(catchAll)

    //Use webhook
    bot.StartWithWebhook("http://mywebhook.com", 8080)
}

func start(update *models.Update){
    bot.SendTextMessage(update.GetChatID(), "Hello, the bot now is running!", 0)
}

func catchAll(update *models.Update) {
    fmt.Println("new message received")
}
```
