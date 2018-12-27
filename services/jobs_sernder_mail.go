package services

import (
	"errors"
	"github.com/eapache/go-resiliency/retrier"
	log "gopkg.in/cihub/seelog.v2"
	"gopkg.in/mailgun/mailgun-go.v1"
	"message_notification_practice/model"
	"time"
)

func NewMailSenderService(mg interface{}) *MailSenderService {
	return &MailSenderService{mg: mg.(mailgun.Mailgun)}
}

type MailSenderService struct {
	mg mailgun.Mailgun
}

func (svc *MailSenderService) Handler(msg *model.UserMsg) error {

	// TODO: not implementation
	log.Debugf("MailSenderService: userMsg: %#v\n", msg)
	return nil

	r := retrier.New(retrier.ExponentialBackoff(5, 20*time.Millisecond), nil)

	// body := template.New() // TODO: HTML mail is not implementation.

	err := r.Run(func() error {

		return errors.New("sender handler happens error")

		log.Info(time.Now().Second())
		resp, id, err := svc.mg.Send(svc.mg.NewMessage(
			"aaa <83214742@qq.com>",
			"Hello",
			msg.Content,
			msg.Email,
		))
		log.Debugf("Handler():resp: %s, id: %s, msg: %#v", resp, id, msg)

		return err

	})

	time.Sleep(5 * time.Second)

	return err
}
