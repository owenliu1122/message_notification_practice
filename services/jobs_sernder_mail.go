package services

import (
	"github.com/eapache/go-resiliency/retrier"
	log "gopkg.in/cihub/seelog.v2"
	"gopkg.in/mailgun/mailgun-go.v1"
	"message_notification_practice"
	"time"
)

//var domain string = "sandboxaaff1b769a3c429daef777dfeed8f173.mailgun.org" // e.g. mg.yourcompany.com
//var privateAPIKey string = "61604cff9615cfa175f4340991b8c713-9b463597-ad5d1076"
//var publicAPIKey string = "pubkey-25f9fdfa7af58880311ec28977c10f6c"

func NewMailSenderService(toolCfg interface{}) *MailSenderService {
	cfg := toolCfg.(map[string]string)

	return &MailSenderService{
		mg: mailgun.NewMailgun(
			cfg["domain"],
			cfg["privateapikey"],
			cfg["publicapikey"],
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
