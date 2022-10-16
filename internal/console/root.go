package console

import (
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/luckyAkbar/jatis-redis-training/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cobra-example",
	Short: "An example of cobra",
	Long: `This application shows how to create modern CLI
			applications in go using Cobra CLI library`,
}

func init() {
	setupLogger()
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func setupLogger() {
	formatter := runtime.Formatter{
		ChildFormatter: &logrus.JSONFormatter{},
		Line:           true,
		File:           true,
	}

	if config.Env() == "development" {
		formatter = runtime.Formatter{
			ChildFormatter: &logrus.TextFormatter{
				ForceColors:   true,
				FullTimestamp: true,
			},
			Line: true,
			File: true,
		}
	}

	logrus.SetFormatter(&formatter)
	logrus.SetOutput(os.Stdout)

	logLevel, err := logrus.ParseLevel(config.LogLevel())
	if err != nil {
		logLevel = logrus.DebugLevel
	}
	logrus.SetLevel(logLevel)

}
