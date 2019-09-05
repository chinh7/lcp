package api

// KVPair is struct for key value pair
type KVPair struct {
	Key   []byte `json:"key,omitempty"`
	Value []byte `json:"value,omitempty"`
}
