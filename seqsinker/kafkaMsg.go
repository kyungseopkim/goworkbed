package main

import (
    "encoding/json"
    "log"
)

type KafkaMessage struct {
    UserName        string      `json:"username"`
    Timestamp       int64       `json:"ts"`
    Topic           string      `json:"topic"`
    Qos             int32       `json:"qos"`
    Payload         string      `json:"payload"`
    Node            string      `json:"node"`
    ClientId        string      `json:"clientid"`
}


func (kafka KafkaMessage) String () string {
    bstr, _ := json.Marshal(kafka)
    return string(bstr)
}

func KafkaMessageFromJson(msg []byte) *KafkaMessage {
    var kafkaMsg KafkaMessage
    err := json.Unmarshal(msg,&kafkaMsg)
    if err != nil {
        log.Println(err)
        return nil
    }
    return &kafkaMsg
}
