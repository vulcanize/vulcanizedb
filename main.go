package main

import (
	"github.com/vulcanize/vulcanizedb/cmd"

	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("vulcanizedb.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	cmd.Execute()
}
