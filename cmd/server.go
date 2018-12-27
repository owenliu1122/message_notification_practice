package cmd

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice/controllers"
	"message_notification_practice/mq"
	"message_notification_practice/pb"
	"message_notification_practice/services"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var serverCmdPort int

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start up message notification server",
	Run:   serverProc,
}

func serverProc(cmd *cobra.Command, args []string) {

	log.Debugf("redis：%#v\n", viper.GetStringMap("server"))
	cfg := viper.GetStringMapString(cmd.Use)

	log.Debug("Start serverProc")
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	db, err := gorm.Open("mysql", cfg["mysql"])
	if err != nil {
		fmt.Printf("init mysql failed, err: %s", err)
		return
	}

	defer db.Close()

	time.Sleep(2 * time.Second) // TODO: remove is

	// grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverCmdPort))
	handleInitError(err, "net")

	gs := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time: 10 * time.Minute,
		}),
	)

	defer gs.GracefulStop()

	mqCli := mq.NewMq(cfg["rabbitmq"])
	if e := mqCli.InitConnection(); e != nil {
		log.Error("InitConnection failed, err: ", e)
	}
	defer mqCli.Close()

	if e := mqCli.InitProducer(cfg["mqexchange"], cfg["mqroutingkey"]); e != nil {
		log.Error("InitProducer failed, err: ", e)
	}

	ctl := controllers.NewServerController(services.NewNotificationService(db, mqCli))
	pb.RegisterMsgNotificationServer(gs, ctl)
	go gs.Serve(lis)

	log.Info("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Debug("Exit serverProc")
}

func init() {

	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().IntVarP(&serverCmdPort, "port", "p", 3000, "Port to listen")

}
