package shared

import (
	"bytes"
	// "fmt"
	g "github.com/nareshganesan/services/globals"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/smtp"
	"strings"
)

// Email entity
type Email struct {
	From     string
	To       []string
	Subject  string
	Body     string
	Password string
}

const (
	// MIME type for email entity
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

// ComposeEmail used to compose new email entity
func ComposeEmail(to []string, subject string) *Email {
	l := g.Gbl.Log
	l.WithFields(logrus.Fields{
		"to":      to,
		"subject": subject,
	}).Info("creating email object")
	return &Email{
		To:      to,
		Subject: subject,
	}
}

// parseTemplate helper to parse html email template file given filepath
func (e *Email) parseTemplate(fileName string, data interface{}) error {
	l := g.Gbl.Log
	l.Info("Parsing Email template")
	t, err := template.ParseFiles(fileName)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error":    err,
			"fileName": fileName,
			"data":     data,
		}).Error("Error parsing email template")
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		l.WithFields(logrus.Fields{
			"error":    err,
			"fileName": fileName,
			"data":     data,
		}).Error("Error binding email template")
		return err
	}
	e.Body = buffer.String()
	return nil
}

// SendHTMLEmail used for sending html email from email entity
func (e *Email) SendHTMLEmail() bool {
	l := g.Gbl.Log
	body := "From: " + e.From + "\n" + "To: " + strings.Join(e.To, ",") + "\r\nSubject: " + e.Subject + "\r\n" + MIME + "\r\n" + e.Body
	SMTP := g.Config.SMTP.Host + ":" + g.Config.SMTP.Port
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", e.From, e.Password, g.Config.SMTP.Host), e.From, e.To, []byte(body)); err != nil {
		l.WithFields(logrus.Fields{
			"error": err,
			"from":  e.From,
			"to":    e.To,
		}).Error("Email could not be sent!")
		return false
	}
	l.WithFields(logrus.Fields{
		"from": e.From,
		"to":   e.To,
	}).Info("Email sent!")
	return true
}

// Send used for sending plain/text email from email entity
func (e *Email) Send(templateName string, data interface{}) bool {
	l := g.Gbl.Log
	err := e.parseTemplate(templateName, data)
	if err != nil {
		l.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error sending email!")
		return false
	}
	if ok := e.SendHTMLEmail(); ok {
		l.WithFields(logrus.Fields{
			"to": e.To,
		}).Info("Email has been sent")
		return true
	}
	l.WithFields(logrus.Fields{
		"to": e.To,
	}).Error("Failed to send email")
	return false
}
