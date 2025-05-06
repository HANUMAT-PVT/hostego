package constants

// Response constants
const (
	SUCCESS = "SUCCESS"
	FAIL    = "FAIL"
)
const (
	AUTOMATION         string = "AUTOMATION"
	SYSTEM             string = "SYSTEM"
	AUTOMATION_ENABLED string = "true"
	ONE_BOX            string = "ONE_BOX"
)
const (
	REQUEST_FORMAT_XML            = "XML"
	REQUEST_FORMAT_JSON           = "JSON"
	REQUEST_FORMAT_PIPE_DELIMITED = "PIPE_DELIMITED"
	REQUEST_METHOD_GET            = "GET"
	REQUEST_METHOD_POST           = "POST"
	REQUEST_METHOD_PUT            = "PUT"
	REQUEST_METHOD_OPTIONS        = "OPTIONS"
	REQUEST_HEADER_X_SECURE_TOKEN = "X-Secure-Token"
	SERVICE_NAME                  = "futures_service"
)

var RETRY_STATUS_CODE = 502

const CS_INDIA_USER_DETAIL_PATH = "api/v1/payment-service/get-user-details"
const Algo_Trading_Financial_Info = "Please use Financial Info API for this information"

var DISABLE_EXTERNAL_URL_LOGS = []string{CS_INDIA_USER_DETAIL_PATH}

const USDT = "usdt"
const INR = "inr"
const USDT_INR = "USDT_INR"
const ORDER_TIME_INFORCE = "GTC"
const REDIS_PREFIX = "{drogon}"

const GSTFactor float64 = 0.18
const PercentageValue int64 = 100
const LimitValue int64 = 10
const MaxLimitValue int64 = 50

const FactorType string = "factor"
const ASSET_SERVICE_CLIENT_ID_SPOT = "1"
const ASSET_SERVICE_CLIENT_ID_FUTURES = "2"

const (
	WAZIRX_EXCHANGE  = "WAZIRX"
	BINANCE_EXCHANGE = "BINANCE"
	BYBIT_EXCHANGE   = "BYBIT"
)

const (
	ExchangeIdContextKey = "exchange_id"
	ORDER_SOURCE_BOOST   = "BOOST"
)
const BLACKLISTED_COINS_CACHE_KEY = "blacklisted_coins_data_"

const ORDER_SOURCE_CSPRO = "CS_PRO"
