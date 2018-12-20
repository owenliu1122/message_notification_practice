package controllers

import (
	"errors"
	"github.com/eapache/go-resiliency/retrier"
	log "gopkg.in/cihub/seelog.v2"
	"gopkg.in/mailgun/mailgun-go.v1"
	"html/template"
	"message_notification_practice/model"
	"time"
)

type MailSenderController struct {
	mg mailgun.Mailgun
}

func NewMailSenderController(domain, apiKey, pubKey string) *MailSenderController {
	return &MailSenderController{
		mg: mailgun.NewMailgun(domain, apiKey, pubKey),
	}
}

func (ctl *MailSenderController) Handler(msg *model.UserMsg) error {

	r := retrier.New(retrier.ExponentialBackoff(5, 20*time.Millisecond), nil)
	body := template.New()
	err := r.Run(func() error {

		return errors.New("sender handler happens error")

		log.Info(time.Now().Second())
		resp, id, err := ctl.mg.Send(ctl.mg.NewMessage(
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
