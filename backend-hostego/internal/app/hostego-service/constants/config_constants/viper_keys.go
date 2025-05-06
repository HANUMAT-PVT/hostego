package config_constants

const (
	VKEYS_HOST_IP   = "host.ip"
	VKEYS_HOST_PORT = "host.port"
	VKEYS_HOST_TYPE = "host.type"

	VKEYS_DATABASE_POSTGRES_SOURCE_HOST              = "database.postgres.source.host"
	VKEYS_DATABASE_POSTGRES_SOURCE_PORT              = "database.postgres.source.port"
	VKEYS_DATABASE_POSTGRES_SOURCE_DB_NAME           = "database.postgres.source.db_name"
	VKEYS_DATABASE_POSTGRES_SOURCE_PASSWORD          = "database.postgres.source.password"
	VKEYS_DATABASE_POSTGRES_SOURCE_USER              = "database.postgres.source.user"
	VKEYS_DATABASE_POSTGRES_SOURCE_MAX_IDLE_CONN     = "database.postgres.max_idle_connections"
	VKEYS_DATABASE_POSTGRES_SOURCE_MAX_OPEN_CONN     = "database.postgres.max_open_connections"
	VKEYS_DATABASE_POSTGRES_SOURCE_MAX_CONN_LIFETIME = "database.postgres.max_connection_lifetime"

	VKEYS_CORS_ORIGINS              = "cors.origins"
	VKEYS_SERVER_ALLOW_HEADERS      = "server.header_allows"
	VKEYS_ALLOW_METHODS             = "server.method_allows"
	VKEYS_EXPOSED_HEADERS           = "server.exposed_headers"
	VKEYS_IDLE_TIMEOUT_SERVER       = "server.idle_timeout"
	VKEYS_READ_WRITE_TIMEOUT_SERVER = "server.read_write_timeout"

	VKEYS_NEWRELIC_LICENSE  = "new_relic.license"
	VKEYS_NEWRELIC_ENABLED  = "new_relic.enabled"
	VKEYS_NEWRELIC_APP_NAME = "new_relic.app_name"

	VKEYS_REDIS_HOST = "redis_cluster.host_cluster_url"

	VKEYS_REDIS_CLUSTERS_HOST_URL             = "redis_cluster.host_cluster_url"
	VKEYS_REDIS_CLUSTERS_POOL_SIZE            = "redis_cluster.pool_size"
	VKEYS_INTERCEPTOR_REDIS_CLUSTERS_HOST_URL = "interceptor_redis_cluster.host_cluster_url"

	VKEYS_LOGGING_FORMAT = "logging.format"
	VKEYS_LOGGING_LEVEL  = "logging.level"

	VKEYS_ADMIN_TOKEN = "admin.token"

	VKEYS_AUTH_SERVICE_BASE_URL = "auth_service.base_url"
	VKEYS_AUTH_SERVICE_TOKEN    = "auth_service.token"

	VKEYS_CS_INDIA_BASE_URL        = "cs_india.base_url"
	VKEYS_CS_INDIA_AUTH_TOKEN      = "cs_india.auth_token"
	VKEYS_CS_INDIA_PORTFOLIO_TOKEN = "cs_india.portfolio_auth_token"

	VKEYS_RATE_SERVICE_BASE_URL = "rate_service.base_url"
	VKEYS_RATE_SERVICE_TOKEN    = "rate_service.token"

	ZAPPA_AUTH_TOKEN        = "zappa_token"
	VKEYS_DG_CONFIG_TIMEOUT = "db_config_timeout"

	PROXY_USERNAME = "proxy.username"
	PROXY_PASSWORD = "proxy.password"
	PROXY_URL      = "proxy.url"

	BINANCE_URL        = "binance.base_url"
	BINANCE_API_KEY    = "binance.api_key"
	BINANCE_API_SECRET = "binance.api_secret"

	VKEYS_INVOICE_BASE_URL   = "cs_india.base_url"
	VKEYS_INVOICE_AUTH_TOKEN = "cs_india.auth_token"

	VKEYS_DROGON_SER_AUTH_TOKEN = "drogon_service_auth_token.token"

	AWS_REGION = "aws.region"

	KMS_KEY_ID = "kms.key_id"

	VKEYS_ASSET_BASE_URL              = "asset_service.base_url"
	VKEYS_ASSET_AUTH_TOKEN            = "asset_service.auth_token"
	VKEYS_ASSET_SERVICE_REDIS_HOST    = "asset_service.redis_host"
	VKEYS_ASSET_SERVICE_STREAM_ENABLE = "asset_service.stream_enable"
	VKEYS_ASSET_REDIS_STREAM_NAME     = "asset_service.stream_name"

	VKEYS_ORDER_SERVICE_BASE_URL       = "order_service.base_url"
	VKEYS_ORDER_SERVICE_API_KEY_ID     = "order_service.api_key_id"
	VKEYS_ORDER_SERVICE_API_KEY_SECRET = "order_service.api_key_secret"
	VKEYS_ORDER_SERVICE_TENANT         = "order_service.tenant"
	VKEYS_ORDER_SERVICE_CLIENT_ID      = "order_service.client_id"
	VKEYS_ORDER_SERVICE_CLIENT_TOKEN   = "order_service.client_token"

	VKEYS_TDS_SERVICE_BASE_URL   = "tds_service.base_url"
	VKEYS_TDS_SERVICE_AUTH_TOKEN = "tds_service.auth_token"

	VKEYS_PRO_BACKEND_SERVICE_BASE_URL = "pro_backend.base_url"
	VKEYS_PRO_BACKEND_AUTH_TOKEN       = "pro_backend.auth_token"

	LEDGER_SECRET             = "ledger.secret"
	LEDGER_ENDPOINT           = "ledger.endpoint"
	LEDGER_READ_ONLY_ENDPOINT = "ledger.read_only_endpoint"

	CONFIG_READER_HOST            = "config_service.config_reader.host"
	CONFIG_READER_HOST_TYPE       = "config_service.config_reader.host_type"
	CONFIG_READER_TENANT          = "config_service.config_reader.tenant"
	CONFIG_READER_SERVICE         = "config_service.config_reader.service_name"
	CONFIG_READER_REQUEST_TIMEOUT = "config_service.config_reader.timeout"
	CONFIG_READER_LOCAL_FILE_PATH = "config_service.config_reader.local_file_path"

	CONFIG_WRITER_HOST            = "config_service.config_writer.host"
	CONFIG_WRITER_HOST_TYPE       = "config_service.config_writer.host_type"
	CONFIG_WRITER_TENANT          = "config_service.config_writer.tenant"
	CONFIG_WRITER_SERVICE         = "config_service.config_writer.service_name"
	CONFIG_WRITER_REQUEST_TIMEOUT = "config_service.config_writer.timeout"
	CONFIG_WRITER_LOCAL_FILE_PATH = "config_service.config_writer.local_file_path"

	CUSTOM_METRICS_SAMPLE_PERCENTAGE = "custom_metrics.sample_percentage"

	VKEYS_ALPHA_BASE_URL = "alpha_service.base_url"
	VKEYS_ALPHA_TOKEN    = "alpha_service.auth_token"

	VkeysFuturesOSBaseUrl   = "futures_os.base_url"
	VkeysFuturesOSAuthToken = "futures_os.auth_token"
	VkeysFuturesTimeout     = "futures_os.timeout"
)
