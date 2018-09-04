package main

import (
	"github.com/vulcanize/vulcanizedb/cmd"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"github.com/spf13/viper"
)

func main() {
	tracer.Start(tracer.WithServiceName(viper.GetString("datadog.name")))

	cmd.Execute()

	defer tracer.Stop()
}
