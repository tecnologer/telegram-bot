package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tecnologer/telegram-bot-api/models"
)

type webhookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

var webhookChannel chan *models.Update

func (b *Bot) setWebHook(webHookURL string) (chan *models.Update, error) {
	endpoint := urlBase + "/setWebhook"

	v := url.Values{}
	if webHookURL != "" {
		v.Add("url", webHookURL)
	}

	res, err := http.PostForm(endpoint, v)
	if err != nil {
		return nil, errors.Wrap(err, "configuring the url of webhook")
	}

	resBody, err := decodeWebhookResponse(res)
	if err != nil {
		return nil, errors.Wrap(err, "error sending webhook to telegram")
	}

	if !resBody.OK {
		return nil, fmt.Errorf("telegram set webhook response: %s", resBody.Description)
	}

	b.isWebhookSet = true
	webhookChannel = make(chan *models.Update)
	go http.ListenAndServe(":8088", http.HandlerFunc(weebHookHandler))
	logrus.Info("listening messages on port 8088")
	return webhookChannel, nil

}

func weebHookHandler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	body := &models.Update{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		logrus.WithError(err).Error("webhook: could not decode request body")
		return
	}

	webhookChannel <- body
}
