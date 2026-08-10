package core

type MailerI interface {
	SendText(to []string, subject, body string) error
	SendHTML(to []string, subject, body string) error
}
