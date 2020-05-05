package main

import (
	"encoding/json"
	"log"
)

type Signal struct {
	MsgID		int32		`json:"msgId"`
	Timestamp 	int64		`json:"timestamp"`
	Epoch 		int32		`json:"epoch"`
	Usec 		int32		`json:"usec"`
	Vlan 		string		`json:"vlan"`
	Vin 		string		`json:"vin"`
	MsgName 	string		`json:"msgName"`
	SignalName	string		`json:"signalName"`
	Value 		float32		`json:"value"`
}

func SignalFromJson(data []byte) *Signal  {
	var signal Signal
	err := json.Unmarshal(data, &signal)
	if err != nil {
		log.Printf("json parse error %s", data)
		return nil
	}
	return &signal
}

func (signal Signal) JSON () []byte {
	if v, err := json.Marshal(signal); err != nil {
		log.Println(err)
		return nil
	} else {
		return v
	}
	return nil
}

func (signal Signal) String() string {
	return string(signal.JSON())
}