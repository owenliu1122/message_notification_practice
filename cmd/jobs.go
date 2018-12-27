package cmd

import (
	"fmt"
	"github.com/spf13/viper"
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

	cfg := viper.GetStringMapString(cmd.Use)
	cfgCM := viper.Sub(cmd.Use).GetStringMapString("consumer")
	cfgPCMap := viper.Sub(cmd.Use).GetStringMap("producer")

	log.Debugf("notification CFG: -->  %#v\n", cfg)
	log.Debugf("notification cfgCM: -->  %#v\n", cfgCM)
	log.Debugf("notification cfgPCMap: -->  %#v\n", cfgPCMap)

	db, err := gorm.Open("mysql", cfg["mysql"])
	if err != nil {
		fmt.Printf("init mysql failed, err: %s", err)
		return
	}

	defer db.Close()

	mqCli := mq.NewMq(cfg["rabbitmq"])
	if e := mqCli.InitConnection(); e != nil {
		log.Error("InitConnection failed, err: ", e)
	}
	defer mqCli.Close()

	if e := mqCli.InitProducer("", ""); e != nil {
		log.Error("InitProducer failed, err: ", e)
	}

	mqSendSvc := services.NewMqSendService(mqCli, services.NewGroupUserRelationService(db))

	for k, v := range cfgPCMap {
		log.Debugf("k: %s, v: %#v\n", k, v)
		mqSendSvc.RegisterExchangeRouting(k, mq.BaseProducer{
			Exchange:   v.(map[string]interface{})["mqexchange"].(string),
			RoutingKey: v.(map[string]interface{})["mqroutingkey"].(string),
		})
	}

	ctl := controllers.NewNotificationController(mqSendSvc)

	for i := 0; i < jobsCmdNum; i++ {

		mqCli.RegisterConsumer(
			fmt.Sprintf("notification[%d]", i),
			ctl.Handler,
			mq.BaseConsumer{
				Queue:   cfgCM["queue"],
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

	time.Sleep(500 * time.Millisecond)

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

	cfg := viper.GetStringMapString(cmd.Use)
	cfgCM := viper.Sub(cmd.Use).Sub(jobsCmdType).GetStringMapString("consumer")
	cfgSendSvc := viper.Sub(cmd.Use).Sub(jobsCmdType).GetStringMapString("sendsvc")

	mqCli := mq.NewMq(cfg["rabbitmq"])
	if e := mqCli.InitConnection(); e != nil {
		log.Error("InitConnection failed, err: ", e)
	}
	defer mqCli.Close()

	// TODO: 需要使用统一的接口，这里暂时时候 mail 接口测试
	sendSvc := services.NewSenderService(jobsCmdType, cfgSendSvc)
	ctl := controllers.NewSenderController(sendSvc)

	for i := 0; i < jobsCmdNum; i++ {

		mqCli.RegisterConsumer(
			fmt.Sprintf("notification[%d]", i),
			ctl.Handler,
			mq.BaseConsumer{
				Queue:   cfgCM["queue"],
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

	time.Sleep(500 * time.Millisecond)

	log.Debug("Exit Jobs Sender!")
}

func init() {
	rootCmd.AddCommand(jobsCmd)

	jobsCmd.AddCommand(notificationCmd)
	jobsCmd.AddCommand(senderCmd)

	jobsCmd.PersistentFlags().IntVarP(&jobsCmdNum, "number", "n", 3, "jobs number")

	senderCmd.PersistentFlags().StringVarP(&jobsCmdType, "type", "t", "mail", "jobs number")
}
