package postoffice

import ()

var defaultPostman *Postman

func init() {
	defaultPostman = NewPostman()
}

func SendMail(subject, content, from string, to, cc, bcc []string, contentType, charset string) error {
	em := NewEmail(subject, content, from, to, cc, bcc, contentType, charset)
	return defaultPostman.Send(em)
}

func SendSimpleEmail(subject, content, from string, to []string) error {
	em := NewSimpleEmail(subject, content, from, to)
	return defaultPostman.Send(em)
}
