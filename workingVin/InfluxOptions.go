package main

import (
	"encoding/json"
	"strconv"
)

type InfluxDBOptions struct {
	InfluxServer string
	InfluxDB     string
	Measurement  string
	Retention    string
	Update       string
	Batch 		 string
}

func GetInfluxDBOptions() InfluxDBOptions {
	options := InfluxDBOptions{ InfluxServer: "localhost:8086", InfluxDB: "kafka_signal",
		Measurement: "working", Retention: "oneDay", Update: "500", Batch: "5",
	}
	getElse("INFLUXDB_SERVER", &options.InfluxServer)
	getElse("INFLUXDB_DB", &options.InfluxDB)
	getElse("INFLUXDB_MEASUREMENT", &options.Measurement)
	getElse("INFLUXDB_RETENTION", &options.Retention)
	getElse("INFLUXDB_UPDATE", &options.Update)
	getElse("INFLUXDB_BATCH", &options.Batch)
	return options
}

func (option InfluxDBOptions) String () string {
	if content, err := json.Marshal(option); err == nil {
		return string(content)
	}
	return ""
}

func (option InfluxDBOptions) getBatch() float64 {
	if batch, err := strconv.ParseFloat(option.Batch, 32); err == nil {
		return batch
	}
	return 0.0
}