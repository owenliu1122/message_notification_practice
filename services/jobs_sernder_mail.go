package services

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/owenliu1122/notice/mq"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/owenliu1122/notice"
	log "gopkg.in/cihub/seelog.v2"
	"gopkg.in/mailgun/mailgun-go.v1"
)

// NewMailSenderService return a mail sender service.
func NewMailSenderService(cfg notice.SendService, pc *mq.Producer, exRouting notice.Producer) *MailSenderService {
	//cfg := toolCfg.(map[string]string)

	domain, _ := b64.StdEncoding.DecodeString(cfg.Domain)
	privateapikey, _ := b64.StdEncoding.DecodeString(cfg.PrivateAPIKey)
	publicapikey, _ := b64.StdEncoding.DecodeString(cfg.PublicAPIKey)

	log.Info("mailgun domain: ", string(domain))
	log.Info("mailgun apikey: ", string(privateapikey))
	log.Info("mailgun pubkey: ", string(publicapikey))

	return &MailSenderService{
		pc:        pc,
		exRouting: exRouting,
		mg: mailgun.NewMailgun(
			string(domain),
			string(privateapikey),
			string(publicapikey),
		),
	}
}

// MailSenderService is a mail sender service.
type MailSenderService struct {
	mg        mailgun.Mailgun
	pc        *mq.Producer
	exRouting notice.Producer
}

// Handler parse a email message that needs to be sent.
func (svc *MailSenderService) Handler(msg *notice.UserMessage) error {

	log.Debugf("MailSenderService: userMsg: %#v\n", msg)
	r := retrier.New(retrier.ExponentialBackoff(3, 20*time.Millisecond), nil)

	err := r.Run(func() error {

		return errors.New("sender handler happens error")

		_, _, err := svc.mg.Send(svc.mg.NewMessage(
			"aaa <83214742@qq.com>",
			"Hello",
			msg.Content,
			msg.Destination,
		))
		//log.Debugf("Handler():resp: %s, id: %s, msg: %#v", resp, id, msg)

		return err

	})

	if err != nil {
		jsonBytes, err := json.Marshal(msg)
		if err != nil {
			log.Error("publish to retr, marshal msg Body failed, err: ", err)
			return err
		}

		if err = svc.pc.Publish(svc.exRouting.Exchange, svc.exRouting.RoutingKey, jsonBytes); err != nil {
			return err
		}
	}

	return err
}
