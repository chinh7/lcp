package crypto

import (
	"errors"
	"reflect"
	"testing"
)

func TestAddressFromString(t *testing.T) {
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
			want: AddressFromBytes([]byte{
				0x00, 0x58, 0x07, 0x2a, 0x26, 0x0b, 0x42, 0xa7, 0xcb,
				0x04, 0x2b, 0x32, 0xd3, 0xe8, 0x6f, 0xc3, 0x20, 0x53,
				0xe5, 0x14, 0x30, 0x42, 0x00, 0x11, 0xf8, 0x3b, 0xcd,
				0x8b, 0xf6, 0xa0, 0x9c, 0x8a, 0x33, 0x48, 0xa3, 0xbb}),
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddressFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
