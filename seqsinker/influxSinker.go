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

func InfluxSinker(influx client.Client, data []*KafkaMessage, options SinkerOptions) {
    conf := client.BatchPointsConfig{Precision: "ms",
        Database:         options.InfluxDB,
        RetentionPolicy:  options.Retension,
        WriteConsistency: "any",
    }

    v2 := make(V2Vehicle)
    v2.FromEnv()

    bp, err := client.NewBatchPoints(conf)
    if err != nil {
        log.Fatalln(err)
    }

    points := make([]*client.Point, 0)
    for _, msg := range data {
        payload := PayloadFromKafkaMessage(msg, v2)
        if payload == nil {
            continue
        }
        tags := make(map[string]string)

        tags["vin"] = payload.Vin
        tags["vlan"] = strconv.Itoa(int(payload.Vlan))
        landingSec:= payload.LandingTs / 1000
        landing := time.Unix(landingSec, 0)
        tags["landingTs"] = fmt.Sprintf("%4d-%02d-%02d", landing.Year(), landing.Month(), landing.Day())
        fields := make(map[string]interface{})
        fields["value"] = payload.Seq
        sec := payload.Ts / 1000000
        ns := (payload.Ts % 1000000) * 1000
        tz := time.Unix(sec, ns)

        point, er := client.NewPoint(options.Measurement, tags, fields, tz)
        if er != nil {
            log.Println(er)
        }
        points = append(points, point)
    }

    if len(points) > 0 {
        bp.AddPoints(points)
        log.Println(fmt.Sprintf("writing %d items", len(points)))
        err = influx.Write(bp)
        if err != nil {
            log.Println(err)
        }
    }
}

func main() {
    options := GetSinkerOptions()
    log.Println(options)

    c, err := kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers": options.Brokers,
        "group.id":          options.GroupID,
        "auto.offset.reset": options.OffsetReset,
        "security.protocol": options.Protocol,
        //"max.poll.records":  options.MaxRecords,
    })
    if err != nil {
        panic(err)
    }
    defer c.Close()


    influxUrl := fmt.Sprintf("http://%s", options.InfluxServer)
    influx, er := client.NewHTTPClient(client.HTTPConfig{
        Addr: influxUrl,
    })
    if er != nil {
        log.Println("Error creating InfluxDB Client: ", err.Error())
    }
    defer influx.Close()

    maxBuffer, er1 := strconv.Atoi(options.Update)
    if er1 != nil {
        log.Fatalln(err)
    }

    err = c.SubscribeTopics([]string{options.Topic}, nil)
    if err != nil {
        log.Fatalln(err)
    }

    limit, err := strconv.ParseInt(options.Duration, 10, 32)
    if err != nil {
        log.Fatal(err)
    }

    prevCheck := time.Now()
    prevSync := time.Now()

    buffer := make([]*KafkaMessage, 0, maxBuffer)
    for {
        msg, err := c.ReadMessage(-1)
        if err == nil {
            kafkaMsg := KafkaMessageFromJson(msg.Value)
            if kafkaMsg == nil {
                continue
            }
            buffer = append(buffer, kafkaMsg)
            current := time.Now()
            if len(buffer) == maxBuffer || current.Sub(prevSync).Seconds() > float64(limit) {
                bufferClone := make([]*KafkaMessage, len(buffer))
                copy(bufferClone, buffer)
                go InfluxSinker(influx, bufferClone, options)
                buffer = buffer[:0]
                now := time.Now()
                if now.Sub(prevCheck).Seconds() > 10 {
                    //log.Println("working:" + signal.String())
                    prevCheck = now
                }
                prevSync = current
            }
        } else {
            // The client will automatically try to recover from all errors.
            log.Printf("Consumer error: %v (%v)\n", err, msg)
        }
    }
}
