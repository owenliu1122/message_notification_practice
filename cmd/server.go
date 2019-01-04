package cmd

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fpay/foundation-go/database"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/controllers"
	"github.com/owenliu1122/notice/pb"
	"github.com/owenliu1122/notice/services"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/fpay/foundation-go/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var serverCmdPort int

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start up message notification server",
	Run:   serverProc,
}

func serverProc(cmd *cobra.Command, args []string) {

	var cfg notice.ServerConfig

	if err := viper.Sub(cmd.Use).Unmarshal(&cfg); err != nil {
		fmt.Printf("%s service Unmarshal configuration is failed, err: %s", cmd.Use, err.Error())
		return
	}

	logger := log.NewLogger(cfg.Logger, os.Stdout)

	logger.Info("Start serverProc")
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	db, err := database.NewDatabase(cfg.MySQL)
	if err != nil {
		fmt.Printf("init mysql failed, err: %s", err)
		return
	}

	defer db.Close()

	// grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverCmdPort))
	handleInitError(err, "net")

	gs := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 10 * time.Minute,
		}),
	)

	defer gs.GracefulStop()

	mqConnection, err := services.NewMQConnection(cfg.RabbitMQ)
	if err != nil {
		logger.Error("new rabbitmq connection failed, err: ", err)
		return
	}
	defer mqConnection.Close()

	producer, err := services.NewProducer("server producer", mqConnection)
	if err != nil {
		logger.Error("NewProducer failed, err: ", err)
		return
	}
	defer producer.Close()

	svc := services.NewNotificationService(logger, db,
		producer,
		cfg.Producer.Exchange,
		cfg.Producer.RoutingKey)

	ctl := controllers.NewServerController(logger, svc)
	pb.RegisterMsgNotificationServer(gs, ctl)
	go gs.Serve(lis)

	logger.Info("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Exit serverProc")
}

func init() {

	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().IntVarP(&serverCmdPort, "port", "p", 3000, "Port to listen")

}
