package main

import (
	"github.com/makerdao/vulcanizedb/cmd"

	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint: true,
	})
	file, err := os.OpenFile("vulcanizedb.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	cmd.Execute()
}
