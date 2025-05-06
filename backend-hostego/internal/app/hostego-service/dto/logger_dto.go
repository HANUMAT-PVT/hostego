package dto

import (
	"strconv"
)

type LoggerWithFields struct {
	UserId              int
	UserCustodyWalletId int
	Message             string
	ReqCtx
}

func (loggerWithFields LoggerWithFields) String() string {
	baseString := "pro-backend Log: "
	if loggerWithFields.ReqCtx.ReqId != "" {
		baseString += " X-Request-Id: " + loggerWithFields.ReqCtx.ReqId
	}
	if loggerWithFields.ReqCtx.SpanReqId != "" {
		baseString += " X-Span-Request-Id: " + loggerWithFields.ReqCtx.SpanReqId
	}
	if loggerWithFields.ReqCtx.AmazonTraceId != "" {
		baseString += " X-Amzn-Trace-Id: " + loggerWithFields.ReqCtx.AmazonTraceId
	}
	if loggerWithFields.UserId != 0 {
		baseString += " user_id: " + strconv.Itoa(loggerWithFields.UserId)
	}
	if loggerWithFields.Url != "" {
		baseString += " url: " + loggerWithFields.Url
	}
	return baseString + " Message: " + loggerWithFields.Message
}
