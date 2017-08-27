package main

import (
	"log"
	"net/http"
	"strconv"
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
	{echo.GET, "/operation/bydoctor", OperationByDoctorHandler},
	{echo.GET, "/operation/byweekday", OperationByWeekdayHandler},
	{echo.GET, "/operation/byroom", OperationByRoomHandler},
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
	dbname := c.Query("dbname")
	startQuery := c.Query("start")
	limitQuery := c.Query("limit")

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		limit = 10
	}

	start, err := strconv.Atoi(startQuery)
	if err != nil {
		start = 0
	}

	InitMgo()
	defer CloseMgo()

	operations, err := QueryOperation(dbname, start, limit)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusBadRequest)
	}

	operations["name"] = dbname
	return c.JSON(operations)
}

// OperationByDoctorHandler returns
func OperationByDoctorHandler(c echo.Context) error {
	dbname := c.Query("dbname")
	stat, err := OperationByDoctorStat(dbname)
	if err != nil {
		return err
	}
	return c.JSON(stat)
}

// OperationByWeekdayHandler returns weekday and department
func OperationByWeekdayHandler(c echo.Context) error {
	dbname := c.Query("dbname")
	stat, err := OperationByWeekdayStat(dbname)
	if err != nil {
		return err
	}
	return c.JSON(stat)
}

// OperationByRoomHandler returns ...
func OperationByRoomHandler(c echo.Context) error {
	dbname := c.Query("dbname")
	data, err := OperationByTimStat(dbname)
	c.Logger().Info(data)
	if err != nil {
		return err
	}

	return c.JSON(dbname)
}
