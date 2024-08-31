package mailer

import (
	"fmt"
	"time"

	"github.com/go-mail/mail"
)

// TODO: Templating

type Mailer struct {
	mail   *mail.Dialer
	sender string
}

func New(m *mail.Dialer, senderEmail string) *Mailer {
	return &Mailer{m, senderEmail}
}

func (m *Mailer) SendResetCode(targetEmail string, code string, validity time.Time) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", targetEmail)
	msg.SetHeader("Subject", "Forgot Password")
	msg.SetBody("text/html", fmt.Sprintf("Your reset code is: %v\nValid before: %v\n", code, validity.Format("2006-01-02 15:04:05 MST")))

	return m.mail.DialAndSend(msg)
}
