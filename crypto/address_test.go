package crypto

import (
	"errors"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestAddressFromString(t *testing.T) {
	validAddress, err := AddressFromBytes([]byte{
		0x58, 0x07, 0x2a, 0x26, 0x0b, 0x42, 0xa7, 0xcb, 0x04,
		0x2b, 0x32, 0xd3, 0xe8, 0x6f, 0xc3, 0x20, 0x53, 0xe5,
		0x14, 0x30, 0x42, 0x00, 0x11, 0xf8, 0x3b, 0xcd, 0x8b,
		0xf6, 0xa0, 0x9c, 0x8a, 0x33, 0x48, 0xa3, 0xbb})
	if err != nil {
		panic(err)
	}
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		want Address
		err  error
	}{
		{
			name: "valid address",
			args: args{address: "LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53"},
			want: validAddress,
		},
		{
			name: "invalid checksum",
			args: args{address: "LADXUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53"},
			want: Address{},
			err:  errors.New("invalid checksum"),
		},
		{
			name: "invalid base32",
			args: args{address: "LADabc"},
			want: Address{},
			err:  errors.New("base32 decode failed: illegal base32 data at input byte 3"),
		},
		{
			name: "invalid version",
			args: args{address: "BADXUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53"},
			want: Address{},
			err:  errors.New("Unexpected version 8"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddressFromString(tt.args.address)
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("AddressFromString() err = %v, want %v", err, tt.err)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddressFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddressFromBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
		err  error
	}{
		{
			name: "valid address",
			args: args{b: []byte{
				0x58, 0x07, 0x2a, 0x26, 0x0b, 0x42, 0xa7, 0xcb, 0x04,
				0x2b, 0x32, 0xd3, 0xe8, 0x6f, 0xc3, 0x20, 0x53, 0xe5,
				0x14, 0x30, 0x42, 0x00, 0x11, 0xf8, 0x3b, 0xcd, 0x8b,
				0xf6, 0xa0, 0x9c, 0x8a, 0x33, 0x48, 0xa3, 0xbb}},
			want: "LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53",
		},
		{
			name: "invalid checksum",
			args: args{b: []byte{
				0x58, 0x07, 0x2a, 0x26, 0x0b, 0x42, 0xa7, 0xcb, 0x04,
				0x2b, 0x32, 0xd3, 0xe8, 0x6f, 0xc3, 0x20, 0x53, 0xe5,
				0x14, 0x30, 0x42, 0x00, 0x11, 0xf8, 0x3b, 0xcd, 0x8b,
				0xf6, 0xa0, 0x9c, 0x8a, 0x33, 0x48, 0xa3, 0xba}},
			err: errors.New("invalid checksum"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddressFromBytes(tt.args.b)
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("AddressFromBytes() err = %v, want %v", err, tt.err)
			}
			if err == nil && tt.want != got.String() {
				t.Errorf("AddressFromBytes() = %s, want %s", got.String(), tt.want)
			}
		})
	}
}

func TestNewDeploymentAddress(t *testing.T) {
	sender, _ := AddressFromString("LADSUJQLIKT4WBBLGLJ6Q36DEBJ6KFBQIIABD6B3ZWF7NIE4RIZURI53")
	contract, _ := AddressFromString("LB5EPP7RST6IROFHLNKTLGKAFQTXGNY45CEAXPTGVT3K53ZXFMMAW575")
	contract2, _ := AddressFromString("LADAUIL4G5BB6DXOZPG4ES6UHVK4DJND4GADTMW7TDRI4P2B4O7NLJYF")
	type args struct {
		senderAddress Address
		senderNonce   uint64
	}
	tests := []struct {
		name string
		args args
		want Address
	}{{
		args: args{
			senderAddress: sender,
			senderNonce:   0,
		},
		want: contract,
	}, {
		args: args{
			senderAddress: sender,
			senderNonce:   1234,
		},
		want: contract2,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeploymentAddress(tt.args.senderAddress, tt.args.senderNonce); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeploymentAddress() = %v, want %v", got.String(), tt.want.String())
			}
		})
	}
}
