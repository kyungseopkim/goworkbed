package main

import (
    "os"
    "strings"
)

type V2Vehicle map[string]struct{}

var exists = struct{}{}

func (v2 V2Vehicle) String() string {
    return v2.String()
}

func (v2 V2Vehicle) Contains(key string) bool {
    _, ok := v2[key]
    if ok {
        return true
    }
    return false
}

func (v2 V2Vehicle) FromString(cars string ) {
    for _, item := range strings.Split(cars, ",") {
        v2[item]=exists
    }
}

func (v2 V2Vehicle) FromEnv() {
    values := os.Getenv("V2VEHICLES")
    v2.FromString(values)
}

