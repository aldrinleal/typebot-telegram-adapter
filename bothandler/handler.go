package bothandler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jmespath/go-jmespath"
	"github.com/joomcode/errorx"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)
import "github.com/shurcooL/go-goon"

type Handler struct {
	telegramApiToken  string
	typebotApiToken   string
	typebotId         string
	typebotApiBaseUrl string
	client            *resty.Client
	sessionMap        map[int64]string
}

type TypeBotMessage struct {
	Id      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Content struct {
		Type     string        `json:"type,omitempty"`
		RichText []interface{} `json:"richText,omitempty"`
		Markdown []interface{} `json:"markdown,omitempty"`
	} `json:"content,omitempty"`
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
		sessionMap:        make(map[int64]string),
	}

	return handler, nil
}

func (h *Handler) HandlerFunc(ctx context.Context, b *bot.Bot, update *models.Update) {
	log.Infof("Got update: %s", goon.Sdump(update))

	if textResult, err := jmespath.Search("message.text", update); nil == err {
		typeBotApiResponse := make(map[string]interface{})

		if text, ok := textResult.(string); ok {
			fromId := update.Message.From.ID
			sessionId, err := h.LookupSessionId(fromId)

			if ENotFound == err {
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
					err := h.RegisterSession(fromId, sessionId)

					if nil != err {
						log.Warnf("Oops: %s", errorx.Decorate(err, "registering session Id (%s / %d)", sessionId, fromId))
					}
				}
			} else if nil != err {
				log.Warnf("Oops: %s", err)
			} else if nil == err {
				// Valid Session Id, valid Stuff

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

		if messages, err := jmespath.Search("messages", typeBotApiResponse); nil == err {
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

	//  "github.com/jmespath/go-jmespath"
}

func messageToSendMessage(update *models.Update, message TypeBotMessage) *bot.SendMessageParams {
	reply := &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
	}

	if message.Content.Type == "richText" {
		reply.Text = getRichTextFor(message)
	}

	return reply
}

func getRichTextFor(message TypeBotMessage) string {
	resultArray := []string{}

	if 1 != len(message.Content.RichText) {
		log.Warnf("Unexpected length for RichText: %d", len(message.Content.RichText))

		return ""
	}

	children, _ := jmespath.Search("[0].children | []", message.Content.RichText)

	childrenArr := (children).([]interface{})

	for _, v := range childrenArr {
		if text, err := jmespath.Search("text", v); err == nil && text != nil {
			resultArray = append(resultArray, text.(string))
		} else if childTextArr, err := jmespath.Search("children[*].children[*].text | []", v); nil == err {
			childTextArray := childTextArr.([]interface{})

			for _, v2 := range childTextArray {
				resultArray = append(resultArray, v2.(string))
			}
		} else {
			log.Warnf("Unexpected Node: %s", goon.Sdump(v))
		}
	}

	return strings.Join(resultArray, "")
}

func typeBotMessageFrom(v interface{}) TypeBotMessage {
	byteArray, _ := json.Marshal(v)

	result := &TypeBotMessage{}

	json.Unmarshal(byteArray, result)

	return *result
}

var ENotFound = errors.New("sessionId not found.")

func (h *Handler) DeleteSessionId(id int64) error {
	delete(h.sessionMap, id)

	return nil
}

func (h *Handler) LookupSessionId(id int64) (string, error) {
	if sessionId, found := h.sessionMap[id]; found {
		return sessionId, nil
	}

	return "", ENotFound
}

func (h *Handler) RegisterSession(id int64, sessionId string) error {
	h.sessionMap[id] = sessionId

	return nil
}
