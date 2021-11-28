package xemail

import (
	"encoding/json"

	"gopkg.in/gomail.v2"
)

type config struct {
	User string `json:"user" yaml:"user"`
	Pass string `json:"pass" yaml:"pass"`
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

type Client struct {
	User string
	*gomail.Dialer
}

func New(confStr string) *Client {
	cfg := config{}
	if err := json.Unmarshal([]byte(confStr), &cfg); err != nil {
		panic(err.Error())
	}
	return &Client{
		User:   cfg.User,
		Dialer: gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Pass),
	}
}

// SendMail 发送单封邮件
func (c Client) SendMail(mailTo []string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "<"+c.User+">")
	m.SetHeader("To", mailTo...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	return c.DialAndSend(m)
}
