package postoffice

import (
	"testing"
)

func TestPostman_Send(t *testing.T) {
	from := "zhangsan@domain.com"
	to := []string{"lisi@domain.com"}
	cc := []string{}
	bcc := []string{}

	em := NewEmail(
		"test mail subject测试邮件主题",
		"test mail content测试邮件内容",
		from,
		to,
		cc,
		bcc,
		ContentTypeTextPlain,
		CharsetUtf8,
	)
	pm := NewPostman()

	err := pm.Send(em)
	if err != nil {
		t.Errorf("send mail failed, error: %v", err)
	}
}
