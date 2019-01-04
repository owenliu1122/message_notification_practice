package cmd

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/fpay/foundation-go/database"

	"github.com/owenliu1122/notice"
	"github.com/owenliu1122/notice/controllers"
	"github.com/owenliu1122/notice/services"

	"github.com/fpay/foundation-go/log"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dashboardCmd represents the server command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Start up message notification dashboard server",
	RunE: func(cmd *cobra.Command, args []string) error {

		var cfg notice.DashboardConfig

		if err := viper.Sub(cmd.Use).Unmarshal(&cfg); err != nil {
			return errors.Wrapf(err, "%s service unmarshal configuration is failed, err: %s", cmd.Use)
		}

		logger := log.NewLogger(cfg.Logger, os.Stdout)

		logger.Info("Start serverProc")

		cache, err := services.NewRedisCli(logger, cfg.Redis, json.Marshal, json.Unmarshal)
		if err != nil {
			return errors.Wrap(err, "init redis failed")
		}

		db, err := database.NewDatabase(cfg.MySQL)
		if err != nil {
			return errors.Wrap(err, "init mysql failed")
		}

		defer db.Close()

		e := echo.New()

		// Middleware
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		grpCtl := controllers.NewGroupController(logger, services.NewGroupService(logger, db, cache))
		usrCtl := controllers.NewUserController(logger, services.NewUserService(logger, db, cache))

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
		e.GET("/dashboard/group_user_relations", grpCtl.ListMembers)
		e.GET("/dashboard/group_user_relations/available_members", grpCtl.AvailableMembers)
		e.POST("/dashboard/group_user_relations", grpCtl.AddMembers)
		//e.PUT("/dashboard/users", gurCtl.Update)
		//e.PUT("/dashboard/users", gurCtl.Update)
		e.DELETE("/dashboard/group_user_relations", grpCtl.DeleteMembers)

		// Start server
		e.Logger.Fatal(e.Start(":8000"))

		logger.Info("server started")

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		return nil
	},
}

func init() {

	rootCmd.AddCommand(dashboardCmd)
}
