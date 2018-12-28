package services

import (
	b64 "encoding/base64"
	"errors"
	"github.com/eapache/go-resiliency/retrier"
	log "gopkg.in/cihub/seelog.v2"
	"gopkg.in/mailgun/mailgun-go.v1"
	"message_notification_practice"
	"time"
)

func NewMailSenderService(cfg map[string]string) *MailSenderService {
	//cfg := toolCfg.(map[string]string)

	domain, _ := b64.StdEncoding.DecodeString(cfg["domain"])
	privateapikey, _ := b64.StdEncoding.DecodeString(cfg["privateapikey"])
	publicapikey, _ := b64.StdEncoding.DecodeString(cfg["publicapikey"])

	log.Info("mailgun domain: ", string(domain))
	log.Info("mailgun apikey: ", string(privateapikey))
	log.Info("mailgun pubkey: ", string(publicapikey))

	return &MailSenderService{
		mg: mailgun.NewMailgun(
			string(domain),
			string(privateapikey),
			string(publicapikey),
		),
	}
}

type MailSenderService struct {
	mg mailgun.Mailgun
}

func (svc *MailSenderService) Handler(msg *root.UserMsg) error {

	log.Debugf("MailSenderService: userMsg: %#v\n", msg)
	//return nil

	r := retrier.New(retrier.ExponentialBackoff(5, 20*time.Millisecond), nil)

	log.Debugf("MailSenderService: mg: %#v", svc.mg)

	err := r.Run(func() error {

		return errors.New("sender handler happens error")

		resp, id, err := svc.mg.Send(svc.mg.NewMessage(
			"aaa <83214742@qq.com>",
			"Hello",
			msg.Content,
			msg.Email,
		))
		log.Debugf("Handler():resp: %s, id: %s, msg: %#v", resp, id, msg)

		return err

	})

	time.Sleep(2 * time.Second)

	return err
}
