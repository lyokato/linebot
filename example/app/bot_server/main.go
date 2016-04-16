package main

import (
	"runtime"

	"github.com/lyokato/linebot"
	"github.com/lyokato/linebot/example/config"
	"github.com/lyokato/linebot/example/handler"
	"github.com/lyokato/linebot/example/util"

	"github.com/gin-gonic/gin"
	"github.com/thoas/stats"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	args := util.GetArgs()
	util.InitLogSetting(args)
	conf := config.LoadFromFile(*args.ConfigFilePath)
	runServer(conf, *args.Debug)
}

func runServer(conf *config.Config, debug bool) {

	bc := conf.Bot

	cw := linebot.NewClientWorker(bc.ChannelId,
		bc.ChannelSecret, bc.MID, bc.ClientWorkerQueueSize)
	cw.Run()
	evh := handler.New(cw)

	bot := linebot.NewServer()

	g := gin.Default()
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	g.POST("/", gin.WrapF(bot.HTTPHandler(bc.ChannelSecret, evh, bc.EventDispatcherQueueSize)))

	/*
		admin := g.Group("/admin")
		{
			admin.GET("/login", admin_controller.ShowLoginPage)
			admin.POST("/login", admin_controller.Login)
			admin.GET("/", admin_controller.Index)
		}
	*/

	s := stats.New()
	g.GET("/stats", func(c *gin.Context) {
		c.JSON(200, s.Data())
	})
	g.Use(func(c *gin.Context) {
		beginning, recorder := s.Begin(c.Writer)
		c.Next()
		s.End(beginning, recorder)
	})

	g.Run(conf.Web.Address())
}
