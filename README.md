### Sendgrid Buffalo Sender

This is a [buffalo](https://github.com/gobuffalo/buffalo) sender for the [Sendgrid](https://sendgrid.com) email service.

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

#### Add Custom Args

If you want to send a mail with `custom args`, you need to add it into `mail.Message.Data` as `map[string]interface{}` this will be into `mail.Message.Personalizations` field :

```go
import (
    ... 
    ssender "github.com/paganotoni/sendgrid-sender"
    bmail "github.com/gobuffalo/buffalo/mail"
)

var sender mail.Sender

func init() {
    APIKey := envy.Get("SENDGRID_API_KEY", "")
    sender = ssender.NewSendgridSender(APIKey)
    ...
    message := bmail.NewMessage()
    ...
    customArgs := map[string]interface{}{
        "custom_arg_0": "custom_value_0",
        "custom_arg_1": 100,
        "custom_arg_2": []string{"val_0", "val_1", "val_2"},
        "custom_arg_3": map[string]string{
            "firstName": "John",
            "lastName":  "Smith",
            "age":       "24",
            },
        ...
    }

    message.Data = customArgs
    ...

}
```

#### Test mode

Whenever the GO_ENV variable is set to be `test` this sender will use [mocksmtp](https://github.com/stanislas-m/mocksmtp) sender to send messages, you can read values in your tests within the property `TestSender` of the SendgridSender.
