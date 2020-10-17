package main

import (
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"github.com/influxdata/influxdb-client-go/v2"
	"log"
	"strconv"
	"time"
	//"github.com/panjf2000/ants/v2"
)

func InfluxSinker(influx influxdb2.Client, data []*Signal, options InfluxDBOptions)  {
	bucket := fmt.Sprintf("%s/%s", options.InfluxDB, options.Retention)
	writeAPI := influx.WriteAPI("",bucket)

	vins := make(map[string]map[string]int32)
	for _, signal := range data {
		event := time.Unix(int64(signal.Epoch), 0)
		loc, _ := time.LoadLocation("America/Los_Angeles")
		pst:= event.In(loc)
		eventDay := pst.Format(time.RFC3339)[:10]

		vin, ok := vins[signal.Vin]
		if ok {
			count, ok1 := vin[eventDay]
			if ok1 {
				vin[eventDay] = count + 1
			} else {
				vin[eventDay] = 1
			}
		} else {
			events := make(map[string]int32)
			events[eventDay]= 1
			vins[signal.Vin] = events
		}
	}

	for vin, events := range vins {
		for event, count := range events {
			tags := make(map[string]string)
			tags["vin"] = vin
			tags["eventDay"] = event
			fields := make(map[string]interface{})
			fields["count"] = count
			tz := time.Now()
			point := influxdb2.NewPoint(options.Measurement, tags, fields, tz)
			writeAPI.WritePoint(point)
		}
	}

	writeAPI.Flush()
	log.Println(fmt.Sprintf("%d - Okay", len(data)))

}

func main() {
	kafkaOptions := GetKafkaOptions()
	influxOptions := GetInfluxDBOptions()
	log.Println(kafkaOptions)
	log.Println(influxOptions)
	//defer ants.Release()

	//pool, _ := ants.NewPoolWithFunc(10, func(args interface{}) {
	//	params := args.(map[string]interface{})
	//	client := params["client"].(influxdb2.Client)
	//	data := params["data"].([]*Signal)
	//	options := params["options"].(InfluxDBOptions)
	//	InfluxSinker(client, data, options)
	//})

	//defer pool.Release()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaOptions.Brokers,
		"group.id":          kafkaOptions.GroupID,
		"auto.offset.reset": kafkaOptions.OffsetReset,
		"security.protocol": kafkaOptions.Protocol,
	})

	if err != nil {
		panic(err)
	}
	defer c.Close()


	influxUrl := fmt.Sprintf("http://%s", influxOptions.InfluxServer)

	maxBuffer, er1 := strconv.Atoi(influxOptions.Update)
	if er1 != nil {
		log.Fatalln(err)
	}
	client := influxdb2.NewClientWithOptions(influxUrl, "my-token",
		influxdb2.DefaultOptions().SetBatchSize(uint(maxBuffer)))
	defer client.Close()

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
				//bufferClone := make([]*Signal, len(buffer))
				//copy(bufferClone, buffer)
				InfluxSinker(client, buffer, influxOptions)
				//go InfluxSinker(client, bufferClone, influxOptions)
				//params := make(map[string]interface{})
				//params["client"] = client
				//params["data"] = bufferClone
				//params["options"] = influxOptions
				//pool.Invoke(params)
				buffer = buffer[:0]
				nowCheck := time.Now()
				if nowCheck.Sub(prevCheck) > 5 {
					//log.Println("working :" + signal.String())
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
