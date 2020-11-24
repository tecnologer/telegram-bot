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

var webhookChannel chan *models.Update

func (b *Bot) setWebHook(webHookURL string, port int) (chan *models.Update, error) {
	v := url.Values{}
	if webHookURL != "" {
		v.Add("url", webHookURL)
	}

	res, err := http.PostForm(b.getEndpoint("setWebhook"), v)
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

	go http.ListenAndServe(fmt.Sprintf(":%s", port), http.HandlerFunc(weebHookHandler))
	logrus.Infof("listening messages on port %d", port)
	return webhookChannel, nil

}

func weebHookHandler(res http.ResponseWriter, req *http.Request) {
	body := &models.Update{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		logrus.WithError(err).Error("webhook: could not decode request body")
		return
	}

	webhookChannel <- body
}
