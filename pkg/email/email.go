package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
)

type Email struct {
	NickName string
	From     string
	Host     string
	Html     string
	IsSsl    bool
	Port     int
	Secret   string
	Subject  string
}

// Init 初始化邮箱
func (E *Email) Init() {
	E.NickName = viper.GetString("email.nickname")
	E.From = viper.GetString("email.from")
	E.Host = viper.GetString("email.host")
	E.Port = viper.GetInt("email.port")
	E.Secret = viper.GetString("email.secret")
	E.Html = viper.GetString("email.html")
	E.Subject = viper.GetString("email.subject")
	E.IsSsl = viper.GetBool("email.is_ssl")
}

// SendEmailVerify 发送验证码
func SendEmailVerify(toEmail string, code int64) error {
	em := Email{}
	em.Init()
	client := email.NewEmail()
	client = &email.Email{
		From:    fmt.Sprintf("%s <%s>", em.NickName, em.From),
		To:      []string{toEmail},
		Subject: em.Subject,
		HTML:    []byte(fmt.Sprintf(em.Html, code)),
	}
	auth := smtp.PlainAuth("", em.From, em.Secret, em.Host)

	if em.IsSsl {
		cfg := &tls.Config{
			ServerName: em.Host,
		}
		return client.SendWithTLS(fmt.Sprintf("%s:%d", em.Host, em.Port), auth, cfg)
	}
	return client.Send(fmt.Sprintf("%s:%d", em.Host, em.Port), auth)
}
