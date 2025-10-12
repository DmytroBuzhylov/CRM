package email_sender

import (
	"bytes"
	"embed"
	"gopkg.in/gomail.v2"
	"text/template"
)

type MailService struct {
	dialer *gomail.Dialer
	from   string
}

func NewMailService(host string, port int, username, password, from string) *MailService {
	return &MailService{
		dialer: gomail.NewDialer(host, port, username, password),
		from:   from,
	}
}

//go:embed email_templates/*.html
var emailTemplatesFS embed.FS

func (ms *MailService) Send(to string, subject string, templateName string, data any) error {
	body, err := renderTemplate(templateName, data)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", ms.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return ms.dialer.DialAndSend(m)
}

func renderTemplate(templateName string, data any) (string, error) {
	tmpl, err := template.ParseFS(emailTemplatesFS, "email_templates/"+templateName)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
