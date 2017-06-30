package rest

import (
	. "../../config"
	"github.com/kataras/iris"
	"log"
)

type Rest struct {
}

func ListenAndServe() {

	app := iris.New()

	injectCORS(app)
	injectLogger(app)
	v2 := apiV2Router(app)

	injectAclV2(v2, app)
	injectCallmeV2(v2, app)

	injectErrors(app)
	if err := app.Run(iris.Addr(Config.Get("server:host") + ":" + Config.Get("server:port"))); err != nil {
		log.Println(err.Error())
	}
}
