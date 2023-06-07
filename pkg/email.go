package pkg

import (
	"bytes"
	"crypto/tls"
	"errors"
	"hitshop/internal/config"
	"hitshop/internal/core"
	"html/template"
	"os"
	"path/filepath"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	VerifyURL      string
	UnsubscribeURL string

	AccountEmail string
	Subject      string
}

// 👇 Email template parser

func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}

func SendEmail(acc *core.Account, accEmail, code string, config *config.Cfg) error {
	//* 👇 Send Email
	data := EmailData{
		VerifyURL:      config.ClientOrigin + "/verifyemail/" + code,
		UnsubscribeURL: config.ClientOrigin + "/unsubscribe/" + code,
		AccountEmail:   accEmail,
		Subject:        "Код проверки вашей учетной записи",
	}

	//* Sender data.
	from := config.EmailFrom
	smtpPass := config.SMTPPass
	smtpUser := config.SMTPUser
	to := acc.Email
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort
	var body bytes.Buffer

	template, err := ParseTemplateDir("templates")
	if err != nil {
		return errors.New("Could not parse template: " + err.Error())
	}

	err = template.ExecuteTemplate(&body, "verificationCode.html", &data)
	if err != nil {
		return errors.New("Could not execute template: " + err.Error())
	}

	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetHeader("List-Unsubscribe", data.UnsubscribeURL)

	m.SetBody("text/html", body.String())
	m.AddAlternative("text/html`", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return errors.New("Could not send email: " + err.Error())
	}
	return nil
}
