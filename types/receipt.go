package types

// Event is emitted while executing transactions
type Event struct {
	Name string
	Data string
}

// TransactionReceipt reflects corresponding Transaction execution result
type TransactionReceipt struct {
	TransactionHash string
	Result          uint64
	GasUsed         uint64
	Events          []*Event
}
