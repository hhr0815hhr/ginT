package mail

import (
	"fmt"

	gmail "github.com/go-mail/mail/v2"
	"github.com/hhr0815hhr/gint/internal/config"
	"github.com/hhr0815hhr/gint/internal/pkg/i18n"
)

var client *gmail.Dialer

func Instance() *gmail.Dialer {
	client = gmail.NewDialer(
		config.Conf.Server.Mail.Host,
		config.Conf.Server.Mail.Port,
		config.Conf.Server.Mail.User,
		config.Conf.Server.Mail.Passwd,
	)
	client.StartTLSPolicy = gmail.MandatoryStartTLS
	return client
}

func SendHyperVerifyText(locale string, verificationCode string, to string) error {
	verificationURL := fmt.Sprintf("%s?code=%s", config.Conf.Server.Mail.VerifyUrl, verificationCode)

	htmlBody := fmt.Sprintf(i18n.Tl(locale, "mail.verifyTemplate"), verificationCode, verificationURL)
	m := gmail.NewMessage()
	m.SetHeader("From", config.Conf.Server.Mail.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", i18n.Tl(locale, "mail.verifyHeader"))
	m.SetBody("text/plain", fmt.Sprintf(i18n.Tl(locale, "mail.verifyBody"), verificationCode, verificationURL)) // 纯文本备选
	m.AddAlternative("text/html", htmlBody)
	if err := Instance().DialAndSend(m); err != nil {
		return err
	}
	return nil
}
