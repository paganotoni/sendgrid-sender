package sender

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	nmail "net/mail"

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

	mm, err := buildMail(m)
	if err != nil {
		return err
	}

	response, err := ps.client.Send(mm)
	if response.StatusCode != 202 {
		return fmt.Errorf("Error sending email, code %v body %v", response.StatusCode, response.Body)
	}

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

func buildMail(m mail.Message) (*smail.SGMailV3, error) {
	mm := new(smail.SGMailV3)
	from, err := nmail.ParseAddress(m.From)
	if err != nil {
		return &smail.SGMailV3{}, fmt.Errorf("invalid from (%s): %s", from, err.Error())
	}
	mm.SetFrom(smail.NewEmail(from.Name, from.Address))
	mm.Subject = m.Subject

	p := smail.NewPersonalization()
	for _, toEmail := range m.To {
		to, err := nmail.ParseAddress(toEmail)
		if err != nil {
			return &smail.SGMailV3{}, fmt.Errorf("invalid to (%s): %s", toEmail, err.Error())
		}
		p.AddTos(smail.NewEmail(to.Name, to.Address))
	}

	for k, v := range m.Data {
		p.SetCustomArg(k, v.(string))
	}

	html := smail.NewContent("text/html", m.Bodies[0].Content)
	text := smail.NewContent("text/plain", m.Bodies[1].Content)
	mm.AddPersonalizations(p)
	mm.AddContent(text, html)

	for _, a := range m.Attachments {
		b := new(bytes.Buffer)
		if n, err := b.ReadFrom(a.Reader); err != nil {
			return &smail.SGMailV3{}, fmt.Errorf("Error attaching file: n %v error %v", n, err)
		}
		disposition := "attachment"
		if a.Embedded {
			disposition = "inline"
		}

		attachment := smail.NewAttachment()
		attachment.SetFilename(a.Name)
		attachment.SetContentID(a.Name)
		attachment.SetContent(b.String())
		attachment.SetType(a.ContentType)
		attachment.SetDisposition(disposition)
		mm.AddAttachment(attachment)
	}

	return mm, nil
}
