package slack

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jordi-reinsma/bagel/model"
)

const (
	urlConversationMembers = "https://slack.com/api/conversations.members"
	urlConversationOpen    = "https://slack.com/api/conversations.open"
	urlPostMessage         = "https://slack.com/api/chat.postMessage"
	conversationTemplate   = "👋 Bom dia! Tá na hora do encontro do <#%s>! 🥯 Combinem um horário para vocês tomarem um café e colocarem o papo em dia. 😊"
	channelTemplate        = "🥯 Bom dia! O bagel começou! Receberam seus convites para um café? 🥯"
)

var (
	slackToken      = os.Getenv("SLACK_TOKEN")
	slackChannelIDs = strings.Split(os.Getenv("SLACK_CHANNEL_IDS"), ",")
	slackBotID      = os.Getenv("SLACK_BOT_ID")
)

type Slacker interface {
	GetChannelUUIDs() []string
	GetUserUUIDs(string) ([]model.User, error)
	SendMessage(model.Match, string) error
	SendChannelMessage(string) error
}

type Client struct{}

func New(mock bool) Slacker {
	if mock {
		return Mock{}
	}
	return Client{}
}

func (c Client) GetChannelUUIDs() []string {
	return slackChannelIDs
}

func (c Client) GetUserUUIDs(channelID string) ([]model.User, error) {
	body := url.Values{}
	body.Set("channel", channelID)

	data, err := callAPI(urlConversationMembers, body)
	if err != nil {
		return nil, err
	}

	var users []model.User
	for _, user := range data["members"].([]interface{}) {
		if user.(string) == slackBotID {
			continue
		}
		users = append(users, model.User{
			UUID: user.(string),
		})
	}

	return users, nil
}

func (c Client) SendMessage(match model.Match, channelID string) error {
	body := url.Values{}
	body.Set("users", fmt.Sprintf("%s,%s", match.A.UUID, match.B.UUID))

	data, err := callAPI(urlConversationOpen, body)
	if err != nil {
		return err
	}

	conversationID := data["channel"].(map[string]interface{})["id"].(string)

	body = url.Values{}
	body.Set("channel", conversationID)
	body.Set("text", fmt.Sprintf(conversationTemplate, channelID))

	_, err = callAPI(urlPostMessage, body)
	return err
}

func (c Client) SendChannelMessage(channelID string) error {
	body := url.Values{}
	body.Set("channel", channelID)
	body.Set("text", channelTemplate)

	_, err := callAPI(urlPostMessage, body)
	return err
}

func callAPI(url string, body url.Values) (map[string]interface{}, error) {
	body.Set("token", slackToken)

	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if !data["ok"].(bool) {
		return nil, errors.New(data["error"].(string))
	}

	return data, nil
}
