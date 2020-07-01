package main

import (
	"encoding/json"
	"os"
)

type SinkerOptions struct {
	Brokers      string
	Protocol     string
	GroupID      string
	OffsetReset  string
	Topic        string
	InfluxServer string
	InfluxDB     string
	Measurement  string
	Retension    string
	Update       string
	Duration     string
	MaxRecords   string
	V2Vehicle    string
}

func getElse(variable string, value *string) {
	v := os.Getenv(variable)
	if len(v) > 0 {
		*value = v
	}
}

func GetSinkerOptions() SinkerOptions {
	options := SinkerOptions{Brokers: "localhost:9092", Protocol: "PLAINTEXT", GroupID: "GO-INFLUX-SINKER",
		OffsetReset: "latest", Topic: "raw_car_messages", InfluxServer: "localhost:8086",
		InfluxDB: "mqtt_messages",
		Measurement: "seqcounter", Retension: "oneMonth", Update: "500", Duration: "10", MaxRecords: "500",
		V2Vehicle: "000014",
	}
	getElse("KAFKA_BROKERS", &options.Brokers)
	getElse("KAFKA_PROTOCOL", &options.Protocol)
	getElse("KAFKA_GROUP_ID", &options.GroupID)
	getElse("KAFKA_OFFSET_RESET", &options.OffsetReset)
	getElse("KAFKA_LISTEN_TOPIC", &options.Topic)
	getElse("INFLUXDB_SERVER", &options.InfluxServer)
	getElse("INFLUXDB_DB", &options.InfluxDB)
	getElse("INFLUXDB_MEASUREMENT", &options.Measurement)
	getElse("INFLUXDB_RETENTION", &options.Retension)
	getElse("INFLUXDB_UPDATE", &options.Update)
	getElse("INFLUXDB_DURATION", &options.Duration)
	getElse("KAFKA_MAX_RECORDS", &options.MaxRecords)
	getElse("V2VEHICLES", &options.V2Vehicle)
	return options
}

func (option SinkerOptions) String() string {
	if content, err := json.Marshal(option); err == nil {
		return string(content)
	}
	return ""
}
