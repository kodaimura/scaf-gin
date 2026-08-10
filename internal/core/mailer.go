package core

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net/smtp"
	"strings"

	"scaf-gin/config"
)

type MailerI interface {
	SendText(to []string, subject, body string) error
	SendHTML(to []string, subject, body string) error
}

// SmtpMailer implements MailerI using SMTP. MailHog also uses SMTP without auth.
type SmtpMailer struct {
	from     string
	host     string
	port     string
	username string
	password string
}

func NewMailer() MailerI {
	switch strings.ToLower(config.MailProvider) {
	case "smtp":
		return NewSmtpMailer()
	case "mailhog":
		return NewMailHogMailer()
	default:
		return NewMailHogMailer()
	}
}

func NewMailHogMailer() MailerI {
	return &SmtpMailer{
		from: config.MailFrom,
		host: "mailhog",
		port: "1025",
	}
}

func NewSmtpMailer() MailerI {
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
	return s.send(to, msg)
}

func (s *SmtpMailer) SendHTML(to []string, subject, body string) error {
	msg := s.composeMessage(to, subject, "text/html", body)
	return s.send(to, msg)
}

func (s *SmtpMailer) address() string {
	return fmt.Sprintf("%s:%s", s.host, s.port)
}

func (s *SmtpMailer) auth() smtp.Auth {
	if s.username == "" || s.password == "" {
		return nil
	}
	return smtp.PlainAuth("", s.username, s.password, s.host)
}

func (s *SmtpMailer) send(to []string, msg []byte) error {
	if config.MailProvider != "smtp" {
		return smtp.SendMail(s.address(), nil, s.from, to, msg)
	}
	return s.sendSMTP(to, msg)
}

func (s *SmtpMailer) sendSMTP(to []string, msg []byte) error {
	c, err := smtp.Dial(s.address())
	if err != nil {
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12}); err != nil {
			return err
		}
	}
	if auth := s.auth(); auth != nil {
		if err := c.Auth(auth); err != nil {
			return err
		}
	}
	if err := c.Mail(s.from); err != nil {
		return err
	}
	for _, recipient := range to {
		if err := c.Rcpt(recipient); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.Close()
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
