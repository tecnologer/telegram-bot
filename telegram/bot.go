package telegram

import (
	"bytes"
	"context"
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

var (
	urlBase string = "https://api.telegram.org/bot%s"
)

var allMsg func(context.Context, *models.Update)

type Bot struct {
	token           string
	shutdownChannel chan interface{}
}

func NewBot(token string) *Bot {
	urlBase = fmt.Sprintf(urlBase, token)
	return &Bot{
		token: token,
	}
}

// GetMe A simple method for testing your bot's auth token. Requires no parameters.
// Returns basic information about the bot in form of a User object.
func (b *Bot) GetMe() (*models.User, error) {
	endpoint := urlBase + "/getMe"

	res, err := http.Get(endpoint)

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

func (b *Bot) SendTextMessage(chatID int, text string, replyTo int, parseMode models.ParseMode) error {
	return b.SendMessage(&models.SendMessage{
		ChatID:                   chatID,
		Text:                     text,
		ReplyToMessageID:         replyTo,
		ParseMode:                parseMode,
		AllowSendingWithoutReply: true,
		ReplyMarkup:              "{}",
	})
}

func (b *Bot) SendMessage(message *models.SendMessage) error {
	endpoint := urlBase + "/sendMessage"

	body, err := json.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "parsing message to json body")
	}

	res, err := http.Post(endpoint, "application/json", bytes.NewBuffer(body))

	if err != nil {
		return errors.Wrap(err, "sending message")
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("error seding message. Status code: %d", res.StatusCode)
	}

	return nil
}

type UpdateConfig struct {
	Offset  int
	Limit   int
	Timeout int
}
type UpdatesChannel struct{}

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
	endpoint := urlBase + "/getUpdates"
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

	resp, err := http.Post(endpoint, "", nil)
	if err != nil {
		return nil, err
	}

	var updates []*models.Update
	var bodyRes *APIResponse
	err = json.NewDecoder(resp.Body).Decode(&bodyRes)
	json.Unmarshal(bodyRes.Result, &updates)

	return updates, nil
}

//SetCommand set action to a command starting with /
func (b *Bot) SetCommand(cmd string, callback func(context.Context, *models.Update)) {
	if callback == nil {
		logrus.Info("callback function is required for a command")
		return
	}

	if commands == nil {
		commands = make(map[string]func(context.Context, *models.Update))
	}

	if !strings.HasPrefix(cmd, "/") {
		cmd = "/" + cmd
	}
	commands[cmd] = callback
}

func (b *Bot) Start(ctx context.Context) {
	chanUpdates, err := b.getUpdatesChan(UpdateConfig{})
	if err != nil {
		logrus.WithError(err).Error("register for updates")
		return
	}
	for update := range chanUpdates {
		validateCmd(ctx, update)

		if allMsg != nil {
			allMsg(ctx, update)
		}
	}
}

func (b *Bot) AllMessage(fn func(context.Context, *models.Update)) {
	allMsg = fn
}
