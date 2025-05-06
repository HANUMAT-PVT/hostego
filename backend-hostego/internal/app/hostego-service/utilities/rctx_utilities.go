package utilities

import (
	"backend-hostego/internal/app/hostego-service/constants"
	"backend-hostego/internal/app/hostego-service/dto"
	"backend-hostego/internal/pkg/logger"
	"context"
	"fmt"
	"strconv"

	goBricksLogger "backend-hostego/internal/app/hostego-service/hostego-logger/go-bricks/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/sirupsen/logrus"
)

func AddGoBricksLogger(rctx *dto.ReqCtx) {
	ctx := rctx.Context.(*gin.Context)
	logger, ok := ctx.Value(constants.ContextLogger).(goBricksLogger.Logger)
	if !ok {
		logger = goBricksLogger.GetNewLogger(rctx.Context, constants.Zap)
	}
	commonFields, _ := ctx.Value(constants.CommonFieldsKey).(map[string]interface{})

	// User Id
	userId := ctx.Request.Header.Get(string(constants.HeaderUserId))
	ctxUserId := ctx.GetInt(constants.UserIdCtxKey)
	if ctxUserId != 0 {
		userId = fmt.Sprintf("%v", ctxUserId) // override userId from context
	}
	if userId != "" {
		userIdInt, err := strconv.Atoi(userId)
		if err == nil {
			rctx.UserId = userIdInt
			commonFields[string(constants.HeaderUserId)] = userId
		}
	}

	// Custody Wallet Id
	custodyWalletId := ctx.Request.Header.Get(string(constants.HeaderCustodyWalletId))
	ctxCustodyWalletId := ctx.GetInt(constants.UserCustodyWalletCtxKey)
	if ctxCustodyWalletId != 0 {
		custodyWalletId = fmt.Sprintf("%v", ctxCustodyWalletId) // override custodyWalletId from context
	}
	if custodyWalletId != "" {
		custodyWalletIdInt, err := strconv.Atoi(custodyWalletId)
		if err == nil {
			rctx.UserCustodyWalletId = int64(custodyWalletIdInt)
			commonFields[string(constants.HeaderCustodyWalletId)] = custodyWalletIdInt
		}
	}

	// Source
	source := ctx.GetString(constants.SourceContextKey)
	if source != "" {
		rctx.Source = source
		commonFields[string(constants.HeaderPlatformSource)] = source
	}

	// subaccount
	logger.AddFields(commonFields)
	rctx.GoBricksLog = logger
}

func GetRCtx(ctx *gin.Context, txnName string) dto.ReqCtx {
	myLogger := logger.GetLogger()
	myLogger.WithFields(logrus.Fields{
		"request_id":     ctx.Request.Header.Get("X-Request-Id"),
		"request_path":   ctx.Request.URL.Path,
		"request_method": ctx.Request.Method,
		"amzn_trace_id":  ctx.Request.Header.Get("X-Amzn-Trace-Id"),
	})
	rCtx := dto.ReqCtx{
		ReqId:         ctx.Request.Header.Get("X-Request-Id"),
		AmazonTraceId: ctx.Request.Header.Get("X-Amzn-Trace-Id"),
		SpanReqId:     ctx.GetString("X-Span-Request-Id"),
		NewRelicTxn:   nrgin.Transaction(ctx),
		Context:       ctx,
		Log:           myLogger,
	}

	AddGoBricksLogger(&rCtx)

	if rCtx.NewRelicTxn != nil {
		rCtx.NewRelicTxn.SetName(txnName)
		rCtx.NrTraceId = rCtx.NewRelicTxn.GetTraceMetadata().TraceID
	}
	return rCtx
}

func GetRCtxNonWeb(txnName string) dto.ReqCtx {
	myLogger := logger.GetLogger()

	rCtx := dto.ReqCtx{
		ReqId:         uuid.New().String(),
		AmazonTraceId: uuid.New().String(),
		SpanReqId:     uuid.New().String(),
		// NewRelicTxn:   GetNonWebNewRelicTxn(txnName),
		Log:     myLogger,
		Context: context.Background(),
	}

	if rCtx.NewRelicTxn != nil {
		rCtx.NewRelicTxn.SetName(txnName)
		rCtx.NrTraceId = rCtx.NewRelicTxn.GetTraceMetadata().TraceID
	}
	return rCtx
}
