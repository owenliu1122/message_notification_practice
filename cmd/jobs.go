package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"message_notification_practice"
	"message_notification_practice/controllers"
	"message_notification_practice/mq"
	"message_notification_practice/redis"
	"message_notification_practice/services"
	"os"
	"os/signal"
	"syscall"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "gopkg.in/cihub/seelog.v2"
)

var (
	jobsCmdNum  int
	jobsCmdType string
)

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
	defer cancel()

	var cfg notice.Config

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("%s service Unmarshal configuration is failed, err: %s", cmd.Use, err.Error())
		return
	}

	cache, err := redis.NewRedisCli(cfg.Notification.Redis, json.Marshal, json.Unmarshal)
	if err != nil {
		fmt.Printf("init redis failed, err: %s", err)
		return
	}

	db, err := gorm.Open("mysql", cfg.Notification.MySQL)
	if err != nil {
		fmt.Printf("init mysql failed, err: %s", err)
		return
	}

	defer db.Close()

	mqCli := mq.NewMq(cfg.Notification.RabbitMQ)
	if e := mqCli.InitConnection(); e != nil {
		log.Error("InitConnection failed, err: ", e)
	}
	defer mqCli.Close()

	if e := mqCli.InitProducer("", ""); e != nil {
		log.Errorf("InitProducer failed, err: %s\n", e.Error())
	}

	mqSendSvc := services.NewMqSendService(mqCli, services.NewGroupUserRelationService(db, cache))

	for k, v := range cfg.Notification.Producer {
		log.Debugf("k: %s, v: %#v\n", k, v)
		mqSendSvc.RegisterExchangeRouting(k, mq.BaseProducer{
			Exchange:   v.Exchange,
			RoutingKey: v.RoutingKey,
		})
	}

	ctl := controllers.NewNotificationController(mqSendSvc)

	for i := 0; i < jobsCmdNum; i++ {

		mqCli.RegisterConsumer(
			fmt.Sprintf("notification[%d]", i),
			ctl.Handler,
			mq.BaseConsumer{
				Queue:   cfg.Notification.Consumer.Queue,
				AutoAck: true,
			})
	}

	if e := mqCli.StartConsumer(ctx); e != nil {
		log.Error("InitConnection failed, err: ", e)
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

	ctx, cancel := context.WithCancel(context.Background())

	var cfg notice.Config

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("%s service Unmarshal configuration is failed, err: %s", cmd.Use, err.Error())
		return
	}

	mqCli := mq.NewMq(cfg.Sender.RabbitMQ)
	if e := mqCli.InitConnection(); e != nil {
		log.Error("InitConnection failed, err: ", e)
	}
	defer mqCli.Close()

	// TODO: 需要使用统一的接口，这里暂时时候 mail 接口测试
	sendSvc := services.NewSenderService(jobsCmdType, cfg.Sender.SendService)
	ctl := controllers.NewSenderController(sendSvc)

	for i := 0; i < jobsCmdNum; i++ {

		mqCli.RegisterConsumer(
			fmt.Sprintf("notification[%d]", i),
			ctl.Handler,
			mq.BaseConsumer{
				Queue:   cfg.Sender.Consumer[jobsCmdType].Queue,
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
