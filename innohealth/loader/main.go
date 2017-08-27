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
	if len(input) == 4 {
		clock, err := jodaTime.Parse("HHmm", input)
		if err != nil {
			log.Println(err)
			return InnoTime{Exist: false}
		}
		return InnoTime{true, clock.Hour(), clock.Minute()}
	}

	if len(input) == 3 {
		hour, err := strconv.Atoi(string(input[0]))
		if err != nil {
			log.Println(err)
			return InnoTime{Exist: false}
		}
		minute, err := strconv.Atoi(input[1:])
		if err != nil {
			log.Println(err)
			return InnoTime{Exist: false}
		}

		if hour < 0 || hour > 24 {
			return InnoTime{Exist: false}
		}
		if minute < 0 || minute > 60 {
			return InnoTime{Exist: false}
		}
		return InnoTime{true, hour, minute}
	}
	return InnoTime{Exist: false}
}

func setValue2Operation(elem reflect.Value, field string, value interface{}) {
	n := FieldMap[field]
	f := elem.Elem().FieldByName(n)

	switch n {
	case "SID", "OperationKindID", "AnesthesiaID":
		sid, err := strconv.Atoi(value.(string))
		if err != nil {
			log.Panic(err)
		}
		f.SetInt(int64(sid))
	case "DayOfWeek":
		dow := DayOfWeek[value.(string)]
		f.Set(reflect.ValueOf(InnoWeekday(dow)))
	case "ID", "DoctorID", "DoctorName", "DepartmentID", "DepartmentName", "OperationName",
		"DiagnosisKind", "OperationKindName", "AnesthesiaName", "WardContact", "OperationRoom":
		f.SetString(value.(string))
	case "FirstVist", "ApplicationDate", "HospitalizedDate", "ReservedDate", "OperationDate":
		date := parseDate(value.(string))
		f.Set(reflect.ValueOf(date))
	case "FrontArrivedTime", "RoomEntranceTime", "AnesthesiaStart", "AnesthesiaReady", "OperationStart",
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
	//fmt.Printf("%+v", operation)
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

func printOperation(dbname string, data []*Operation) {
	fmt.Println(dbname)
	for _, item := range data {
		//fmt.Printf("%+v\n", item)
		fmt.Println(item.RoomEntranceTime, item.RoomOutTime)
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
	//printOperation(dbname, data)
	UpdateOPeration(dbname, data)
}
