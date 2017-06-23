package rest

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
)

func injectLogger(app *iris.Application) {
	customLogger := logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
	})

	app.Use(customLogger)
}
