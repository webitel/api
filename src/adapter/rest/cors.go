package rest

import (
	"github.com/kataras/iris"
	"github.com/rs/cors"
)

func injectCORS(app *iris.Application) {
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Content-Type", "X-Key", "X-Access-Token"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		Debug:            false,
	}
	corsWrapper := cors.New(corsOptions).ServeHTTP

	app.WrapRouter(corsWrapper)
}
