package constants

type Header string

const (
	HeaderAccessToken     Header = "X-Access-Token"
	HeaderSecureToken     Header = "X-Secure-Token"
	HeaderForwardedFor    Header = "X-Forwarded-For"
	HeaderTraceReqId      Header = "X-Request-Id"
	HeaderTraceSpanID     Header = "X-Span-Request-Id"
	HeaderTraceAmznID     Header = "X-Amzn-Trace-Id"
	HeaderNrTraceId       Header = "nr-Trace-ID"
	HeaderDebugMode       Header = "X-Debug-Mode"
	HeaderUserId          Header = "X-USER-ID"
	HeaderCustodyWalletId Header = "X-USER-CUSTODY-WALLET-ID"
	MessageId             Header = "stream-message-id"
	StreamName            Header = "stream-name"
	HeaderSource          Header = "X-ORDER-SOURCE"
	HeaderPlatformSource  Header = "X-Platform-Id"
	HeaderExchangeId      Header = "X-EXCHANGE-ID"
)

const (
	HeaderContentType     = "Content-Type"
	HeaderContentTypeJson = "application/json"
	HeaderXAuthToken      = "X-Auth-Token"
	HeaderXRequestId      = "X-Request-Id"
)

const (
	UserCustodyWalletIdCtxKey = "ctx_user_custody_wallet_id"
	SourceCtxKey              = "ctx_source"
)
