package main

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"strings"
)

func main()  {
	options := GetFilterOptions()
	log.Println(options)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": options.Brokers,
		"group.id":          options.GroupID,
		"auto.offset.reset": options.OffsetReset,
		"security.protocol": options.Protocol,
	})
	defer c.Close()

	p, er := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": options.Brokers,
		"security.protocol": options.Protocol,
		})

	if er != nil {
		log.Fatal(er)
	}

	defer p.Close()

	if err != nil {
		panic(err)
	}

	filter := make(map[string]bool)
	for _, v := range strings.Split(options.SelectedSignals, ",") {
		filter[v] = true
	}

	c.SubscribeTopics([]string{ options.Topic }, nil)

	topic := options.TargetTopic
	counter := 0
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			 signal := SignalFromJson(msg.Value)
			 if _ , ok := filter[signal.SignalName]; ok {
					p.Produce(&kafka.Message{
						TopicPartition: kafka.TopicPartition{Topic:&topic , Partition: kafka.PartitionAny},
						Value:          signal.JSON(),
					},nil)
					if counter == 10000 {
						log.Println(signal)
						counter = 0
					}
					counter++
			 }
		} else {
			// The client will automatically try to recover from all errors.
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}
