package main

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "strings"
    "testing"
    "time"
)

func TestString(t *testing.T) {
    r := strings.NewReader("Hello1")
    b := make([]byte, 6)
    n, err := r.Read(b)
    if err != io.EOF {
        fmt.Printf("%v - %v", n, b)
    }
}

func TestGetByte(t *testing.T) {
    data := []byte { 6, 72, 101, 108, 108, 111, 49 }
    r := bytes.NewReader(data)
    n, _ := r.ReadByte()
    fmt.Println(uint8(n))
}

func TestPayloadFromKafkaMessage(t *testing.T) {
    dat, err := ioutil.ReadFile("/tmp/v2logs/1590105614.7820-43.baby")
    if err != nil {
        panic(err)
    }
    v2 := make(V2Vehicle)
    v2.FromS3()

    payload := GetPayload(dat, v2, time.Now().Unix(), "abcdef")
    fmt.Println(payload)
}