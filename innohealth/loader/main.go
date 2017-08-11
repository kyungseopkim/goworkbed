package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/vjeantet/jodaTime"
)

func parseDate(input string) time.Time {
	date, err := jodaTime.Parse("yyyyMMdd", input)
	if err != nil {
		log.Println(err)
		date = time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)
	}
	return date
}

func parseTime(input string) time.Time {
	ptime, err := jodaTime.Parse("HHmm", input)
	if err != nil {
		ptime, err = jodaTime.Parse("Hmm", input)
		if err != nil {
			log.Println(err)
			ptime = time.Date(1950, 1, 1, 0, 0, 0, 0, time.UTC)
		}
	}
	ptime = time.Date(1950, 1, 1, ptime.Hour(), ptime.Minute(), 0, 0, time.UTC)
	return ptime
}

func setValue2Operation(elem reflect.Value, field string, value interface{}) {
	n := FieldMap[field]
	f := elem.Elem().FieldByName(n)

	switch n {
	case "DayOfWeek":
		dow := DayOfWeek[value.(string)]
		f.Set(reflect.ValueOf(dow))
	case "SID", "OperationKindID", "AnesthesiaID", "OperationRoom":
		sid, err := strconv.Atoi(value.(string))
		if err != nil {
			log.Panic(err)
		}
		f.SetInt(int64(sid))
	case "ID", "DoctorID", "DoctorName", "DepartmentID", "DepartmentName", "OperationName",
		"DiagnosisKind", "OperationKindName", "AnesthesiaName", "WardContact":
		f.SetString(value.(string))
	case "FirstVist", "ApplicationDate", "HospitalizedDate", "ReservedDate", "OperationDate":
		date := parseDate(value.(string))
		f.Set(reflect.ValueOf(date))
	case "FrontArrivedTime", "RoomEtranceTime", "AnesthesiaStart", "AnesthesiaReady", "OperationStart",
		"OperationEnd", "AnesthesiaAwaken", "RoomOutTime":
		time := parseTime(value.(string))
		f.Set(reflect.ValueOf(time))
	}

}

func record2Operation(fields []string, record []string) *Operation {
	if len(fields) != len(record) {
		log.Fatal("Number of Fields is mismatched")
	}
	operation := new(Operation)
	for i, field := range fields {
		setValue2Operation(reflect.ValueOf(operation), strings.TrimSpace(field), record[i])
	}
	return operation
}

func csvLoader(filename string) {
	file := filename
	fileReader, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	reader := csv.NewReader(fileReader)
	// read header
	fields, err := reader.Read()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		record2Operation(fields, record)
	}
}

func main() {
	InitMgo("testdb")
	defer CloseMgo()

	fid, err := os.Open("/Users/kkim/workspace/innohealth/operation.csv")
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(fid)
	fields, err := reader.Read()
	if err != nil {
		panic(err)
	}

	data := make([]*Operation, 0)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		data = append(data, record2Operation(fields, record))
	}

	UpdateOPeration(data)
}
