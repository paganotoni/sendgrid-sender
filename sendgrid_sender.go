package sender

import (
	"errors"
	"os"

	"github.com/gobuffalo/buffalo/mail"
	sendgrid "github.com/sendgrid/sendgrid-go"
	smail "github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stanislas-m/mocksmtp"
)

//SendgridSender implements the Sender interface to be used
//within buffalo mailer generated package.
type SendgridSender struct {
	TestSender *mocksmtp.MockSMTP
	client     *sendgrid.Client
}

//Send sends an email to Sendgrid for delivery, it assumes
//bodies[0] is HTML body and bodies[1] is text.
func (ps SendgridSender) Send(m mail.Message) error {
	if len(m.Bodies) < 2 {
		return errors.New("you must specify at least 2 bodies HTML and plain text")
	}

	if os.Getenv("GO_ENV") == "test" {
		return ps.TestSender.Send(m)
	}

	mm := new(smail.SGMailV3)
	mm.SetFrom(smail.NewEmail("", m.From))
	mm.Subject = m.Subject

	p := smail.NewPersonalization()
	for _, to := range m.To {
		p.AddTos(smail.NewEmail("", to))
	}

	html := smail.NewContent("text/html", m.Bodies[0].Content)
	text := smail.NewContent("text/plain", m.Bodies[1].Content)
	mm.AddPersonalizations(p)
	mm.AddContent(text, html)

	_, err := ps.client.Send(mm)
	return err
}

// NewSendgridSender creates a new SendgridSender with
// its own Sendgrid client inside
func NewSendgridSender(APIKey string) SendgridSender {
	client := sendgrid.NewSendClient(APIKey)
	return SendgridSender{
		client:     client,
		TestSender: mocksmtp.New(),
	}
}
