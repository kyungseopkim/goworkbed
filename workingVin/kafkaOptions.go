package main

import (
	"encoding/json"
	"log"
	"strconv"
)

type KafkaOptions struct {
	Brokers				string
	Protocol 			string
	GroupID				string
	OffsetReset			string
	Topic 				string
	MaxPartition		string
	MaxFetch 			string
}

func GetKafkaOptions() *KafkaOptions {
	options := KafkaOptions{Brokers: "localhost:9092", Protocol:"PLAINTEXT", GroupID:"GO-INFLUX-SINKER",
		OffsetReset: "latest", Topic:"signals", MaxPartition: "104857600", MaxFetch: "524288000"}

	getElse("KAFKA_BROKERS", &options.Brokers)
	getElse("KAFKA_PROTOCOL", &options.Protocol)
	getElse("KAFKA_GROUP_ID", &options.GroupID)
	getElse("KAFKA_OFFSET_RESET", &options.OffsetReset)
	getElse("KAFKA_LISTEN_TOPIC", &options.Topic)
	getElse("KAFKA_MAX_PARTITION", &options.MaxPartition)
	getElse( "KAFKA_MAX_FETCH", &options.MaxFetch)
	return &options
}

func (option KafkaOptions) String () string {
	if content, err := json.Marshal(option); err == nil {
		return string(content)
	}
	return ""
}

func str2int(valstr string) int {
	val, err := strconv.Atoi(valstr)
	if err != nil {
		log.Fatal(err)
	}
	return val

}

func (option KafkaOptions) GetMaxFetch() int {
	return str2int(option.MaxFetch)
}

func (option KafkaOptions) GetMaxPartition() int {
	return str2int(option.MaxPartition)
}

