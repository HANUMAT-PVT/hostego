package middlewares

import (
	"backend-hostego/internal/app/hostego-service/constants"
	"backend-hostego/internal/app/hostego-service/constants/api_constants"
	"backend-hostego/internal/app/hostego-service/hostego-logger/go-bricks/logger"
	"backend-hostego/internal/app/hostego-service/utilities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
	// Import Zap logger
)

func beforeRequest(ctx *gin.Context) {
	url := ctx.Request.URL.Path
	method := ctx.Request.Method

	TraceReqID := ctx.GetHeader(string(constants.HeaderTraceReqId))
	TraceAmznID := ctx.GetHeader(string(constants.HeaderTraceAmznID))
	TraceSpanID := uuid.New().String()

	debugMode := ctx.GetHeader(string(constants.HeaderDebugMode))
	logLevel := ""
	if debugMode == "true" {
		logLevel = constants.LogLevelDebug
	}
	log := logger.GetNewLogger(ctx, constants.Zap, logLevel)

	if TraceReqID == "" {
		TraceReqID = uuid.New().String()
		log.Debugf("No RequestId passed setting as : %s", TraceReqID)
	}

	// TODO: Remove this once the Caller access token key attribute is fixed to X-Secure-Key
	if ctx.GetHeader(string(constants.HeaderSecureToken)) != "" {
		log.InfofCf("request api url %s, is passing the header X-Secure-Token", url)
	}

	userId := ctx.GetHeader(string(constants.HeaderUserId))
	UserCustodyWalletId := ctx.GetHeader(string(constants.HeaderCustodyWalletId))

	txn := newrelic.FromContext(ctx)
	var NrTraceId string
	if txn != nil {
		traceMetadata := txn.GetTraceMetadata()
		NrTraceId = traceMetadata.TraceID
	}
	commonFields := map[string]interface{}{
		constants.RequestPath:                   url,
		constants.RequestMethod:                 method,
		string(constants.HeaderTraceReqId):      TraceReqID,
		string(constants.HeaderTraceAmznID):     TraceAmznID,
		string(constants.HeaderTraceSpanID):     TraceSpanID,
		string(constants.HeaderNrTraceId):       NrTraceId,
		string(constants.HeaderUserId):          userId,
		string(constants.HeaderCustodyWalletId): UserCustodyWalletId,
	}
	ctx.Set(constants.CommonFieldsKey, commonFields)
	log.AddFields(commonFields)

	ctx.Set(constants.ContextLogger, log)
	ctx.Set(constants.ContextNrTxn, txn)

	//Not logging Request content for health check api
	if !utilities.IsContains([]string{api_constants.HEALTH, api_constants.METRICS}, ctx.Request.URL.Path) {
		log.WarnfCf("API %+v is hit", url)
	}
}

func afterRequest(ctx *gin.Context) {

}

func RequestRequirements() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		beforeRequest(ctx)
		ctx.Next()
		afterRequest(ctx)
	}
}
