package main

import (
	"time"

	"github.com/webx-top/echo"
)

type router struct {
	Method  string
	Path    string
	Handler echo.HandlerFunc
}

var routing = []router{
	{echo.GET, "/ping", HealthHandler},
	{echo.GET, "/operation", OperationHandler},
}

// HealthHandler is checking server health checking
func HealthHandler(c echo.Context) error {
	message := make(map[string]interface{}, 0)
	message["status"] = "OK"
	message["date"] = time.Now().Format(time.RFC3339)
	return c.JSON(message)
}

// OperationHandler returns Operation related information
func OperationHandler(c echo.Context) error {
	return c.JSON("hello")
}
