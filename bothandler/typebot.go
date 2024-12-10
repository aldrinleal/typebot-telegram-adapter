package bothandler

import (
	"encoding/json"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/jmespath/go-jmespath"
	"github.com/shurcooL/go-goon"
	log "github.com/sirupsen/logrus"
	"strings"
)

type TypeBotMessage struct {
	Id      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Content struct {
		Type     string        `json:"type,omitempty"`
		RichText []interface{} `json:"richText,omitempty"`
		Markdown []interface{} `json:"markdown,omitempty"`
	} `json:"content,omitempty"`
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
