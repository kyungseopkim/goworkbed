package main

import (
    "bytes"
    "encoding/base64"
    "encoding/binary"
    "encoding/json"
    "log"
)

type Payload struct {
    Vin             string      `json:"vin"`
    Ver             int32       `json:"ver"`
    Arxml           int32       `json:"arxml"`
    Seq             int32       `json:"seq"`
    Vlan            int32       `json:"vlan"`
    Ts              int64       `json:"ts"`
    LandingTs       int64       `json:"landing_ts"`
    Payload         string      `json:"payload"`
}

func (p Payload) String () string {
    bstr, _ := json.Marshal(p)
    return string(bstr)
}

func getByte(reader *bytes.Reader) uint8 {
    data, err := reader.ReadByte()
    if err != nil {
        log.Fatalln(err)
    }
    return uint8(data)
}

func GetVin(reader *bytes.Reader) string {
    len := getByte(reader)
    vin := make([]byte, int(len))
    n, err := reader.Read(vin)
    if err != nil {
        log.Fatalln(err)
    }
    if int(len) != n {
        log.Fatalln("length reading failure")
    }

    return string(vin)
}

func GetPacketVer(reader *bytes.Reader) int32 {
    return int32(getByte(reader))
}

func getInt(reader *bytes.Reader) uint32 {
    data := make([]byte, 4)
    _, err := reader.Read(data)
    if err != nil {
        log.Fatalln(err)
    }

    return binary.BigEndian.Uint32(data)
}

func GetArxml(reader *bytes.Reader) int32 {
    return int32(getInt(reader))
}

func GetSeq(reader *bytes.Reader) int32 {
    return int32(getByte(reader))
}

func GetVlan(reader *bytes.Reader) int32 {
    return int32(getByte(reader))
}

func GetTs(reader *bytes.Reader) int64 {
    ts := make([]byte, 8)
    _, err := reader.Read(ts)
    if err != nil {
        log.Fatalln(err)
    }
    return int64(binary.BigEndian.Uint64(ts))
}

func GetUsec(reader *bytes.Reader) int32 {
    return int32(getInt(reader))
}

func PayloadFromKafkaMessage(msg *KafkaMessage, v2 V2Vehicle) *Payload {
    bindata, err := base64.StdEncoding.DecodeString(msg.Payload)
    if err == nil {
        log.Println(err)
        return nil
    }
    reader := bytes.NewReader(bindata)
    vin := GetVin(reader)
    if ! v2.Contains(vin) {
        return nil
    }
    ver := GetPacketVer(reader)
    arxml := GetArxml(reader)
    seq := GetSeq(reader)
    vlan := GetVlan(reader)
    ts := GetTs(reader)
    usec := GetUsec(reader)
    ts = (ts * 1000000) + int64(usec) // microsecs
    return &Payload{vin, ver,arxml, seq, vlan, ts, msg.Timestamp,msg.Payload}
}