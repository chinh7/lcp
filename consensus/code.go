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
	CodeTypeAccountNotExist       uint32 = 7
	CodeTypeInsufficientFee       uint32 = 8
	CodeTypeInvalidData           uint32 = 9
	CodeTypeNonContractAccount    uint32 = 10
	CodeTypeInvalidGasPrice       uint32 = 11
	CodeTypeInvalidPubKey         uint32 = 12
)
