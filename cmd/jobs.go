package cmd

import (
	"fmt"
	"gopkg.in/mailgun/mailgun-go.v1"
	"message_notification_practice/services"

	"context"
	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice/controllers"
	//"message_notification_practice/model"
	"message_notification_practice/mq"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	jobsCmdNum  int
	jobsCmdType string
)

var senderMqCfgMap = map[string]mq.BaseProducer{
	services.MsgTypeMail: {
		//URL:      "amqp://liujx:Liujiaxing@localhost:5672/",
		RoutingKey: "push.msg.q.notification.mail",
		Exchange:   "t.msg.ex.notification",
	},
	services.MsgTypePhone: {
		//URL:      "amqp://liujx:Liujiaxing@localhost:5672/",
		RoutingKey: "push.msg.q.notification.phone",
		Exchange:   "t.msg.ex.notification",
	},
	services.MsgTypeWeChat: {
		//URL:      "amqp://liujx:Liujiaxing@localhost:5672/",
		RoutingKey: "push.msg.q.notification.wechat",
		Exchange:   "t.msg.ex.notification",
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
	ctx, cancel := context.WithCancel(context.Background())

	db, err := gorm.Open("mysql", "root:123456@/msg_notification?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Printf("init mysql failed, err: %s", err)
		return
	}

	defer db.Close()

	// 业务处理协程

	cmCfg := mq.MQCfg{
		URL:      "amqp://liujx:Liujiaxing@localhost:5672/",
		Queue:    "push.msg.q",
		Exchange: "t.msg.ex",
	}

	mqCli := mq.NewMq(cmCfg)
	if e := mqCli.InitConnection(); e != nil {
		log.Error("InitConnection failed, err: ", e)
	}
	defer mqCli.Close()

	if e := mqCli.InitProducer(cmCfg.Exchange, cmCfg.Queue); e != nil {
		log.Error("InitProducer failed, err: ", e)
	}

	mqSendSvc := services.NewMqSendService(mqCli, services.NewGroupUserRelationService(db))

	for k := range senderMqCfgMap {
		mqSendSvc.RegisterExchangeRouting(k, senderMqCfgMap[k])
	}

	ctl := controllers.NewNotificationController(mqSendSvc)

	for i := 0; i < jobsCmdNum; i++ {

		mqCli.RegisterConsumer(
			fmt.Sprintf("notification[%d]", i),
			ctl.Handler,
			mq.BaseConsumer{
				Queue:   cmCfg.Queue,
				AutoAck: true,
			})
	}

	if e := mqCli.StartConsumer(ctx); e != nil {
		log.Error("InitConnection failed, err: ", e)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	cancel()

	time.Sleep(2 * time.Second)

	log.Debug("Exit Jobs Notification!")
}

var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "Start sender job",
	Run:   senderProc,
}

func senderProc(cmd *cobra.Command, args []string) {

	log.Debug("Start Jobs Sender!")

	ctx, cancel := context.WithCancel(context.Background())

	var senderJobCfgMap = map[string]mq.MQCfg{
		services.MsgTypeMail: {
			URL:      "amqp://liujx:Liujiaxing@localhost:5672/",
			Queue:    "push.msg.q.notification.mail",
			Exchange: "t.msg.ex.notification",
		},
		services.MsgTypePhone: {
			URL:      "amqp://liujx:Liujiaxing@localhost:5672/",
			Queue:    "push.msg.q.notification.phone",
			Exchange: "t.msg.ex.notification",
		},
		services.MsgTypeWeChat: {
			URL:      "amqp://liujx:Liujiaxing@localhost:5672/",
			Queue:    "push.msg.q.notification.wechat",
			Exchange: "t.msg.ex.notification",
		},
	}
	// 消费
	cmCfg, ok := senderJobCfgMap[jobsCmdType]
	if !ok {
		log.Errorf("not found %s mq config info", jobsCmdType)
		return
	}

	mqCli := mq.NewMq(cmCfg)
	if e := mqCli.InitConnection(); e != nil {
		log.Error("InitConnection failed, err: ", e)
	}
	defer mqCli.Close()

	// TODO: 需要使用统一的接口，这里暂时时候 mail 接口测试
	sendSvc := services.NewSenderService(jobsCmdType, mailgun.NewMailgun(domain, privateAPIKey, publicAPIKey))
	ctl := controllers.NewSenderController(sendSvc)

	for i := 0; i < jobsCmdNum; i++ {

		mqCli.RegisterConsumer(
			fmt.Sprintf("notification[%d]", i),
			ctl.Handler,
			mq.BaseConsumer{
				Queue:   cmCfg.Queue,
				AutoAck: true,
			})
	}

	if e := mqCli.StartConsumer(ctx); e != nil {
		log.Error("InitConnection failed, err: ", e)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	cancel()

	log.Debug("Exit Jobs Sender!")
}

func init() {
	rootCmd.AddCommand(jobsCmd)

	jobsCmd.AddCommand(notificationCmd)
	jobsCmd.AddCommand(senderCmd)

	jobsCmd.PersistentFlags().IntVarP(&jobsCmdNum, "number", "n", 3, "jobs number")

	senderCmd.PersistentFlags().StringVarP(&jobsCmdType, "type", "t", "mail", "jobs number")
}
