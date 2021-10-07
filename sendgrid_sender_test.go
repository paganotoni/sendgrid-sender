package sender

import (
	"bytes"
	"testing"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/stretchr/testify/require"
)

func Test_build_Mail(t *testing.T) {
	a := require.New(t)

	cases := []struct {
		Correct     bool
		From        string
		To          []string
		Bodies      []mail.Body
		Attachments []mail.Attachment
		Error       string
	}{
		{
			Correct: true,
			From:    "tatan@test.com",
			To:      []string{"email@test.com", "anotheremail@test.com"},
			Bodies: []mail.Body{
				{
					Content:     "<p>Test Content of mail</p>",
					ContentType: "text/html",
				},

				{
					Content:     "Test Content of mail",
					ContentType: "text/plain",
				},
			},
			Attachments: []mail.Attachment{},
			Error:       "",
		},
		{
			Correct: false,
			From:    "",
			To:      []string{"email@test.com", "anotheremail@test.com"},
			Bodies: []mail.Body{
				{
					Content:     "<p>Test Content of mail</p>",
					ContentType: "text/html",
				},

				{
					Content:     "Test Content of mail",
					ContentType: "text/plain",
				},
			},
			Attachments: []mail.Attachment{},
			Error:       "invalid from (<nil>): mail: no address",
		},
		{
			Correct: false,
			From:    "tatan@test.com",
			To:      []string{"", "anotheremail@test.com"},
			Bodies: []mail.Body{
				{
					Content:     "<p>Test Content of mail</p>",
					ContentType: "text/html",
				},

				{
					Content:     "Test Content of mail",
					ContentType: "text/plain",
				},
			},
			Attachments: []mail.Attachment{},
			Error:       "invalid to (): mail: no address",
		},
		// Attachments
		{
			Correct: true,
			From:    "tatan@test.com",
			To:      []string{"email@test.com", "anotheremail@test.com"},
			Bodies: []mail.Body{
				{
					Content:     "<p>Test Content of mail</p>",
					ContentType: "text/html",
				},

				{
					Content:     "Test Content of mail",
					ContentType: "text/plain",
				},
			},
			Attachments: []mail.Attachment{
				{
					Name:        "test_file.pdf",
					Reader:      bytes.NewReader([]byte("TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4gQ3JhcyBwdW12")),
					ContentType: "application/pdf",
					Embedded:    false,
				},
				{
					Name:        "test_image.png",
					Reader:      bytes.NewReader([]byte("R29zIGxvdmVzIHlvdQ==")),
					ContentType: "image/png",
					Embedded:    true,
				},
			},
			Error: "",
		},
	}

	for i, c := range cases {
		m := mail.NewMessage()

		m.From = c.From
		m.Subject = "Test Mail"
		m.To = c.To
		m.Bodies = c.Bodies
		m.Attachments = c.Attachments

		mm, err := buildMail(m)

		if !c.Correct {
			a.Errorf(err, "Loop %d", i)
			a.EqualErrorf(err, c.Error, "Loop: %v", i)
			continue
		}

		a.Equalf(len(m.Attachments), len(mm.Attachments), "Loop %d", i)

		for j, at := range mm.Attachments {
			a.Equalf(m.Attachments[j].Name, at.Filename, "Loop %d", i)
			if !m.Attachments[j].Embedded {
				a.Equalf("attachment", at.Disposition, "Loop %d", i)
			} else {
				a.Equalf("inline", at.Disposition, "Loop %d", i)
			}
			a.Equalf(m.Attachments[j].ContentType, at.Type, "Loop %d", i)
		}

	}
}

func Test_build_Mail_Custom_Args(t *testing.T) {
	a := require.New(t)
	m := mail.NewMessage()

	m.From = "tatan@test.com"
	m.Subject = "Test Mail"
	m.To = []string{"email@test.com", "anotheremail@test.com"}
	m.Bodies = []mail.Body{
		{
			Content:     "<p>Test Content of mail</p>",
			ContentType: "text/html",
		},

		{
			Content:     "Test Content of mail",
			ContentType: "text/plain",
		},
	}
	m.Attachments = []mail.Attachment{
		{
			Name:        "test_file.pdf",
			Reader:      bytes.NewReader([]byte("TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4gQ3JhcyBwdW12")),
			ContentType: "application/pdf",
			Embedded:    false,
		},
		{
			Name:        "test_image.png",
			Reader:      bytes.NewReader([]byte("R29zIGxvdmVzIHlvdQ==")),
			ContentType: "image/png",
			Embedded:    true,
		},
	}

	m.Data = map[string]interface{}{}
	SetCustomArgs(m, CustomArgs{
		"custom_key_0": "custom_value_0",
		"custom_key_1": "custom_value_1",
		"custom_key_2": "[val_2 val_3 val_4]",
	})

	mm, err := buildMail(m)

	a.NoError(err)
	a.Equal(3, len(mm.Personalizations[0].CustomArgs))
	a.Equal("custom_value_0", mm.Personalizations[0].CustomArgs["custom_key_0"])
	a.Equal("custom_value_1", mm.Personalizations[0].CustomArgs["custom_key_1"])
	a.Equal("[val_2 val_3 val_4]", mm.Personalizations[0].CustomArgs["custom_key_2"])

	m.Data = map[string]interface{}{}
	mm, err = buildMail(m)

	a.NoError(err)
	a.Equal(0, len(mm.Personalizations[0].CustomArgs))
	a.Equal("", mm.Personalizations[0].CustomArgs["custom_key_0"])
	a.Equal("", mm.Personalizations[0].CustomArgs["custom_key_1"])
	a.Equal("", mm.Personalizations[0].CustomArgs["custom_key_2"])
}

func Test_build_Mail_Custom_Headers(t *testing.T) {
	a := require.New(t)
	m := mail.NewMessage()

	m.From = "digi@charat.com"
	m.Subject = "Test Header Mail"
	m.To = []string{"email@test.com", "anotheremail@test.com"}
	m.Headers = map[string]string{
		"x-provider": "sendgrid",
		"x-header":   "testing_header",
	}
	m.Bodies = []mail.Body{
		{
			Content:     "<p>Test Content of mail</p>",
			ContentType: "text/html",
		},

		{
			Content:     "Test Content of mail",
			ContentType: "text/plain",
		},
	}

	mm, err := buildMail(m)
	a.NoError(err)
	a.Equal(2, len(mm.Personalizations[0].Headers))
	a.Equal("sendgrid", mm.Personalizations[0].Headers["x-provider"])
	a.Equal("testing_header", mm.Personalizations[0].Headers["x-header"])
}
