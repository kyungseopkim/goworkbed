package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// InnoWeekday weekday
type InnoWeekday time.Weekday

// MarshalJSON of InnoWeekday returns DayOfWeek
func (w InnoWeekday) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", RevDayOfWeek[time.Weekday(w)])), nil
}

// UnmarshalJSON of InnoWeekday return DayofWeek
func (w *InnoWeekday) UnmarshalJSON(b []byte) error {
	var dow string
	if err := json.Unmarshal(b, &dow); err != nil {
		return err
	}
	*w = InnoWeekday(DayOfWeek[dow])
	return nil
}

// InnoDate has Date Type
type InnoDate struct {
	Exist bool
	Year  int
	Month int
	Day   int
}

func (d InnoDate) String() string {
	if d.Exist {
		return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
	}
	return ""
}

// MarshalJSON write InnoDate Type
func (d InnoDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
}

// UnmarshalJSON reads InnoDate Type
func (d *InnoDate) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	reg, err := regexp.Compile(`(\d{4})-(\d{2})-(\d{2})`)
	if err != nil {
		return err
	}

	result := reg.FindStringSubmatch(s)
	d.Exist = true
	d.Year, _ = strconv.Atoi(result[1])
	d.Month, _ = strconv.Atoi(result[2])
	d.Day, _ = strconv.Atoi(result[3])
	return nil
}

// InnoTime represts Time info
type InnoTime struct {
	Exist  bool
	Hour   int
	Minute int
}

func (time InnoTime) String() string {
	if time.Exist {
		return fmt.Sprintf("%02d:%02d", time.Hour, time.Minute)
	}
	return ""
}

// MarshalJSON return json type
func (time InnoTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.String())), nil
}

// UnmarshalJSON returns InnoTime
func (time *InnoTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	reg, err := regexp.Compile(`(\d{2}):(\d{2})`)
	if err != nil {
		return err
	}
	result := reg.FindStringSubmatch(s)
	time.Exist = true
	time.Hour, _ = strconv.Atoi(result[1])
	time.Minute, _ = strconv.Atoi(result[2])

	return nil
}

// Compare return comparison
func (time InnoTime) Compare(other InnoTime) int {
	if time.Exist && !other.Exist {
		return 1
	}

	if !time.Exist && other.Exist {
		return -1
	}

	if time.Hour == other.Hour {
		return time.Hour - other.Hour
	}
	return time.Hour - other.Hour
}

func (time *InnoTime) delta(other InnoTime) InnoTime {
	var result InnoTime

	if !time.Exist {
		return InnoTime{Exist: false}
	}
	minutes := time.Minute + other.Minute
	deltaHour := minutes / 60
	result.Minute = minutes % 60

	hours := time.Hour + other.Hour + deltaHour
	adjust := hours % 24
	if adjust < 0 {
		adjust = 24 - adjust
	}
	result.Hour = adjust

	time.Hour = result.Hour
	time.Minute = result.Minute
	return result
}

// Floor return ...
func (time InnoTime) Floor(interval int) InnoTime {
	quotient := time.Minute / interval
	return InnoTime{true, time.Hour, quotient * interval}
}

// EnumerateTo returns times between start and end
func (time InnoTime) EnumerateTo(end InnoTime, interval int) []InnoTime {
	result := make([]InnoTime, 0)
	for current := time.Floor(interval); current.Compare(end.Floor(interval)) <= 0; current.delta(InnoTime{Minute: interval}) {
		result = append(result, current)
	}
	return result
}
