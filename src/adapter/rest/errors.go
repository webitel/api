package rest

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

type HttpError struct {
	Status string `json:"status"`
	Code   int    `json:"internalCode"`
	Info   string `json:"info"`
}

func NewError(code int, info string) *HttpError {
	return &HttpError{
		Status: "error",
		Code:   code,
		Info:   info,
	}
}

var NO_FOUNR_ERROR = NewError(404, "No found")
var INVALID_TOKEN_ERROR = NewError(401, "Invalid token")
var EXPIRE_TOKEN_ERROR = NewError(401, "Token Expired")
var INVALID_USER_ERROR = NewError(401, "Invalid User")
var INVALID_TOKEN_KEY_ERROR = NewError(401, "Invalid Token or Key")

func injectErrors(app *iris.Application) {
	app.OnErrorCode(iris.StatusNotFound, func(ctx context.Context) {
		ctx.StatusCode(404)
		ctx.JSON(NO_FOUNR_ERROR)
	})
}
