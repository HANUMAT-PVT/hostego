package response

import (
	"backend-hostego/internal/app/hostego-service/constants/string_constants"
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/pkg/logger"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	log = logger.GetLogger()
)

type Response gin.H

func FormatResponse(data interface{}, status bool, err error, message string) *Response {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	res := &Response{
		"success": status,
		"data":    data,
		"code":    errMsg,
		"msg":     message,
	}
	return res
}

func Success(ctx *gin.Context, reqCtx dto.ReqCtx, statusCode int, data interface{}) {
	log := logger.GetLogger()
	if ctx.Writer.Written() {
		log.WarnWithStruct("response body was already written! will not overwrite", dto.LoggerWithFields{ReqCtx: reqCtx})
		return
	}
	res := FormatResponse(data, true, nil, string_constants.EMPTY_STRING)

	ctx.Writer.Header().Set("X-Request-Id", ctx.GetString("X-Request-Id"))
	ctx.Writer.Header().Set("X-Span-Request-Id", ctx.GetString("X-Span-Request-Id"))
	ctx.Writer.Header().Set("X-Amzn-Trace-Id", ctx.GetString("X-Amzn-Trace-Id"))

	log.Info(dto.LoggerWithFields{ReqCtx: reqCtx, Message: fmt.Sprintf("Response:-%v", res)})

	ctx.JSON(statusCode, res)
}

func Fail(ctx *gin.Context, statusCode int, errors []gin.H, msg string) {
	log := logger.GetLogger()
	if ctx.Writer.Written() {
		log.Warn("response body was already written! will not overwrite")
		return
	}
	res := FormatResponse(Response{
		"errors": errors,
	}, false, nil, msg)

	ctx.JSON(statusCode, res)
}

func Error(ctx *gin.Context, reqCtx dto.ReqCtx, statusCode int, errCode ErrorCode, msg string) {
	if ctx.Writer.Written() {
		log.WarnWithStruct("response body was already written! will not overwrite", dto.LoggerWithFields{ReqCtx: reqCtx})
		return
	}
	res := Response{
		"success": false,
		"data":    nil,
		"code":    errCode,
		"msg":     msg,
	}

	ctx.Writer.Header().Set("X-Request-Id", ctx.GetString("X-Request-Id"))
	ctx.Writer.Header().Set("X-Span-Request-Id", ctx.GetString("X-Span-Request-Id"))
	ctx.Writer.Header().Set("X-Amzn-Trace-Id", ctx.GetString("X-Amzn-Trace-Id"))

	log.Info(dto.LoggerWithFields{ReqCtx: reqCtx, Message: fmt.Sprintf("Response:-%v", res)})
	ctx.JSON(statusCode, res)
}
