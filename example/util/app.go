package util

import (
	"flag"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
)

type Args struct {
	ConfigFilePath *string
	Debug          *bool
	LogLevel       *string
}

func GetArgs() *Args {
	args := &Args{}
	flag.Usage = func() {
		fmt.Printf("usage: helper -config=develop.toml -debug=true")
		flag.PrintDefaults()
		os.Exit(2)
	}
	args.ConfigFilePath =
		flag.String("config", "develop.toml", "configuration file path")
	args.Debug =
		flag.Bool("debug", false, "use debug mode")
	args.LogLevel =
		flag.String("log", "debug", "log level")

	flag.Parse()

	if flag.NArg() > 0 {
		flag.Usage()
	}
	return args
}

func InitLogSetting(args *Args) {
	logLevel := *args.LogLevel
	if *args.Debug {
		logLevel = "debug"
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(err.Error())
	}
	logrus.SetLevel(level)
	//logrus.SetFormatter(&logrus.JSONFormatter{})
}
