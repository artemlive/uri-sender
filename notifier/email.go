package notifier

type Email struct {
	message Message
}

func (e *Email) Send() error {
	return nil
}

func NewEmailNotifier(message Message, recipients []string) (*Email, error) {
	return &Email{message: message}, nil
}