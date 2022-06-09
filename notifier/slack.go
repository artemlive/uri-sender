package notifier

import (
	"fmt"
	"github.com/artemlive/uri-sender/config"
	"github.com/rs/zerolog/log"
)
import "github.com/slack-go/slack"


type Slack struct {
	authToken	string
	message    Message
	recipients []string
}


func NewSlackNotifier(message Message, recipients []string, conf config.Config) (*Slack, error) {
	if len(recipients) == 0 {
		return nil, fmt.Errorf("recipients can't be empty")
	}
	token, err := conf.GetSlackApiToken()
	if err != nil {
		return nil, err
	}
	return &Slack{authToken: token, message: message, recipients: recipients}, nil
}
func (s *Slack) Send() error {
	fmt.Println(s.message)
	err := s.sendMessage()
	return err
}

func (s *Slack) sendMessage() error {
	api := slack.New(s.authToken)
	if len(s.message.File) > 0 {
		params := slack.FileUploadParameters{
			Title:    s.message.Title,
			Filetype: "image/png",
			File:     s.message.File,
			Channels: s.recipients,
		}
		file, err := api.UploadFile(params)
		if err != nil {
			return fmt.Errorf("file \"%s\" upload: %s\n", file.Name, err)
		}
	}
	for _, recipient := range s.recipients {
		channelID, timestamp, err := api.PostMessage(
			recipient,
			slack.MsgOptionText(s.message.Message, false),
			slack.MsgOptionAsUser(true), // Add this if you want that the bot would post message as a user, otherwise it will send response using the default slackbot
		)
		if err != nil {
			return err
		}
		log.Info().Msgf("Message successfully sent to channel %s at %s", channelID, timestamp)
	}
	return nil
}