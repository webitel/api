package rest

import (
	"../../models/acl"
	"../../services/auth"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
)

func injectAclV2(apiV2 router.Party, app router.Party) {
	aclRoute := apiV2.Party("/acl")
	{
		aclRoute.Get("/roles", listRoles)
		aclRoute.Get("/roles/{name}", itemRole)
		aclRoute.Put("/roles/{name}", updateRole)
		aclRoute.Delete("/roles/{name}", removeRole)
		aclRoute.Post("/roles", createRole)
	}
}

func listRoles(ctx context.Context) {
	err, data := acl.List(auth.GetSessionFromContext(ctx))
	if err != nil {
		ctx.StatusCode(err.Code)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	ctx.JSON(map[string]interface{}{
		"data":   data,
		"status": "OK",
	})
}

func itemRole(ctx context.Context) {
	err, data := acl.ListPerms(auth.GetSessionFromContext(ctx), ctx.Params().Get("name"))
	if err != nil {
		ctx.StatusCode(err.Code)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	ctx.JSON(map[string]interface{}{
		"data":   data,
		"status": "OK",
	})
}

func updateRole(ctx context.Context) {
	var body map[string][]string
	e := ctx.ReadJSON(&body)
	if e != nil {
		ctx.StatusCode(500)
		ctx.JSON(map[string]interface{}{
			"info":   e.Error(),
			"status": "error",
		})
		return
	}

	err := acl.UpdatePermission(auth.GetSessionFromContext(ctx), ctx.Params().Get("name"), body)
	if err != nil {
		ctx.StatusCode(err.Code)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	ctx.JSON(map[string]interface{}{
		"data":   body,
		"status": "OK",
	})
}

func removeRole(ctx context.Context) {
	err := acl.RemoveGroup(auth.GetSessionFromContext(ctx), ctx.Params().Get("name"))
	if err != nil {
		ctx.StatusCode(err.Code)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	ctx.JSON(map[string]interface{}{
		"info":   "Successful",
		"status": "OK",
	})
}

type roleJson struct {
	Name       string `json:"role"`
	ParentName string `json:"parent"`
}

func createRole(ctx context.Context) {
	var b roleJson
	ctx.ReadJSON(&b)
	err := acl.CreateGroup(auth.GetSessionFromContext(ctx), b.Name, b.ParentName)
	if err != nil {
		ctx.StatusCode(err.Code)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	ctx.JSON(map[string]interface{}{
		"info":   "Successful",
		"status": "OK",
	})
}
