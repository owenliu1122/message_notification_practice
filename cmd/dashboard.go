package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	log "gopkg.in/cihub/seelog.v2"
	"message_notification_practice/controllers"
	"message_notification_practice/redis"
	"message_notification_practice/services"
	"os"
	"os/signal"
	"syscall"
)

// dashboardCmd represents the server command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Start up message notification dashboard server",
	Run: func(cmd *cobra.Command, args []string) {

		log.Debug("Start serverProc")
		cfg := viper.GetStringMapString(cmd.Use)

		cache, err := redis.NewRedisCli(cfg["redis"], json.Marshal, json.Unmarshal)
		if err != nil {
			fmt.Printf("init redis failed, err: %s", err)
			return
		}

		db, err := gorm.Open("mysql", cfg["mysql"])

		if err != nil {
			fmt.Printf("init mysql failed, err: %s", err)
			return
		}

		defer db.Close()

		e := echo.New()

		// Middleware
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		grpCtl := controllers.NewGroupController(services.NewGroupService(db))
		usrCtl := controllers.NewUserController(services.NewUserService(db, cache))
		gurCtl := controllers.NewGroupUserRelationController(services.NewGroupUserRelationService(db, cache))

		// Groups Routes
		e.GET("/dashboard/groups", grpCtl.List)
		e.POST("/dashboard/groups", grpCtl.Create)
		e.PUT("/dashboard/groups", grpCtl.Update)
		e.DELETE("/dashboard/groups", grpCtl.Delete)

		// Users Routes
		e.GET("/dashboard/users", usrCtl.List)
		e.POST("/dashboard/users", usrCtl.Create)
		e.PUT("/dashboard/users", usrCtl.Update)
		e.DELETE("/dashboard/users", usrCtl.Delete)

		// Group and User Relations Routes
		e.GET("/dashboard/group_user_relations", gurCtl.ListMembers)
		e.GET("/dashboard/group_user_relations/available_members", gurCtl.AvailableMembers)
		e.POST("/dashboard/group_user_relations", gurCtl.AddMembers)
		//e.PUT("/dashboard/users", gurCtl.Update)
		//e.PUT("/dashboard/users", gurCtl.Update)
		e.DELETE("/dashboard/group_user_relations", gurCtl.DeleteMembers)

		// Start server
		e.Logger.Fatal(e.Start(":8000"))

		log.Info("server started")

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
	},
}

func init() {

	rootCmd.AddCommand(dashboardCmd)
}
