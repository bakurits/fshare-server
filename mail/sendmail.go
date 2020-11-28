package mail

import (
	gomail "gopkg.in/mail.v2"
)

const MailHost = "smtp.gmail.com"
const MAILPORT = 587

func SendMail(subject, body string, from string, password string, addr ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", addr...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(MailHost, MAILPORT, from, password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
