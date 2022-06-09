package notifier

import (
	"github.com/artemlive/uri-sender/config"
	"github.com/rs/zerolog/log"
)

type Action string

type Message struct {
	Title	string
	Template string
	Message string
	File	string
}

const (
	SLACK Action = "SLACK"
	EMAIL Action = "EMAIL"
)

type Creator interface {
	CreateNotifier(action Action, message Message, recipients []string) (*Notifier, error)
}

type Notifier interface {
	Send() error
}

type NotificatorCreator struct{}

func NewCreator() *NotificatorCreator {
	return &NotificatorCreator{}
}

func (c *NotificatorCreator) CreateNotifier(action Action, message Message, recipients []string, conf config.Config) (Notifier, error) {
	var notifier Notifier
	var err error
	switch action {
	case SLACK:
		notifier, err = NewSlackNotifier(message, recipients, conf)
		if err != nil {return nil, err}
	case EMAIL:
		notifier, err = NewEmailNotifier(message, recipients)
		if err != nil {return nil, err}
	default:
		log.Fatal().Msgf("Unknown Notification method \"%s\"", action)
	}
	return notifier, nil
}

