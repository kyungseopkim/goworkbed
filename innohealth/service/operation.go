package main

import (
	"strings"
	"time"
	"unicode"
)

// Operation is a type of Operation data
type Operation struct {
	SID               int          `json:"SN"`
	ID                string       `json:"등록번호"`
	FirstVist         time.Time    `json:"첫외래일"`
	ApplicationDate   time.Time    `json:"수술예약신청일"`
	DiagnosisKind     string       `json:"초재진"`
	HospitalizedDate  time.Time    `json:"입원일자"`
	ReservedDate      time.Time    `json:"수술예약일"`
	OperationDate     time.Time    `json:"수술확정시행일"`
	DoctorID          string       `json:"집도의"`
	DoctorName        string       `json:"집도의명"`
	DepartmentID      string       `json:"집도과"`
	DepartmentName    string       `json:"집도과명"`
	OperationName     string       `json:"수술명"`
	DayOfWeek         time.Weekday `json:"요일"`
	OperationKindID   int          `json:"수술구분"`
	OperationKindName string       `json:"수술구분명"`
	AnesthesiaID      int          `json:"마취구분"`
	AnesthesiaName    string       `json:"마취구분명"`
	OperationRoom     int          `json:"수술방"`
	WardContact       string       `json:"병동연락"`
	FrontArrivedTime  time.Time    `json:"입구도착"`
	RoomEtranceTime   time.Time    `json:"수술방입실"`
	AnesthesiaStart   time.Time    `json:"마취시작"`
	AnesthesiaReady   time.Time    `json:"마취완료"`
	OperationStart    time.Time    `json:"수술시작"`
	OperationEnd      time.Time    `json:"수술종료"`
	AnesthesiaAwaken  time.Time    `json:"마취종료"`
	RoomOutTime       time.Time    `json:"환자퇴실"`
}

// FieldMap Names
var FieldMap = map[string]string{
	"SN":      "SID",
	"등록번호":    "ID",
	"첫외래일":    "FirstVist",
	"수술예약신청일": "ApplicationDate",
	"초재진":     "DiagnosisKind",
	"입원일자":    "HospitalizedDate",
	"수술예약일":   "ReservedDate",
	"수술확정시행일": "OperationDate",
	"집도의":     "DoctorID",
	"집도의명":    "DoctorName",
	"집도과":     "DepartmentID",
	"집도과명":    "DepartmentName",
	"수술명":     "OperationName",
	"요일":      "DayOfWeek",
	"수술구분":    "OperationKindID",
	"수술구분명":   "OperationKindName",
	"마취구분":    "AnesthesiaID",
	"마취구분명":   "AnesthesiaName",
	"수술방":     "OperationRoom",
	"병동연락":    "WardContact",
	"입구도착":    "FrontArrivedTime",
	"수술방입실":   "RoomEtranceTime",
	"마취시작":    "AnesthesiaStart",
	"마취완료":    "AnesthesiaReady",
	"수술시작":    "OperationStart",
	"수술종료":    "OperationEnd",
	"마취종료":    "AnesthesiaAwaken",
	"환자퇴실":    "RoomOutTime",
}

// ReversFields Struct to Fields
var ReversFields = map[string]string{
	"SID":               "SN",
	"ID":                "등록번호",
	"FirstVist":         "첫외래일",
	"ApplicationDate":   "수술예약신청일",
	"DiagnosisKind":     "초재진",
	"HospitalizedDate":  "입원일자",
	"ReservedDate":      "수술예약일",
	"OperationDate":     "수술확정시행일",
	"DoctorID":          "집도의",
	"DoctorName":        "집도의명",
	"DepartmentID":      "집도과",
	"DepartmentName":    "집도과명",
	"OperationName":     "수술명",
	"DayOfWeek":         "요일",
	"OperationKindID":   "수술구분",
	"OperationKindName": "수술구분명",
	"AnesthesiaID":      "마취구분",
	"AnesthesiaName":    "마취구분명",
	"OperationRoom":     "수술방",
	"WardContact":       "병동연락",
	"FrontArrivedTime":  "입구도착",
	"RoomEtranceTime":   "수술방입실",
	"AnesthesiaStart":   "마취시작",
	"AnesthesiaReady":   "마취완료",
	"OperationStart":    "수술시작",
	"OperationEnd":      "수술종료",
	"AnesthesiaAwaken":  "마취종료",
	"RoomOutTime":       "환자퇴실",
}

// DayOfWeek is
var DayOfWeek = map[string]time.Weekday{
	"일": time.Sunday,
	"월": time.Monday,
	"화": time.Tuesday,
	"수": time.Wednesday,
	"목": time.Thursday,
	"금": time.Friday,
	"토": time.Saturday,
}

// StripWhiteSpace removes white space in string
func StripWhiteSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
