package error_constants

const (
	InvalidInputStatusCode    int = 422
	InternalServerStatusCode  int = 500
	NotFoundStatusCode        int = 404
	BadRequestStatusCode      int = 400
	FailedToBlockLedger       int = 424
	MinMaxQuoteBaseBreached   int = 423
	UnableToBlockLedgerFromOS int = 402
)
