package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/controllers"
	"github.com/owenliu1122/notice/mq"
	"github.com/owenliu1122/notice/redis"
	"github.com/owenliu1122/notice/services"

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

	mqConnection, err := mq.NewConnection(cfg.Notification.RabbitMQ)
	if err != nil {
		log.Error("new rabbitmq connection failed, err: ", err)
		return
	}
	defer mqConnection.Close()

	producer, err := mq.NewProducer("jobs notification producer", mqConnection)
	if err != nil {
		log.Error("create producer failed, err: ", err)
	}
	defer producer.Close()

	mqSendSvc := services.NewMqSendService(producer, services.NewGroupUserRelationService(db, cache), cfg.Notification.Producer)

	ctl := controllers.NewNotificationController(mqSendSvc)

	consumer, err := mq.NewConsumer(ctx,
		"jobs notification consumer",
		cfg.Notification.Consumer.Queue,
		mqConnection,
		jobsCmdNum,
		true,
		ctl.Handler)

	if err != nil {
		log.Error("create consumer failed, err: ", err)
		return
	}
	defer consumer.Close()

	if e := consumer.Start(); e != nil {
		log.Error("start consumer failed, err: ", e)
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
	defer cancel()

	var cfg notice.Config

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("%s service Unmarshal configuration is failed, err: %s", cmd.Use, err.Error())
		return
	}

	mqConnection, err := mq.NewConnection(cfg.Sender.RabbitMQ)
	if err != nil {
		log.Error("new rabbitmq connection failed, err: ", err)
		return
	}
	defer mqConnection.Close()

	// TODO: 需要使用统一的接口，这里暂时时候 mail 接口测试
	sendSvc := services.NewSenderService(jobsCmdType, cfg.Sender.SendService)
	ctl := controllers.NewSenderController(sendSvc)

	consumer, err := mq.NewConsumer(ctx,
		"jobs sender consumer",
		cfg.Sender.Consumer[jobsCmdType].Queue,
		mqConnection,
		jobsCmdNum,
		true,
		ctl.Handler)

	if err != nil {
		log.Error("create consumer failed, err: ", err)
	}
	defer consumer.Close()

	if e := consumer.Start(); e != nil {
		log.Error("start consumer failed, err: ", e)
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
