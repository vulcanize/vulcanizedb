package main

import (
	"github.com/vulcanize/vulcanizedb/cmd"

	"github.com/spf13/viper"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	tracer.Start(tracer.WithServiceName(viper.GetString("datadog.name")))

	cmd.Execute()

	defer tracer.Stop()
}
