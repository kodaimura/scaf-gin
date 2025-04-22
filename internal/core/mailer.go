package core

type MailerI interface {
	SendText(to []string, subject, body string) error
	SendHTML(to []string, subject, body string) error
}

var Mailer MailerI = &noopMailer{}

func SetMailer(m MailerI) {
	Mailer = m
}

type noopMailer struct{}

func (n *noopMailer) SendText(to []string, subject, body string) error { return nil }
func (n *noopMailer) SendHTML(to []string, subject, body string) error { return nil }
