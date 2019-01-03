package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/controllers"
	"github.com/owenliu1122/notice/services"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	cache, err := services.NewRedisCli(cfg.Notification.Redis, json.Marshal, json.Unmarshal)
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

	mqConnection, err := services.NewMQConnection(cfg.Notification.RabbitMQ)
	if err != nil {
		log.Error("new rabbitmq connection failed, err: ", err)
		return
	}
	defer mqConnection.Close()

	producer, err := services.NewProducer("jobs notification producer", mqConnection)
	if err != nil {
		log.Error("create producer failed, err: ", err)
	}
	defer producer.Close()

	mqSendSvc := services.NewMqSendService(producer, services.NewGroupService(db, cache), cfg.Notification.Producer)

	ctl := controllers.NewNotificationController(mqSendSvc)

	consumer, err := services.NewConsumer(ctx,
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

	mqConnection, err := services.NewMQConnection(cfg.Sender.RabbitMQ)
	if err != nil {
		log.Error("new rabbitmq connection failed, err: ", err)
		return
	}
	defer mqConnection.Close()

	producer, err := services.NewProducer("jobs sender producer", mqConnection)
	if err != nil {
		log.Error("create producer failed, err: ", err)
	}
	defer producer.Close()

	producer.DeclareExpiration(cfg.Sender.RetryProducer[jobsCmdType].Exchange,
		cfg.Sender.RetryProducer[jobsCmdType].RoutingKey,
		cfg.Sender.DelayProducer[jobsCmdType].Exchange,
		cfg.Sender.DelayProducer[jobsCmdType].RoutingKey,
		time.Duration(cfg.Sender.RetryDelay)*time.Second)

	sendSvc := services.NewSenderService(jobsCmdType, cfg.Sender.SendService, producer, cfg.Sender.RetryProducer[jobsCmdType])
	ctl := controllers.NewSenderController(sendSvc)

	consumer, err := services.NewConsumer(ctx,
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
