package main

import (
	"encoding/json"
)

type KafkaOptions struct {
	Brokers				string
	Protocol 			string
	GroupID				string
	OffsetReset			string
	Topic 				string
}

func GetKafkaOptions() KafkaOptions {
	options := KafkaOptions{Brokers: "localhost:9092", Protocol:"PLAINTEXT", GroupID:"GO-INFLUX-SINKER",
		OffsetReset: "latest", Topic:"signals"}

	getElse("KAFKA_BROKERS", &options.Brokers)
	getElse("KAFKA_PROTOCOL", &options.Protocol)
	getElse("KAFKA_GROUP_ID", &options.GroupID)
	getElse("KAFKA_OFFSET_RESET", &options.OffsetReset)
	getElse("KAFKA_LISTEN_TOPIC", &options.Topic)
	return options
}

func (option KafkaOptions) String () string {
	if content, err := json.Marshal(option); err == nil {
		return string(content)
	}
	return ""
}
