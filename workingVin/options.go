package main

import "os"

func getElse(variable string, value *string) {
	v := os.Getenv(variable)
	if len(v) > 0 {
		*value = v
	}
}
