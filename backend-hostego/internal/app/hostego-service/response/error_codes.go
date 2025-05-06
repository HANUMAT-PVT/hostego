package response

type ErrorCode string

const (
	CodeUnknown              ErrorCode = "CODE_UNKNOWN"
	CodePanic                ErrorCode = "CODE_PANIC"
	BadRequest               ErrorCode = "BAD_REQUEST"
	ErrorWhileProcessing     ErrorCode = "ERROR_WHILE_PROCESSING"
	Unauthorized             ErrorCode = "UNAUTHORIZED"
	InvalidInput             ErrorCode = "invalid_input"
	InvalidSource            ErrorCode = "invalid_source"
	EmptyUserCustodyWalletId ErrorCode = "empty_user_custody_wallet_id"
	EmptySource              ErrorCode = "empty_source"
)
