package main

import (
	"os"

	"github.com/hyperpilotio/snap-plugin-collector-kafka/kafka"
	"github.com/intelsdi-x/snap/control/plugin"
)

// plugin bootstrap
func main() {
	plugin.Start(
		kafka.Meta(),
		kafka.NewKafkaCollector(),
		os.Args[1],
	)
}
