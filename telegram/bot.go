package telegram

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/tecnologer/telegram-bot-api/models"
)

var (
	urlBase string = "https://api.telegram.org/bot%s"
)

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

type UpdateConfig struct {
	Offset  int
	Limit   int
	Timeout int
}
type UpdatesChannel struct{}

// GetUpdatesChan starts and returns a channel for getting updates.
func (bot *Bot) GetUpdatesChan(config UpdateConfig) (chan *models.Update, error) {
	ch := make(chan *models.Update)

	go func() {
		for {
			select {
			//TODO: Implement stop method
			case <-bot.shutdownChannel:
				close(ch)
				return
			default:
			}

			updates, err := bot.GetUpdates(config)
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

// GetUpdates fetches updates.
// If a WebHook is set, this will not return any data!
//
// Offset, Limit, and Timeout are optional.
// To avoid stale items, set Offset to one higher than the previous item.
// Set Timeout to a large number to reduce requests so you can get updates
// instantly instead of having to wait between requests.
func (bot *Bot) GetUpdates(config UpdateConfig) ([]*models.Update, error) {
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
