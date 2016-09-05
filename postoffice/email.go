package postoffice

import (
	"andals/gobox/encoding"
)

const (
	ContentTypeTextPlain = "text/plain"
	ContentTypeTextHtml  = "text/html"

	CharsetUtf8 = "UTF-8"
)

type Email struct {
	From string
	To   []string
	Cc   []string
	Bcc  []string

	ContentType string
	Charset     string

	Subject string
	Content string
}

func NewEmail(subject, content, from string, to, cc, bcc []string, contentType, charset string) *Email {

	subject = "=?" + charset + "?B?" + string(encoding.Base64Encode([]byte(subject))) + "?="
	subject = subject + " \nContent-Type: " + contentType + ";charset=" + charset

	em := &Email{
		From: from,
		To:   to,
		Cc:   cc,
		Bcc:  bcc,

		ContentType: contentType,
		Charset:     charset,

		Subject: subject,
		Content: content,
	}
	return em
}

func NewSimpleEmail(subject, content, from string, to []string) *Email {
	return NewEmail(subject, content, from, to, nil, nil, ContentTypeTextPlain, CharsetUtf8)
}
