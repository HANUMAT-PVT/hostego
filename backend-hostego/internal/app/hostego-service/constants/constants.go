package constants

// logging
const (
	Zap             = "zap"
	CommonFieldsKey = "commonFields"
	LogLevelDebug   = "debug"
	LogLevelInfo    = "info"
	LogLevelWarn    = "warn"
)

const (
	RequestPath   = "request_path"
	RequestMethod = "request_method"
)

const (
	ContextLogger = "logger"
	ContextNrTxn  = "nr_txn"
)

const (
	UserIdCtxKey            = "ctx_user_id"
	UserCustodyWalletCtxKey = "ctx_user_custody_wallet_id"
	InrFuturesSubaccount    = "ctx_inr_futures_sub_account"
	CustodyWalletLedger     = "ctx_custody_wallet_ledger"
)

const (
	SourceContextKey = "source"
)

const (
	INVALID_USER_AUTH_TOKEN            = "invalid user auth token"
	MALFORMED_REQUEST_DATA             = "Malformed request data"
	UNMARSHALL_ERROR                   = "unable to unmarshall"
	MARSHALL_ERROR                     = "unable to marshall"
	CRITICAL_ISSUE                     = "critical issue"
	SOMETHING_WENT_WRONG_ERROR_MESSAGE = "Something went wrong, Try again later"
	DB_ERROR                           = "error interacting with DB"

	INVALID_CLIENT_ID_ERROR = "invalid client id"

	INVALID_BROKER           = "invalid broker"
	INVALID_TAG_REQUEST_DATA = "please pass tag_id or tag_slug"

	ORDER_RAISE_FAILED        = "Failed to raise order on Exchange"
	ORDER_STATUS_CHECK_FAILED = "Failed to check order status"
	GET_ALL_ORDERS_FAILED     = "Failed to get all orders"
	GET_ORDER_TRADES_FAILED   = "Failed to get order trades"

	SUB_ACCOUNT_DEBIT_MONEY_FAILED   = "Failed to debit money from user subaccount"
	SUB_ACCOUNT_CREDIT_MONEY_FAILED  = "Failed to credit money to user subaccount"
	LEDGER_DEBIT_MONEY_FAILED_ERROR  = "Failed to block margin from ledger"
	LEDGER_CREDIT_MONEY_FAILED_ERROR = "Failed to credit in user ledger"

	FUTURES_DISABLED              = "Currently futures is disabled at the moment"
	ASSET_LEVEL_VALIDATION_FAILED = "asset validation failed"
	INVALID_CLIENT_ORDER_ID       = "Invalid client order id format. Is should be UUID"
	INVALID_EXCHANE               = "Invalid exchange. Exchange not supported"
	INVALID_ORDER_SIDE            = "Invalid order side"
	INVALID_ORDER_TYPE            = "Invalid order type"
	INVALID_PLATFORM              = "Invalid platform."
	ORDER_ALREADY_PLACED          = "Order already placed"
	USER_KYC_NOT_DONE             = "User kyc not completed."
	ERROR_WHILE_FETCHING_ID       = "Error while getting user subscriptions"
	INVALID_REQUEST_PAYLOAD       = "Error while processing payload or Payload input is wrong"
)
