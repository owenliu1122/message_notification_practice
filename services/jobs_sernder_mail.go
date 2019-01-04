package services

import (
	b64 "encoding/base64"
	"encoding/json"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/fpay/foundation-go/log"
	"github.com/owenliu1122/notice"
	"gopkg.in/mailgun/mailgun-go.v1"
)

// NewMailSenderService return a mail sender service.
func NewMailSenderService(logger *log.Logger, cfg notice.SendServiceConfig, pc notice.ProducerInterface, exRouting notice.ProducerConfig) *MailSenderService {
	//cfg := toolCfg.(map[string]string)

	domain, _ := b64.StdEncoding.DecodeString(cfg.Domain)
	privateapikey, _ := b64.StdEncoding.DecodeString(cfg.PrivateAPIKey)
	publicapikey, _ := b64.StdEncoding.DecodeString(cfg.PublicAPIKey)

	logger.Info("mailgun domain: ", string(domain))
	logger.Info("mailgun apikey: ", string(privateapikey))
	logger.Info("mailgun pubkey: ", string(publicapikey))

	return &MailSenderService{
		logger:    logger,
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
	logger    *log.Logger
	mg        mailgun.Mailgun
	pc        notice.ProducerInterface
	exRouting notice.ProducerConfig
}

// Handler parse a email message that needs to be sent.
func (svc *MailSenderService) Handler(msg *notice.UserMessage) error {

	svc.logger.Debugf("MailSenderService: userMsg: %#v\n", msg)
	r := retrier.New(retrier.ExponentialBackoff(3, 20*time.Millisecond), nil)

	err := r.Run(func() error {

		//return errors.New("sender handler happens error")

		_, _, err := svc.mg.Send(svc.mg.NewMessage(
			"aaa <83214742@qq.com>",
			"Hello",
			msg.Content,
			msg.Destination,
		))
		//svc.logger.Debugf("Handler():resp: %s, id: %s, msg: %#v", resp, id, msg)

		return err

	})

	if err != nil {
		var jsonBytes []byte
		jsonBytes, err = json.Marshal(msg)
		if err != nil {
			svc.logger.Error("publish to retr, marshal msg Body failed, err: ", err)
			return err
		}

		if err = svc.pc.Publish(svc.exRouting.Exchange, svc.exRouting.RoutingKey, jsonBytes); err != nil {
			return err
		}
	}

	return err
}
