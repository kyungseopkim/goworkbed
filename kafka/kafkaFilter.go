package main

import (
	"log"
	"strings"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	options := GetFilterOptions()
	log.Println(options)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": options.Brokers,
		"group.id":          options.GroupID,
		"auto.offset.reset": options.OffsetReset,
		"security.protocol": options.Protocol,
	})
	defer c.Close()
	if err != nil {
		panic(err)
	}

	producer, er := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": options.Brokers,
		"security.protocol": options.Protocol,
	})

	if er != nil {
		log.Fatal(er)
	}

	filter := make(map[string]bool)
	for _, v := range strings.Split(options.SelectedSignals, ",") {
		filter[v] = true
	}

	c.SubscribeTopics([]string{options.Topic}, nil)

	topic := options.TargetTopic

	// Delivery report handler for produced messages
	go func() {
		count := 0
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Fatalf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					if count == 100 {
						//log.Printf("message %s\n", ev.Value)
						count = 0
					}
					count++
				}
			}
		}
	}()

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			signal := SignalFromJson(msg.Value)
			if _, ok := filter[signal.SignalName]; ok {
				producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Value:          signal.JSON(),
					Headers:        []kafka.Header{{Key: "kafkaSelectedSignals", Value: []byte("header values are binary")}},
				}, nil)
			}
		} else {
			// The client will automatically try to recover from all errors.
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
