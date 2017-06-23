package rest

import (
	"../../db"
	"../../services"
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
		callback.Put("/{queueId}", queueUpdate)    //TODO
		callback.Delete("/{queueId}", queueDelete) //+

		callback.Get("/{queueId}/members", membersList) //+
		//callback.Post("/{queueId}/members", memberCreate)
		callback.Get("/{queueId}/members/{memberId}", memberItem)      //+
		callback.Put("/{queueId}/members/{memberId}", memberUpdate)    //TODO
		callback.Delete("/{queueId}/members/{memberId}", memberDelete) //+
	}
}

func queueList(ctx context.Context) {
	data, err := services.CallbackQueueList(db.RequestListFromUri(ctx))
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
	data, err := services.CallbackQueueItem(ctx.Params().Get("queueId"))
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
	var q *services.Queue
	err := ctx.ReadJSON(&q)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	cError := services.CallbackQueueCreate(q)

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
	err := services.CallbackQueueDelete(ctx.Params().Get("queueId"))
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
	ctx.Text("queueUpdate " + ctx.Params().Get("queueId"))
}

func membersList(ctx context.Context) {
	data, err := services.CallbackMembersList(ctx.Params().Get("queueId"), db.RequestListFromUri(ctx))
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

	var m *services.Member
	err := ctx.ReadJSON(&m)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"info":   err.Error(),
			"status": "error",
		})
		return
	}

	cError := services.CallbackMemberCreate(ctx.Params().Get("queueId"), m)

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
	data, err := services.CallbackMemberItem(ctx.Params().Get("queueId"), ctx.Params().Get("memberId"))
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
	ctx.Text("memberUpdate " + ctx.Params().Get("queueId") + " " + ctx.Params().Get("memberId"))
}

func memberDelete(ctx context.Context) {
	err := services.CallbackMemberDelete(ctx.Params().Get("queueId"), ctx.Params().Get("memberId"))
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
