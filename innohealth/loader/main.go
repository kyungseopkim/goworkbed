package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/vjeantet/jodaTime"
)

func parseDate(input string) InnoDate {
	date, err := jodaTime.Parse("yyyyMMdd", input)
	if err != nil {
		log.Println(err)
		return InnoDate{Exist: false}
	}
	return InnoDate{true, date.Year(), int(date.Month()), date.Day()}
}

func parseTime(input string) InnoTime {
	time, err := jodaTime.Parse("HHmm", input)
	if err != nil {
		log.Println(err)
		return InnoTime{Exist: false}
	}
	return InnoTime{true, time.Hour(), time.Minute()}
}

func setValue2Operation(elem reflect.Value, field string, value interface{}) {
	n := FieldMap[field]
	f := elem.Elem().FieldByName(n)

	switch n {
	case "SID", "OperationKindID", "AnesthesiaID", "OperationRoom":
		sid, err := strconv.Atoi(value.(string))
		if err != nil {
			log.Panic(err)
		}
		f.SetInt(int64(sid))
	case "DayOfWeek":
		dow := DayOfWeek[value.(string)]
		f.Set(reflect.ValueOf(InnoWeekday(dow)))
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
	case "FirstDiogonasis":
		f.SetBool(value.(bool))
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
	log.SetFlags(log.Lshortfile)
	if len(os.Args) < 3 {
		fmt.Println("inputfile dbname")
		os.Exit(1)
	}

	dbname := os.Args[2]
	InitMgo()
	defer CloseMgo()

	fid, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(fid)
	fields, err := reader.Read()
	if err != nil {
		panic(err)
	}
	var data []*Operation
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
	UpdateOPeration(dbname, data)
}
