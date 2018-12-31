package cmd

import (
	"fmt"
	"message_notification_practice/pb"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	log "gopkg.in/cihub/seelog.v2"
)

var clientCmdHost string

// serverCmd represents the server command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start up message notification client",
	Run:   clientProc,
}

func clientProc(cmd *cobra.Command, args []string) {

	log.Debug("Start clientProc")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, clientCmdHost, grpc.WithInsecure())
	handleInitError(err, "connect")
	defer conn.Close()

	//初始化客户端
	c := pb.NewMsgNotificationClient(conn)

	//调用方法
	reqBody := new(pb.MsgNotificationRequest)

	for i := 3; i < 4; i++ {
		reqBody.NoticeType = []pb.NoticeType{0, 1, 2}
		reqBody.Content = fmt.Sprintf("gRPC liujx: [%d]", i)
		reqBody.Group = uint64(i)
		r, err := c.CheckIn(context.Background(), reqBody)
		if err != nil {
			fmt.Println(err)
		}

		log.Debugf("[%d]: %s", i, r.Status)

	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Debug("Exit clientProc")
}

func init() {

	rootCmd.AddCommand(clientCmd)

	clientCmd.PersistentFlags().StringVarP(&clientCmdHost, "host", "s", "127.0.0.1:3000", "Server host address")

}
