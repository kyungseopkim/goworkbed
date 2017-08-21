package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestOperation(t *testing.T) {
	InitMgo()
	defer CloseMgo()
	data, err := QueryOperation("testdb", 0, 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, item := range data {
		dat, _ := json.Marshal(item)
		t.Log(string(dat))
	}
}

func TestOperationDoctors(t *testing.T) {
	InitMgo()
	defer CloseMgo()

	OperationByDoctorStat("testdb")
}

func TestOperationByWeekday(t *testing.T) {
	InitMgo()
	defer CloseMgo()

	OperationByWeekdayStat("testdb")
}
