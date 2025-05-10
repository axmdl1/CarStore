package email

import "log"

type Sender interface {
	Send(to, subject, body string) error
}

type ConsoleSender struct{}

func NewConsoleSender() *ConsoleSender {
	return &ConsoleSender{}
}

func (s *ConsoleSender) Send(to, subject, body string) error {
	log.Printf("[Email] To: %s | Subject: %s | Body: %s", to, subject, body)
	return nil
}
