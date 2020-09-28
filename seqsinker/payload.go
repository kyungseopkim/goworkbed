package main

import (
    "bytes"
    "encoding/base64"
    "encoding/binary"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/pierrec/lz4"
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

func getByte(reader *bytes.Reader) (uint8, error) {
    dat, err := reader.ReadByte()
    if err != nil {
        return 0, err
    }
    return uint8(dat), nil
}

func getUShort(reader *bytes.Reader) (uint16, error) {
    data := make([]byte, 2)
    _, err := reader.Read(data)
    if err != nil {
        return 0, err
    }
    return binary.BigEndian.Uint16(data), nil
}

func GetVin(reader *bytes.Reader) (string, error) {
    len, err := getByte(reader)
    if err != nil{
        return "", err
    }
    vin := make([]byte, int(len))
    n, err := reader.Read(vin)
    if err != nil {
        return "", err
    }
    if int(len) != n {
        log.Println(fmt.Sprintf("length reading failure %d != %d", int(len), n))
        return "", errors.New("reading error")
    }

    return string(vin), nil
}

func GetPacketVer(reader *bytes.Reader) (int32, error) {
    val , err := getByte(reader)
    if err != nil {
        return 0, err
    }
    return int32(val), nil
}

func getInt(reader *bytes.Reader) (uint32, error) {
    data := make([]byte, 4)
    _, err := reader.Read(data)
    if err != nil {
        return 0, err
    }

    return binary.BigEndian.Uint32(data), nil
}

func GetArxml(reader *bytes.Reader) (int32, error) {
    val, err := getInt(reader)
    if err != nil {
        return 0, err
    }
    return int32(val), nil
}

func GetSeq(ver int32, reader *bytes.Reader) (int32, error) {
    switch ver {
    case 2:
        seq, err := getByte(reader)
        if err != nil {
            return 0, err
        }
        return int32(seq), nil
    default:
        seq, err := getUShort(reader)
        if err != nil {
            return 0, err
        }
        return int32(seq), nil
    }
}

func GetVlan(ver int32, reader *bytes.Reader) (int32, error) {
    if ver > 2 { return 0, nil}

    val, err := getByte(reader)
    if err != nil {
        return 0, err
    }
    return int32(val), nil
}

func GetTs(reader *bytes.Reader) (int64, error) {
    ts := make([]byte, 8)
    _, err := reader.Read(ts)
    if err != nil {
        return 0, nil
    }
    return int64(binary.BigEndian.Uint64(ts)), nil
}

func GetUsec(reader *bytes.Reader) (int32, error) {
    val, err := getInt(reader)
    if err != nil {
        return 0, err
    }
    return int32(val), nil
}

func decompress(content []byte) ([]byte, error) {
    bucket := make([]byte, 10*1024*1024)
    size, err := lz4.UncompressBlock(content, bucket )
    if err != nil {
        log.Println(err)
        return nil, err
    }
    return bucket[:size], nil
}

func GetPayload(data []byte, v2 V2Vehicle,landing int64, payload string) *Payload {
    content, err := decompress(data)
    if err != nil {
        return nil
    }
    reader := bytes.NewReader(content)
    vin, err := GetVin(reader)
    if err != nil {
        log.Println(err)
        return nil
    }
    if ! v2.Contains(vin) {
        return nil
    }
    ver, err := GetPacketVer(reader)
    if err != nil {
        log.Println(err)
        return nil
    }

    if ver < 2 {
      log.Println(fmt.Sprintf("packet ver == %d", ver ))
      return nil
    }

    arxml, err := GetArxml(reader)
    if err != nil {
        log.Println(err)
        return nil
    }
    seq, err := GetSeq(ver, reader)
    if err != nil {
        log.Println(err)
        return nil
    }
    vlan, err := GetVlan(ver, reader)
    if err != nil {
        log.Println(err)
        return nil
    }
    ts, err := GetTs(reader)
    if err != nil {
        log.Println(err)
        return nil
    }
    usec, err := GetUsec(reader)
    if err != nil {
        log.Println(err)
        return nil
    }
    ts = (ts * 1000000) + int64(usec) // microsecs
    return &Payload{vin, ver,arxml, seq, vlan, ts, landing, payload}
}

func PayloadFromKafkaMessage(msg *KafkaMessage, v2 V2Vehicle) *Payload {
    bindata, err := base64.StdEncoding.DecodeString(msg.Payload)
    if err != nil {
        log.Println(err)
        return nil
    }
    return GetPayload(bindata, v2, msg.Timestamp, msg.Payload)
}
