package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tecnologer/telegram-bot-api/models"
)

const (
	urlBase string = "https://api.telegram.org/bot%s"
)

//UpdateConfig configuration for long polling updates
type UpdateConfig struct {
	Offset  int
	Limit   int
	Timeout int
}

var allMsg func(*models.Update)

//Bot base struct for telegram bot
type Bot struct {
	token           string
	shutdownChannel chan interface{}
	isWebhookSet    bool
	url             string
	commands        map[string]func(*models.Update)
}

//NewBot creates a new bot with especific token
func NewBot(token string) *Bot {
	return &Bot{
		token: token,
		url:   fmt.Sprintf(urlBase, token),
	}
}

func (b *Bot) getEndpoint(action string) string {
	separator := ""
	if !strings.HasPrefix(action, "/") {
		separator = "/"
	}
	return fmt.Sprintf("%s%s%s", b.url, separator, action)
}

// GetMe A simple method for testing your bot's auth token. Requires no parameters.
// Returns basic information about the bot in form of a User object.
func (b *Bot) GetMe() (*models.User, error) {
	res, err := http.Get(b.getEndpoint("getMe"))

	if err != nil {
		return nil, errors.Wrap(err, "Get met: error calling the end-point")
	}

	defer res.Body.Close()
	body, err := decodeUserBody(res)
	if err != nil {
		return nil, errors.Wrap(err, "Get met: error decoding body response")
	}

	if !body.Ok {
		return nil, fmt.Errorf("Get me: bot response with error: (%d), %s", body.ErrorCode, body.Description)
	}

	return body.Result, nil
}

//SendTextMessage send text message to the specific chat
func (b *Bot) SendTextMessage(chatID int, text string, replyTo int) error {
	return b.SendMessage(&models.SendMessage{
		ChatID:                   chatID,
		Text:                     text,
		ReplyToMessageID:         replyTo,
		ParseMode:                models.Markdown2,
		AllowSendingWithoutReply: true,
		ReplyMarkup:              "{}",
	})
}

//SendMessage send any kind of message to specific chat
func (b *Bot) SendMessage(message *models.SendMessage) error {
	body, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "parsing message to json body")
	}

	res, err := http.Post(b.getEndpoint("sendMessage"), "application/json", bytes.NewBuffer(body))

	if err != nil {
		return errors.Wrap(err, "sending message")
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("error seding message. Status code: %d", res.StatusCode)
	}

	return nil
}

// getUpdatesChan starts and returns a channel for getting updates.
func (b *Bot) getUpdatesChan(config UpdateConfig) (chan *models.Update, error) {
	ch := make(chan *models.Update)

	go func() {
		for {
			select {
			//TODO: Implement stop method
			case <-b.shutdownChannel:
				close(ch)
				return
			default:
			}

			updates, err := b.getUpdates(config)
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				if update.UpdateID >= config.Offset {
					config.Offset = update.UpdateID + 1
					ch <- update
				}
			}
		}
	}()

	return ch, nil
}

// getUpdates fetches updates.
// If a WebHook is set, this will not return any data!
//
// Offset, Limit, and Timeout are optional.
// To avoid stale items, set Offset to one higher than the previous item.
// Set Timeout to a large number to reduce requests so you can get updates
// instantly instead of having to wait between requests.
func (b *Bot) getUpdates(config UpdateConfig) ([]*models.Update, error) {
	v := url.Values{}
	if config.Offset != 0 {
		v.Add("offset", strconv.Itoa(config.Offset))
	}
	if config.Limit > 0 {
		v.Add("limit", strconv.Itoa(config.Limit))
	}
	if config.Timeout > 0 {
		v.Add("timeout", strconv.Itoa(config.Timeout))
	}

	resp, err := http.Post(b.getEndpoint("getUpdates"), "", nil)
	if err != nil {
		return nil, err
	}

	var updates []*models.Update
	var bodyRes *APIResponse
	err = json.NewDecoder(resp.Body).Decode(&bodyRes)
	json.Unmarshal(bodyRes.Result, &updates)

	return updates, nil
}

//SetCommand assign action to command starting with '/' (I.e: /start)
func (b *Bot) SetCommand(cmd string, callback func(*models.Update)) error {
	if callback == nil {
		return fmt.Errorf("callback function is required for a command")
	}

	if b.commands == nil {
		b.commands = make(map[string]func(*models.Update))
	}

	if !strings.HasPrefix(cmd, "/") {
		cmd = "/" + cmd
	}
	b.commands[cmd] = callback
	return nil
}

//Start starts the bot with long polling
func (b *Bot) Start() {
	if b.isWebhookSet {
		logrus.Error("starting with long polling: weebhook is configured")
		return
	}

	chanUpdates, err := b.getUpdatesChan(UpdateConfig{})
	if err != nil {
		logrus.WithError(err).Error("register for updates")
		return
	}
	for update := range chanUpdates {
		b.validateCmd(update)

		if allMsg != nil {
			allMsg(update)
		}
	}
}

//StartWithWebhook starts the bot using webhook
func (b *Bot) StartWithWebhook(url string, port int) {
	chanUpdates, err := b.setWebHook(url, port)
	if err != nil {
		logrus.WithError(err).Error("register for updates")
		return
	}
	for update := range chanUpdates {
		b.validateCmd(update)

		if allMsg != nil {
			allMsg(update)
		}
	}
}

//AllMessage redirect all messages to the function (fn)
func (b *Bot) AllMessage(fn func(*models.Update)) {
	allMsg = fn
}
