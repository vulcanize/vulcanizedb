package main

import (
	"os"

	"github.com/spf13/viper"

	"github.com/vulcanize/vulcanizedb/cmd"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logfile := viper.GetString("logfile")
	if logfile != "" {
		file, err := os.OpenFile(logfile,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logrus.SetOutput(file)
		} else {
			logrus.SetOutput(os.Stdout)
			logrus.Info("Failed to log to file, using default stdout")
		}
	} else {
		logrus.SetOutput(os.Stdout)
	}
	cmd.Execute()
}
