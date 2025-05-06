package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response gin.H

type Status string

const (
	SUCCESS      Status = "SUCCESS"
	FAIL         Status = "FAIL"
	UNAUTHORIZED Status = "UNAUTHORIZED"
)

type Message string

const (
	SomethingWentWrong Message = "Something went wrong, Please try again later"
)

func FormatResponse(data interface{}, status Status, err error, message Message) *Response {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	res := &Response{
		"status":  status,
		"payload": data,
		"err":     errMsg,
		"msg":     message,
	}
	return res
}

// data argument should be gin.H{} or []gin.H{}
func Success(ctx *gin.Context, statusCode int, data interface{}, msg string) {
	if ctx.Writer.Written() {
		fmt.Println("response body was already written! will not overwrite")
		return
	}
	res := FormatResponse(data, SUCCESS, nil, Message(msg))
	ctx.JSON(statusCode, res)
}

func Fail(ctx *gin.Context, statusCode int, err error, msg Message) {
	status := FAIL
	if statusCode == http.StatusUnauthorized {
		status = UNAUTHORIZED
	}

	if ctx.Writer.Written() {
		fmt.Println("response body was already written! will not overwrite")
		return
	}
	res := FormatResponse(nil, status, err, msg)

	ctx.JSON(statusCode, res)

}
