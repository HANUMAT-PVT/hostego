package config_constants

const (
	CONFIG_PATH                      = "./config"
	CONFIG_PATH_LAMBDA               = "./../../config"
	LOCAL                            = "local"
	YML                              = "yml"
	HOST_TYPE                        = "HOST_TYPE"
	CONFIG_TYPE                      = "CONFIG_TYPE"
	CONNECTION_ADDRESS               = "0.0.0.0:8080"
	LAMBDA                           = "LAMBDA"
	LAMDA_CONFIG_PATH                = "PATH"
	PREPROD_HOST                     = "preprod"
	PROD_HOST                        = "prod"
	AUTHORIZATION                    = "Authorization"
	CONTENT_TYPE_XXX_FORM_URLENCODED = "application/x-www-form-urlencoded"
	REQUEST_FORMAT_JSON              = "JSON"
	REQUEST_FORMAT_XML               = "XML"
	REQUEST_FORMAT_PIPE_DELIMITED    = "PIPE_DELIMITED"
	RETRY_COUNT                      = "Retry-Count"
)

var ProdHosts = []interface{}{PROD_HOST, PREPROD_HOST}
var LocalHosts = []interface{}{LOCAL}
