package response

import (
	"errors"
	"harmoni/internal/pkg/errorx"
	"harmoni/internal/pkg/httpx"
	"harmoni/internal/pkg/reason"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

func HandleResponse(ctx *gin.Context, err error, data interface{}) {
	// no error
	if err == nil {
		ctx.JSON(http.StatusOK, httpx.NewRespBodyData(http.StatusOK, reason.Success, data))
		return
	}

	var myErr *errorx.Error
	// unknown error
	if !errors.As(err, &myErr) {
		log.Error(err, "\n", errorx.LogStack(2, 5))
		ctx.JSON(http.StatusInternalServerError, httpx.NewRespBody(
			http.StatusInternalServerError, reason.UnknownError))
		return
	}

	// log internal server error
	if errorx.IsInternalServer(myErr) {
		log.Error(myErr)
	}

	respBody := httpx.NewRespBodyFromError(myErr)
	if data != nil {
		respBody.Data = data
	}
	ctx.JSON(int(myErr.Code), respBody)
}
