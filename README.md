### Sendgrid Buffalo Sender

This is a [buffalo](github.com/gobuffalo/buffalo) sender for the [Sendgrid](https://https://sendgrid.com//) email service.

#### How to use

In your `mailers.go`

```go
import (
    ... 
    ssender "github.com/paganotoni/sendgrid-sender"
)

var sender mail.Sender

func init() {
	APIKey := envy.Get("SENDGRID_API_KEY", "")
	sender = ssender.NewSendgridSender(APIKey)
}
```

And then in your mailers you would do the same `sender.Send(m)` as this sender matches buffalos [`mail.Sender`](https://github.com/gobuffalo/buffalo/blob/master/mail/mail.go#L4) interface.

#### Test mode

Whenever the GO_ENV variable is set to be `test` this sender will use [mocksmtp](https://github.com/stanislas-m/mocksmtp) sender to send messages, you can read values in your tests within the property `TestSender` of the SendgridSender.