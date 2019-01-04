package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fpay/foundation-go/database"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/controllers"
	"github.com/owenliu1122/notice/services"

	"github.com/fpay/foundation-go/log"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg notice.NotificationConfig

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("%s service Unmarshal configuration is failed, err: %s", cmd.Use, err.Error())
		return
	}

	logger := log.NewLogger(cfg.Logger, os.Stdout)

	logger.Info("Start Jobs Notify!")

	cache, err := services.NewRedisCli(logger, cfg.Redis, json.Marshal, json.Unmarshal)
	if err != nil {
		fmt.Printf("init redis failed, err: %s", err)
		return
	}

	db, err := database.NewDatabase(cfg.MySQL)
	if err != nil {
		fmt.Printf("init mysql failed, err: %s", err)
		return
	}
	defer db.Close()

	mqConnection, err := services.NewMQConnection(cfg.RabbitMQ)
	if err != nil {
		logger.Error("new rabbitmq connection failed, err: ", err)
		return
	}
	defer mqConnection.Close()

	producer, err := services.NewProducer("jobs notification producer", mqConnection)
	if err != nil {
		logger.Error("create producer failed, err: ", err)
	}
	defer producer.Close()

	mqSendSvc := services.NewMqSendService(logger, producer, services.NewGroupService(logger, db, cache), cfg.Producer)

	ctl := controllers.NewNotificationController(logger, mqSendSvc)

	consumer, err := services.NewConsumer(ctx,
		"jobs notification consumer",
		cfg.Consumer.Queue,
		mqConnection,
		jobsCmdNum,
		true,
		ctl.Handler)

	if err != nil {
		logger.Error("create consumer failed, err: ", err)
		return
	}
	defer consumer.Close()

	if e := consumer.Start(); e != nil {
		logger.Error("start consumer failed, err: ", e)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Exit Jobs Notification!")
}

var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "Start sender job",
	Run:   senderProc,
}

func senderProc(cmd *cobra.Command, args []string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg notice.SenderConfig

	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Printf("%s service Unmarshal configuration is failed, err: %s", cmd.Use, err.Error())
		return
	}

	logger := log.NewLogger(cfg.Logger, os.Stdout)

	logger.Info("Start Jobs Sender!")

	mqConnection, err := services.NewMQConnection(cfg.RabbitMQ)
	if err != nil {
		logger.Error("new rabbitmq connection failed, err: ", err)
		return
	}
	defer mqConnection.Close()

	producer, err := services.NewProducer("jobs sender producer", mqConnection)
	if err != nil {
		logger.Error("create producer failed, err: ", err)
		return
	}
	defer producer.Close()

	if err = producer.DeclareExpiration(cfg.RetryProducer[jobsCmdType].Exchange,
		cfg.RetryProducer[jobsCmdType].RoutingKey,
		cfg.DelayProducer[jobsCmdType].Exchange,
		cfg.DelayProducer[jobsCmdType].RoutingKey,
		time.Duration(cfg.RetryDelay)*time.Second); err != nil {
		logger.Error("declare delay queue producer failed, err: ", err)
		return

	}

	sendSvc := services.NewSenderService(logger, jobsCmdType, cfg.SendService, producer, cfg.RetryProducer[jobsCmdType])
	ctl := controllers.NewSenderController(logger, sendSvc)

	consumer, err := services.NewConsumer(ctx,
		"jobs sender consumer",
		cfg.Consumer[jobsCmdType].Queue,
		mqConnection,
		jobsCmdNum,
		true,
		ctl.Handler)

	if err != nil {
		logger.Error("create consumer failed, err: ", err)
		return
	}
	defer consumer.Close()

	if e := consumer.Start(); e != nil {
		logger.Error("start consumer failed, err: ", e)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Exit Jobs Sender!")
}

func init() {
	rootCmd.AddCommand(jobsCmd)

	jobsCmd.AddCommand(notificationCmd)
	jobsCmd.AddCommand(senderCmd)

	jobsCmd.PersistentFlags().IntVarP(&jobsCmdNum, "number", "n", 3, "jobs number")

	senderCmd.PersistentFlags().StringVarP(&jobsCmdType, "type", "t", "mail", "jobs number")
}
