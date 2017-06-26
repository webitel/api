package rest

import (
	. "../../config"
	"../../services/auth"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/errors"
	"github.com/kataras/iris/core/router"
	"io/ioutil"
)

var hmacSampleSecret interface{}

func init() {
	var err error
	hmacSampleSecret, err = ioutil.ReadFile(Config.Get("application:auth:tokenSecretKey"))
	if err != nil {
		panic(err)
	}
}

type tokenKey struct {
	Token string `json:"access_token"`
	Key   string `json:"x_key"`
}

func apiV2Router(app *iris.Application) router.Party {
	return app.Party("/api/v2", func(ctx context.Context) {
		var key, token string
		token = ctx.GetHeader("x-access-token")
		key = ctx.GetHeader("x-key")

		if token == "" {
			token = ctx.FormValue("access_token")
			key = ctx.FormValue("x_key")

			if token == "" && ctx.Method() != "GET" && ctx.Method() != "HEAD" {
				var jsonData tokenKey
				if err := ctx.ReadJSON(&jsonData); err == nil && jsonData.Token != "" {
					token = jsonData.Token
					key = jsonData.Key
				}
			}
		}

		if token == "" {
			ctx.StatusCode(INVALID_TOKEN_ERROR.Code)
			ctx.JSON(INVALID_TOKEN_ERROR)
			return
		}

		c, err := decodeToken(token)
		if err != nil {
			ctx.StatusCode(INVALID_TOKEN_ERROR.Code)
			ctx.JSON(INVALID_TOKEN_ERROR)
			return
		}

		if c.Version == 2 && c.Type == "domain" {
			e, role := auth.AclFindDomain(c.Id, c.Domain)
			if e != nil {
				ctx.StatusCode(INVALID_TOKEN_ERROR.Code)
				ctx.JSON(INVALID_TOKEN_ERROR)
				return
			}
			c.RoleName = role
		} else if key != "" {
			e := auth.AclFindAuth(key, c)
			if e != nil {
				ctx.StatusCode(INVALID_TOKEN_ERROR.Code)
				ctx.JSON(INVALID_TOKEN_ERROR)
				return
			}
		} else {
			ctx.StatusCode(INVALID_TOKEN_ERROR.Code)
			ctx.JSON(INVALID_TOKEN_ERROR)
			return
		}

		ctx.Values().Set("user", c)
		ctx.Next()
	})
}

func decodeToken(data string) (c *auth.Session, err error) {
	var token *jwt.Token
	token, err = jwt.ParseWithClaims(data, &auth.Session{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret, nil
	})

	if err != nil {
		return
	}

	var ok bool
	if c, ok = token.Claims.(*auth.Session); ok && token.Valid {
		return
	} else {
		err = errors.New("Bad token")
	}
	return
}
