package slack

import "github.com/jordi-reinsma/bagel/model"

type Mock struct{}

func (m Mock) GetChannelUUIDs() []string {
	return []string{"general"}
}

func (m Mock) GetUserUUIDs(channelID string) ([]model.User, error) {
	// comment and uncomment some users to test different executions
	users := []model.User{
		{UUID: "111"},
		{UUID: "222"},
		{UUID: "333"},
		{UUID: "444"},
		{UUID: "555"},
		{UUID: "666"},
		{UUID: "777"},
		{UUID: "888"},
		{UUID: "999"},
		{UUID: "000"},
		{UUID: "XXX"},
		{UUID: "LLL"},
		{UUID: "PPP"},
	}

	return users, nil
}

func (m Mock) SendMessage(match model.Match, channelID string) error {
	return nil
}

func (m Mock) SendChannelMessage(channelID string) error {
	return nil
}
