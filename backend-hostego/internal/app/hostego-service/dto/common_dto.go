package dto

import (
	"backend-hostego/internal/pkg/logger"
	"context"

	goBricksLogger "backend-hostego/internal/app/hostego-service/hostego-logger/go-bricks/logger"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type ReqCtx struct {
	Url                 string
	Method              string
	ReqId               string
	SpanReqId           string
	UserId              int
	UserCustodyWalletId int64
	Source              string
	AmazonTraceId       string
	NewRelicTxn         *newrelic.Transaction
	GoBricksLog         goBricksLogger.Logger
	NrTraceId           string
	Log                 *logger.MyLogger
	context.Context
}

type QtyDto struct {
	Value string `json:"value"`
	Unit  string `json:"unit"`
}
