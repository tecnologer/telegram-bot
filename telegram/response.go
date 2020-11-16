package telegram

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tecnologer/telegram-bot-api/models"
)

type APIResponse struct {
	Ok          bool                `json:"ok"`
	Result      json.RawMessage     `json:"result"`
	ErrorCode   int                 `json:"error_code"`
	Description string              `json:"description"`
	Parameters  *ResponseParameters `json:"parameters"`
}

type ResponseParameters struct {
	MigrateToChatID int64 `json:"migrate_to_chat_id"` // optional
	RetryAfter      int   `json:"retry_after"`        // optional
}

type Response struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

type userResponse struct {
	*Response
	Result *models.User `json:"result"`
}

func decodeUserBody(res *http.Response) (body *userResponse, err error) {
	body = &userResponse{}
	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		return nil, errors.Wrap(err, "parsing response (json) to user")
	}

	return
}
