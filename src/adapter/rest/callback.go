package rest

import (
	"../../db"
	"../../services/auth"
	"../../services/callback"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
)

func injectCallmeV2(api router.Party, app router.Party) {
	// Public api
	app.Post("/api/v2/callback/{queueId}/members", memberCreate)
	callback := api.Party("/callback")
	{
		callback.Get("/", queueList)               //+
		callback.Post("/", queueCreate)            //+
		callback.Get("/{queueId}", queueItem)      //+
		callback.Put("/{queueId}", queueUpdate)    //+
		callback.Delete("/{queueId}", queueDelete) //+

		callback.Get("/{queueId}/members", membersList) //+
		//callback.Post("/{queueId}/members", memberCreate)
		callback.Get("/{queueId}/members/{memberId}", memberItem)                 //+
		callback.Put("/{queueId}/members/{memberId}", memberUpdate)               //+
		callback.Delete("/{queueId}/members/{memberId}", memberDelete)            //+
		callback.Post("/{queueId}/members/{memberId}/comments", memberAddComment) //+
	}
}

func queueList(ctx context.Context) {

	data, err := callback.CallbackQueueList(auth.GetSessionFromContext(ctx), db.RequestListFromUri(ctx))
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

func queueItem(ctx context.Context) {
	data, err := callback.CallbackQueueItem(ctx.Params().Get("queueId"))
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

func queueCreate(ctx context.Context) {
	var q *callback.Queue
	err := ctx.ReadJSON(&q)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	cError := callback.CallbackQueueCreate(q)

	if cError != nil {
		ctx.StatusCode(cError.Code)
		ctx.JSON(map[string]interface{}{
			"info":   cError.Error(),
			"status": "error",
		})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(map[string]interface{}{
		"data":   q,
		"status": "OK",
	})
	return
}

func queueDelete(ctx context.Context) {
	err := callback.CallbackQueueDelete(ctx.Params().Get("queueId"))
	if err != nil {
		ctx.StatusCode(err.Code)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	ctx.JSON(map[string]interface{}{
		"info":   "Success",
		"status": "OK",
	})
}

func queueUpdate(ctx context.Context) {
	var q map[string]interface{}
	err := ctx.ReadJSON(&q)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	e := callback.CallbackQueueUpdate(auth.GetSessionFromContext(ctx), ctx.Params().Get("queueId"), q)
	if e != nil {
		ctx.StatusCode(e.Code)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	ctx.JSON(map[string]interface{}{
		"info":   "Success",
		"status": "OK",
	})
}

func membersList(ctx context.Context) {
	data, err := callback.CallbackMembersList(ctx.Params().Get("queueId"), db.RequestListFromUri(ctx))
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

func memberCreate(ctx context.Context) {

	var m *callback.Member
	err := ctx.ReadJSON(&m)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	cError := callback.CallbackMemberCreate(ctx.Params().Get("queueId"), m)

	if cError != nil {
		ctx.StatusCode(cError.Code)
		ctx.JSON(map[string]interface{}{
			"info":   cError.Error(),
			"status": "error",
		})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(map[string]interface{}{
		"data":   m,
		"status": "OK",
	})
	return
}

func memberItem(ctx context.Context) {
	data, err := callback.CallbackMemberItem(ctx.Params().Get("queueId"), ctx.Params().Get("memberId"))
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

func memberUpdate(ctx context.Context) {
	var q map[string]interface{}
	err := ctx.ReadJSON(&q)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	e := callback.CallbackMemberUpdate(auth.GetSessionFromContext(ctx), ctx.Params().Get("queueId"), ctx.Params().Get("memberId"), q)
	if e != nil {
		ctx.StatusCode(e.Code)
		ctx.JSON(map[string]interface{}{
			"info":   e.Error(),
			"status": "error",
		})
		return
	}

	ctx.JSON(map[string]interface{}{
		"info":   "Success",
		"status": "OK",
	})
}

func memberDelete(ctx context.Context) {
	err := callback.CallbackMemberDelete(ctx.Params().Get("queueId"), ctx.Params().Get("memberId"))
	if err != nil {
		ctx.StatusCode(err.Code)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	ctx.JSON(map[string]interface{}{
		"info":   "Success",
		"status": "OK",
	})
}

func memberAddComment(ctx context.Context) {
	var comment *callback.Comment
	err := ctx.ReadJSON(&comment)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	cError := callback.CallbackMemberCommentAdd(
		auth.GetSessionFromContext(ctx),
		ctx.Params().Get("queueId"),
		ctx.Params().Get("memberId"),
		comment,
	)

	if cError != nil {
		ctx.StatusCode(cError.Code)
		ctx.JSON(map[string]interface{}{
			"info":   cError.Error(),
			"status": "error",
		})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(map[string]interface{}{
		"data":   comment,
		"status": "OK",
	})
	return
}

func listCall(ctx context.Context) {
	//db.RequestFromUri(ctx)
	//data, err := services.FindGetCall(nil)
	//if err != nil {
	//	ctx.StatusCode(500)
	//	ctx.JSON(map[string]interface{}{
	//		"info":   err.Error(),
	//		"status": "error",
	//	})
	//	return
	//}

	ctx.JSON(map[string]interface{}{
		//"data":   data,
		"status": "OK",
	})
}
