package email_sender

type Sender interface {
	Send(to string, subject string, templateName string, data any) error
}
