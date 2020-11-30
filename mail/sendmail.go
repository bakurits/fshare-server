package mail

import (
	gomail "gopkg.in/mail.v2"
)

const MAILHOST = "smtp.gmail.com"
const MAILPORT = 587

type Sender struct {
	Email    string
	Password string
}

func (s *Sender) SendMail(subject, body string, addr ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.Email)
	m.SetHeader("To", addr...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(MAILHOST, MAILPORT, s.Email, s.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
