package main

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/middleware"
)

func main() {
	InitMgo()
	defer CloseMgo()
	e := echo.New()

	// Middleware
	e.Use(middleware.Log())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	rid := 1
	for _, router := range routing {
		e.Router().Add(&echo.Route{Method: router.Method,
			Path:    router.Path,
			Handler: router.Handler}, rid)
		rid++
	}

	e.Run(fasthttp.New(":1323"))
}
