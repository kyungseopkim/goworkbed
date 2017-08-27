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

func ExampleGetDistictValue() {
	InitMgo()
	defer CloseMgo()
	session := mgoSession.Clone()
	values, _ := GetDistictValue(session.DB("testdb"), "operationroom")
	fmt.Println(values)
	//Output: [1 2 3 5 6 7 8 9 10 11 12 13 15 16 17 18 19 0 20 21 22]
}

func ExampleOperationByTimStat() {
	InitMgo()
	defer CloseMgo()

	fmt.Println(OperationByTimStat("testdb"))
	//Output: -
}
