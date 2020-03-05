package consensus

// Return codes for the examples
const (
	CodeTypeOK                    uint32 = 0
	CodeTypeEncodingError         uint32 = 1
	CodeTypeBadNonce              uint32 = 2
	CodeTypeUnauthorized          uint32 = 3
	CodeTypeUnknownError          uint32 = 4
	CodeTypeExceedTransactionSize uint32 = 5
	CodeTypeInvalidSignature      uint32 = 6
	CodeTypeContractNotFound      uint32 = 7
	CodeTypeInsufficientFee       uint32 = 8
	CodeTypeInvalidData           uint32 = 9
)
