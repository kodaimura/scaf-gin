package mailer

import (
	"fmt"
	"net/smtp"
	"strings"
	"mime"
	"goscaf/config"
	"goscaf/internal/core"
)

type SmtpMailer struct {
	from     string
	host     string
	port     string
	username string
	password string
}

func NewSmtpMailer() core.MailerI {
	return &SmtpMailer{
		from:     config.MailFrom,
		host:     config.SMTPHost,
		port:     config.SMTPPort,
		username: config.SMTPUser,
		password: config.SMTPPass,
	}
}

func (s *SmtpMailer) SendText(to []string, subject, body string) error {
	msg := s.composeMessage(to, subject, "text/plain", body)
	return smtp.SendMail(s.address(), s.auth(), s.from, to, msg)
}

func (s *SmtpMailer) SendHTML(to []string, subject, body string) error {
	msg := s.composeMessage(to, subject, "text/html", body)
	return smtp.SendMail(s.address(), s.auth(), s.from, to, msg)
}

func (s *SmtpMailer) address() string {
	return fmt.Sprintf("%s:%s", s.host, s.port)
}

func (s *SmtpMailer) auth() smtp.Auth {
	return smtp.PlainAuth("", s.username, s.password, s.host)
}

func (s *SmtpMailer) composeMessage(to []string, subject, contentType, body string) []byte {
	header := s.defaultHeader(to, subject)
	header += fmt.Sprintf("Content-Type: %s; charset=UTF-8\r\n", contentType)
	return []byte(header + "\r\n" + body)
}

func (s *SmtpMailer) defaultHeader(to []string, subject string) string {
	return "From: " + s.from + "\r\n" +
		"To: " + strings.Join(to, ", ") + "\r\n" +
		"Subject: " + mime.QEncoding.Encode("UTF-8", subject) + "\r\n" +
		"MIME-Version: 1.0\r\n"
}
