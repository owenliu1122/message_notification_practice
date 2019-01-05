package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fpay/foundation-go"

	"github.com/fpay/foundation-go/database"
	"github.com/fpay/foundation-go/job"

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

	if err := viper.Sub(cmd.Use).Unmarshal(&cfg); err != nil {
		fmt.Printf("%s service Unmarshal configuration is failed, err: %s", cmd.Use, err.Error())
		return
	}

	logger := log.NewLogger(cfg.Logger, os.Stdout)

	logger.Info("Start Jobs Notify!")

	cache, err := services.NewRedisCli(logger, cfg.Redis, json.Marshal, json.Unmarshal)
	if err != nil {
		fmt.Printf("init redis failed, redis: %#v, err: %s", cfg, err)
		return
	}

	db, err := database.NewDatabase(cfg.MySQL)
	if err != nil {
		fmt.Printf("init mysql failed, err: %s", err)
		return
	}
	defer db.Close()

	jobManager := job.NewJobManager(cfg.RabbitMQ)

	mqSendSvc := services.NewMqSendService(logger, jobManager, services.NewGroupService(logger, db, cache), cfg.Producer)

	ctl := controllers.NewNotificationController(logger, mqSendSvc)

	for i := 0; i < jobsCmdNum; i++ {
		go func(ctx context.Context, jobManager foundation.JobManager, queue string, handler foundation.JobHandler) {
			if err := jobManager.Do(ctx, queue, handler); err != nil {
				logger.Errorf("job manager DO() function return err: %s", err)
			}
		}(ctx, jobManager, cfg.Consumer.Queue, ctl.Handler)
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

	if err := viper.Sub(cmd.Use).Unmarshal(&cfg); err != nil {
		fmt.Printf("%s service Unmarshal configuration is failed, err: %s", cmd.Use, err.Error())
		return
	}

	logger := log.NewLogger(cfg.Logger, os.Stdout)

	logger.Info("Start Jobs Sender!")

	jobManager := job.NewJobManager(cfg.RabbitMQ)

	sendSvc := services.NewSenderService(logger, jobsCmdType, cfg.SendService, jobManager)
	ctl := controllers.NewSenderController(logger, cfg.RetryDelay, sendSvc)

	//jobManager.Do(ctx, cfg.Consumer[jobsCmdType].Queue, ctl.Handler)

	for i := 0; i < jobsCmdNum; i++ {
		go func(ctx context.Context, jobManager foundation.JobManager, queue string, handler foundation.JobHandler) {
			if err := jobManager.Do(ctx, queue, handler); err != nil {
				logger.Errorf("job manager DO() function return err: %s", err)
			}
		}(ctx, jobManager, cfg.Consumer[jobsCmdType].Queue, ctl.Handler)
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
