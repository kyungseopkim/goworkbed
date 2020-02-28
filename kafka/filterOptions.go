package main

import (
	"encoding/json"
	"os"
)

type FilterOptions struct {
	Brokers				string
	Protocol 			string
	GroupID				string
	OffsetReset			string
	Topic 				string
	SelectedSignals 	string
	TargetTopic 		string
}

func getElse(variable string, value *string) {
	v := os.Getenv(variable)
	if len(v) > 0 {
		*value = v
	}
}

func GetFilterOptions() FilterOptions {
	options := FilterOptions{Brokers:"localhost:9092", Protocol:"PLAINTEXT", GroupID:"filter",
		OffsetReset: "latest", Topic:"signals", SelectedSignals: "IBMU_ModBlkVolt_M42B07", TargetTopic: "selected"}
	getElse("KAFKA_BROKERS", &options.Brokers)
	getElse("KAFKA_PROTOCOL", &options.Protocol)
	getElse("KAFKA_GROUP_ID", &options.GroupID)
	getElse("KAFKA_OFFSET_RESET", &options.OffsetReset)
	getElse("KAFKA_LISTEN_TOPIC", &options.Topic)
	getElse("FILTER_SIGNALS", &options.SelectedSignals)
	getElse("KAFKA_PRODUCE_TOPIC", &options.TargetTopic)
	return options
}

func (option FilterOptions) String () string {
	if content, err := json.Marshal(option); err == nil {
		return string(content)
	}
	return ""
}
