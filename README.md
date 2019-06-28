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

To add custom args, you must add values using `CustomArgs` type and add it into Message.Data ussing `CustomArgsKey` key

```go
CustomArgsKey = "sendgrid_custom_args_key"
...
type CustomArgs map[string]string
```

#### How to use Sendgrid customargs

One thing you could need is to add customArgs to the message you're sending through sendgrid, to do this you would be using `SetCustomArgs` function, passing your `mail.Message`with the `CustomArgs` you want to add.

```go
import (
	...
    ssender "github.com/paganotoni/sendgrid-sender"
)

func main() {
	APIKey := envy.Get("SENDGRID_API_KEY", "")
    sender = ssender.NewSendgridSender(APIKey)

	m := mail.NewMessage()
    ...
    ssender.SetCustomArgs(m, ssender.CustomArgs{
        "custom_arg_0": "custom_value_0",
        "custom_arg_1": "custom_value_1",
        ...
    })

    if err := sender.Send(m); err != nil{
        ...
    }
}
```

#### Test mode

Whenever the GO_ENV variable is set to be `test` this sender will use [mocksmtp](https://github.com/stanislas-m/mocksmtp) sender to send messages, you can read values in your tests within the property `TestSender` of the SendgridSender.
