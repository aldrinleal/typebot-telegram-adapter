package bothandler

import (
	"context"
	"github.com/aldrinleal/typebot-telegram-adapter/session"
	"github.com/go-resty/resty/v2"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jmespath/go-jmespath"
	"github.com/joomcode/errorx"
	"github.com/shurcooL/go-goon"
	log "github.com/sirupsen/logrus"
	"os"
)

type Handler struct {
	telegramApiToken  string
	typebotApiToken   string
	typebotId         string
	typebotApiBaseUrl string
	client            *resty.Client
	SessionManager    session.SessionManager
}

func NewHandler() (*Handler, error) {
	telegramApiToken := os.Getenv("TELEGRAM_APITOKEN")
	typebotApiToken := os.Getenv("TYPEBOT_API_TOKEN")
	typebotId := os.Getenv("TYPEBOT_ID")
	typebotApiBaseUrl := os.Getenv("TYPEBOT_API_BASE_URL")

	handler := &Handler{
		telegramApiToken:  telegramApiToken,
		typebotApiToken:   typebotApiToken,
		typebotId:         typebotId,
		typebotApiBaseUrl: typebotApiBaseUrl,
		client:            resty.New(),
		SessionManager:    session.NewInMemorySessionManager(),
	}

	return handler, nil
}

func (h *Handler) HandlerFunc(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Infof("Got update: %s", goon.Sdump(update))

	if textResult, err := jmespath.Search("message.text", update); nil == err {
		if text, ok := textResult.(string); ok {
			typeBotApiResponse := make(map[string]interface{})

			h.handleTextUpdate(update, text, typeBotApiResponse)

			if messages, ok := typeBotApiResponse["messages"]; ok {
				if nil != messages {
					if messageAsArray, ok := messages.([]interface{}); ok {
						for _, v := range messageAsArray {
							message := typeBotMessageFrom(v)

							b.SendMessage(ctx, messageToSendMessage(update, message))
						}
					}
				}
			}
		}
	}
}

func (h *Handler) handleTextUpdate(update *models.Update, text string, typeBotApiResponse map[string]interface{}) {
	fromId := update.Message.From.ID
	sessionId, err := h.SessionManager.LookupSession(fromId)

	if session.ENotFound == err {
		startChatUrl := h.typebotApiBaseUrl + "/v1/typebots/" + h.typebotId + "/startChat"

		resp, err := h.client.R().
			SetDebug(true).
			SetAuthToken(h.typebotApiToken).
			SetBody(map[string]map[string]interface{}{
				"message": {
					"type": "text",
					"text": text,
				},
			}).
			SetResult(&typeBotApiResponse).
			Post(startChatUrl)

		if nil != err {
			log.Warnf("Oops: %s", errorx.Decorate(err, "while opening chat to '%s' for userId %d", startChatUrl, fromId))
		}

		if resp.StatusCode() < 200 || resp.StatusCode() >= 400 {
			log.Warnf("Unexpexted statusCode: %03d", resp.StatusCode())
		}

		if sessionId, ok := typeBotApiResponse["sessionId"].(string); ok {
			err := h.SessionManager.RegisterSession(fromId, sessionId)

			if nil != err {
				log.Warnf("Oops: %s", errorx.Decorate(err, "registering session Id (%s / %d)", sessionId, fromId))
			}
		}
	} else if nil != err {
		log.Warnf("Oops: %s", err)
	} else if nil == err {
		// Valid Session Id, valid Stuff - just refresh

		h.SessionManager.RefreshSession(fromId)

		continueChatUrl := h.typebotApiBaseUrl + "/v1/sessions/" + sessionId + "/continueChat"

		resp, err := h.client.R().
			SetDebug(true).
			SetAuthToken(h.typebotApiToken).
			SetBody(map[string]map[string]interface{}{
				"message": {
					"type": "text",
					"text": text,
				},
			}).
			SetResult(&typeBotApiResponse).
			Post(continueChatUrl)

		if nil != err {
			log.Warnf("Oops: %s", errorx.Decorate(err, "while opening chat to '%s' for userId %d", continueChatUrl, fromId))
		}

		if resp.StatusCode() < 200 || resp.StatusCode() >= 400 {
			log.Warnf("Unexpexted statusCode: %03d", resp.StatusCode())
		}
	}
}
