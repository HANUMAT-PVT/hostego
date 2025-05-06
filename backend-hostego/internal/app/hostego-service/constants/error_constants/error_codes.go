package error_constants

const (
	InvalidInput           string = "invalid_input"
	UNABLE_TO_GET_RESPONSE string = "Unable to get response or incorrect response"
	SomethingWentWrongCode string = "Something Went Wrong, Please Try Again"
	InternalServerError    string = "Internal Server Error"
	OrderNotFound          string = "Not Found"
	LowBalanceError        string = "Balance is low to block ledger"
	WalletDisabledCode     string = "Wallet Disabled"
	INVALID_ORDER_ID       string = "Invalid Order Id Format. It should be UUID"
)
