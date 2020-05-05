package main

import (
	"fmt"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"strconv"
	"time"
)

func InfluxSinker(influx client.Client, data []*Signal, options InfluxDBOptions)  {
	conf := client.BatchPointsConfig{Precision:"ms",
		Database: options.InfluxDB,
		RetentionPolicy: options.Retention,
		WriteConsistency: "any",
	}

	bp, err := client.NewBatchPoints(conf)
	if err != nil {
		log.Fatalln(err)
	}

	points := make([]*client.Point, 0)
	for _, signal := range data {
		tags := make(map[string]string)
		tags["vin"] = signal.Vin
		tags["signalName"] = signal.SignalName
		event := time.Unix(int64(signal.Epoch), 0)
		loc, _ := time.LoadLocation("America/Los_Angeles")
		pst:= event.In(loc)
		tags["eventDay"] = pst.Format(time.RFC3339)[:10]
		fields := make(map[string]interface{})
		fields["value"] = signal.Value
		tz := time.Now()
		point, er := client.NewPoint(options.Measurement, tags, fields, tz)
		if er != nil {
			log.Println(er)
		}
		points = append(points, point)
	}

	bp.AddPoints(points)
	err = influx.Write(bp)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	kafkaOptions := GetKafkaOptions()
	influxOptions := GetInfluxDBOptions()
	log.Println(kafkaOptions)
	log.Println(influxOptions)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaOptions.Brokers,
		"group.id":          kafkaOptions.GroupID,
		"auto.offset.reset": kafkaOptions.OffsetReset,
		"security.protocol": kafkaOptions.Protocol,
	})

	defer c.Close()
	if err != nil {
		panic(err)
	}

	influxUrl := fmt.Sprintf("http://%s", influxOptions.InfluxServer)
	influx, er := client.NewHTTPClient(client.HTTPConfig{
		Addr: influxUrl,
	})
	if er != nil {
		log.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer influx.Close()

	maxBuffer, er1 := strconv.Atoi(influxOptions.Update)
	if er1 != nil {
		log.Fatalln(err)
	}

	err = c.SubscribeTopics([]string{kafkaOptions.Topic}, nil)
	if err != nil {
		log.Fatalln(err)
	}

	buffer := make([]*Signal, 0, maxBuffer)
	prevTime := time.Now()
	prevCheck := time.Now()
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			signal := SignalFromJson(msg.Value)
			buffer = append(buffer, signal)
			now := time.Now()
			if len(buffer) == maxBuffer || now.Sub(prevTime).Seconds() > influxOptions.getBatch() {
				bufferClone := make([]*Signal, len(buffer))
				copy(bufferClone, buffer)
				go InfluxSinker(influx, bufferClone, influxOptions)
				buffer = buffer[:0]
				nowCheck := time.Now()
				if nowCheck.Sub(prevCheck) > 5 {
					log.Println("working :" + signal.String())
					prevCheck= nowCheck
				}
				prevTime = now
			}
		} else {
			// The client will automatically try to recover from all errors.
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
