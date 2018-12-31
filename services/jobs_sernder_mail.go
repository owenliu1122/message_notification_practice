package services

import (
	b64 "encoding/base64"
	"message_notification_practice"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	log "gopkg.in/cihub/seelog.v2"
	"gopkg.in/mailgun/mailgun-go.v1"
)

// NewMailSenderService return a mail sender service.
func NewMailSenderService(cfg notice.SendService) *MailSenderService {
	//cfg := toolCfg.(map[string]string)

	domain, _ := b64.StdEncoding.DecodeString(cfg.Domain)
	privateapikey, _ := b64.StdEncoding.DecodeString(cfg.PrivateAPIKey)
	publicapikey, _ := b64.StdEncoding.DecodeString(cfg.PublicAPIKey)

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

// MailSenderService is a mail sender service.
type MailSenderService struct {
	mg mailgun.Mailgun
}

// Handler parse a email message that needs to be sent.
func (svc *MailSenderService) Handler(msg *notice.UserMessage) error {

	log.Debugf("MailSenderService: userMsg: %#v\n", msg)

	r := retrier.New(retrier.ExponentialBackoff(5, 20*time.Millisecond), nil)

	log.Debugf("MailSenderService: mg: %#v", svc.mg)

	err := r.Run(func() error {

		//return errors.New("sender handler happens error")

		resp, id, err := svc.mg.Send(svc.mg.NewMessage(
			"aaa <83214742@qq.com>",
			"Hello",
			msg.Content,
			msg.Email,
		))
		log.Debugf("Handler():resp: %s, id: %s, msg: %#v", resp, id, msg)

		return err

	})

	return err
}
