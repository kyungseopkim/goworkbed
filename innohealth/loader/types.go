package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// InnoDate is a type
type InnoDate string

// InnoTime is time type
type InnoTime string

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

// // innoDate has Date Type
// type innoDate struct {
// 	Year  int
// 	Month int
// 	Day   int
// }

// func (d innoDate) String() string {
// 	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
// }

// // MarshalJSON write innoDate Type
// func (d innoDate) MarshalJSON() ([]byte, error) {
// 	return []byte(fmt.Sprintf("\"%s\"", d.String())), nil
// }

// // UnmarshalJSON reads innoDate Type
// func (d *innoDate) UnmarshalJSON(b []byte) error {
// 	var s string
// 	if err := json.Unmarshal(b, &s); err != nil {
// 		return err
// 	}

// 	reg, err := regexp.Compile(`(\d{4})-(\d{2})-(\d{2})`)
// 	if err != nil {
// 		return err
// 	}

// 	result := reg.FindStringSubmatch(s)

// 	d.Year, err = strconv.Atoi(result[1])
// 	if err != nil {
// 		return err
// 	}
// 	d.Month, err = strconv.Atoi(result[2])
// 	if err != nil {
// 		return err
// 	}
// 	d.Day, err = strconv.Atoi(result[3])
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
