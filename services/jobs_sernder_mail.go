package services

import (
	"context"
	b64 "encoding/base64"
	"time"

	"github.com/fpay/foundation-go"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/fpay/foundation-go/log"
	"github.com/owenliu1122/notice"
	"gopkg.in/mailgun/mailgun-go.v1"
)

// NewMailSenderService return a mail sender service.
func NewMailSenderService(logger *log.Logger, cfg notice.SendServiceConfig, pc foundation.JobManager) *MailSenderService {

	domain, _ := b64.StdEncoding.DecodeString(cfg.Domain)
	privateapikey, _ := b64.StdEncoding.DecodeString(cfg.PrivateAPIKey)
	publicapikey, _ := b64.StdEncoding.DecodeString(cfg.PublicAPIKey)

	logger.Info("mailgun domain: ", string(domain))
	logger.Info("mailgun apikey: ", string(privateapikey))
	logger.Info("mailgun pubkey: ", string(publicapikey))

	return &MailSenderService{
		logger: logger,
		pc:     pc,
		mg: mailgun.NewMailgun(
			string(domain),
			string(privateapikey),
			string(publicapikey),
		),
	}
}

// MailSenderService is a mail sender service.
type MailSenderService struct {
	logger *log.Logger
	mg     mailgun.Mailgun
	pc     foundation.JobManager
}

// Handler parse a email message that needs to be sent.
func (svc *MailSenderService) Handler(ctx context.Context, msg *notice.UserMessage) error {

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

	return err
}
