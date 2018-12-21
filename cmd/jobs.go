package cmd

import (
	"encoding/json"
	"fmt"
	//"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	"github.com/streadway/amqp"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice/controllers"
	"message_notification_practice/model"
	"message_notification_practice/mq"
	"message_notification_practice/pb"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	jobsCmdNum  int
	jobsCmdType string
)

var senderMqCfgMap = map[string]mq.MQCfg{
	model.MsgTypeMail: {
		URL:      "amqp://liujx:Liujiaxing@localhost:5672/",
		Queue:    "push.msg.q.notification.mail",
		Exchange: "t.msg.ex.notification",
	},
	model.MsgTypePhone: {
		URL:      "amqp://liujx:Liujiaxing@localhost:5672/",
		Queue:    "push.msg.q.notification.phone",
		Exchange: "t.msg.ex.notification",
	},
	model.MsgTypeWeChat: {
		URL:      "amqp://liujx:Liujiaxing@localhost:5672/",
		Queue:    "push.msg.q.notification.wechat",
		Exchange: "t.msg.ex.notification",
	},
}

var domain string = "sandboxaaff1b769a3c429daef777dfeed8f173.mailgun.org" // e.g. mg.yourcompany.com
var privateAPIKey string = "61604cff9615cfa175f4340991b8c713-9b463597-ad5d1076"
var publicAPIKey string = "pubkey-25f9fdfa7af58880311ec28977c10f6c"

var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Start job for notification or sender",
}

var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Start notification job",
	Run:   notificationProc,
}

func notificationProc(cmd *cobra.Command, args []string) {

	log.Debug("Start Jobs Notify!")

	db, err := gorm.Open("mysql", "root:123456@/msg_notification?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("init mysql failed, err: %s", err)
		return
	}

	defer db.Close()
	// 消费
	chanDepth := 10 * jobsCmdNum
	rcvDataChan := make(chan interface{}, chanDepth)
	//sendDataChan := make(chan interface{}, chanDepth)

	consumerInfo := mq.MQInfo{
		Cfg: mq.MQCfg{
			URL:      "amqp://liujx:Liujiaxing@localhost:5672/",
			Queue:    "push.msg.q",
			Exchange: "t.msg.ex",
		},
		MsgChan: rcvDataChan,
	}

	if err := mq.ConsumerStart(jobsCmdNum, consumerInfo); err != nil {
		log.Errorf("ConsumerStart err：%v", err)
		return
	}

	sendChanMap := map[string]chan interface{}{
		model.MsgTypeMail:   make(chan interface{}, chanDepth),
		model.MsgTypePhone:  make(chan interface{}, chanDepth),
		model.MsgTypeWeChat: make(chan interface{}, chanDepth),
	}

	for tp, cfg := range senderMqCfgMap {
		if err := mq.ProducerStart(jobsCmdNum, mq.MQInfo{Cfg: cfg, MsgChan: sendChanMap[tp]}); err != nil {
			log.Errorf("Jobs nitificaiton start producers failed, err：%v", err)
			return
		}
	}

	// 业务处理协程
	busiRoutineNum := 3 * jobsCmdNum

	for i := 0; i < busiRoutineNum; i++ {
		go func(id int, rcvDataChan chan interface{}, sendDataChan map[string]chan interface{}) {

			for {
				select {

				case msg, ok := <-rcvDataChan:

					if ok {
						rq := &pb.MsgNotificationRequest{}

						err := json.Unmarshal(msg.(amqp.Delivery).Body, rq)
						if err != nil {
							log.Error("Unmarshal MsgNotificationRequest Body failed, err: ", err)
						}

						var users []model.User

						// TODO: 需要整理放入 service 层
						if e := db.Raw("select * from users where id in (select user_id from group_user_relations where group_id = ?)",
							rq.Group).Scan(&users).Error; e != nil {
							log.Error("get group_user_relations failed, err: ", e)
						}

						for _, user := range users {

							userMsg := &model.UserMsg{
								ID:      user.ID,
								Name:    user.Name,
								Content: rq.Content,
								Email:   user.Email,
								Phone:   user.Phone,
								WeChat:  user.Wechat,
							}

							bytes, err := json.Marshal(&userMsg)
							if err != nil {
								log.Error("Email marshal UserMsg failed, err: ", err)
							}

							if len(user.Email) > 0 {
								if ch, ok := sendDataChan[model.MsgTypeMail]; ok {
									ch <- bytes
								} else {
									log.Error("Email send channel dose not exist.")
								}
							}

							if len(user.Phone) > 0 {
								if ch, ok := sendDataChan[model.MsgTypePhone]; ok {
									ch <- bytes
								} else {
									log.Error("Phone send channel dose not exist.")
								}
							}

							if len(user.Wechat) > 0 {

								if ch, ok := sendDataChan[model.MsgTypeWeChat]; ok {
									ch <- bytes
								} else {
									log.Error("Wechat send channel dose not exist.")
								}
							}
						}

						log.Debugf("group_id: %d, %#v\n", rq.Group, users)

						time.Sleep(2 * time.Second) // TODO: remove debug
					}
				}
			}
		}(i, rcvDataChan, sendChanMap)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Debug("Exit Jobs Notification!")
}

var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "Start sender job",
	Run:   senderProc,
}

func senderProc(cmd *cobra.Command, args []string) {

	log.Debug("Start Jobs Sender!")

	// 消费
	chanDepth := 10 * jobsCmdNum
	rcvDataChan := make(chan interface{}, chanDepth)

	mqCfg, ok := senderMqCfgMap[jobsCmdType]
	if !ok {
		log.Errorf("not found %s mq config info", jobsCmdType)
		return
	}

	consumerInfo := mq.MQInfo{
		Cfg:     mqCfg,
		MsgChan: rcvDataChan,
	}

	if err := mq.ConsumerStart(jobsCmdNum, consumerInfo); err != nil {
		log.Errorf("ConsumerStart err：%v", err)
		return
	}

	ctl := controllers.NewSenderController(jobsCmdType, domain, privateAPIKey, publicAPIKey)
	//ctl := controllers.NewSenderController(jobsCmdType)

	// 业务处理协程

	for i := 0; i < jobsCmdNum; i++ {

		go func(id int, rcvDataChan chan interface{}, handler controllers.SenderHandler) {

			for {
				select {
				case msg, ok := <-rcvDataChan:
					if ok {
						userMsg := model.UserMsg{}

						if err := json.Unmarshal(msg.(amqp.Delivery).Body, &userMsg); err != nil {
							log.Error("Unmarshal MsgNotificationRequest Body failed, err: ", err)
						}

						r := retrier.New(retrier.ExponentialBackoff(5, 2*time.Second), nil)

						err := r.Run(func() error {
							log.Info(time.Now().Second())
							return handler(&userMsg)
						})

						if err != nil {
							log.Error("get an error, handle it, err: ", err)
						}
					}
				}
			}

		}(i, rcvDataChan, ctl.Handler)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Debug("Exit Jobs Sender!")
}

func init() {
	rootCmd.AddCommand(jobsCmd)

	jobsCmd.AddCommand(notificationCmd)
	jobsCmd.AddCommand(senderCmd)

	jobsCmd.PersistentFlags().IntVarP(&jobsCmdNum, "number", "n", 3, "jobs number")

	senderCmd.PersistentFlags().StringVarP(&jobsCmdType, "type", "t", "mail", "jobs number")
}
