package response

type ErrorMessages string

const (
	EmSomethingWentWrong     ErrorMessages = "Something went wrong, Please try again later"
	SymbolNotSuported        ErrorMessages = "Symbol not supported"
	MinQuoteQuantityError    ErrorMessages = "Minimum quote quantity breached error"
	MaxQuoteQuantityError    ErrorMessages = "Maximum quote quantity breached error"
	MinBaseQuantityError     ErrorMessages = "Minimum base quantity breached error"
	MaxBaseQuantityError     ErrorMessages = "Maximum base quantity breached error"
	UnableToGetResponseError ErrorMessages = "Unable to get response or incorrect response"
	IncorrectExchange        ErrorMessages = "Please provide valid exchange name"
)
