// Package cbor_refmt provides a `Marshaler` implementation for CBOR backed by go-ipld-cbor and
// refmt packages.
package cbor_refmt

import (
	"bytes"
	"errors"

	"github.com/daotl/go-marsha"
	cbor "github.com/ipfs/go-ipld-cbor"
)

var (
	ErrNotRegistered = errors.New("model type not registered")
)

// CBORRefmtMarshaler is a Marshaler implementation for CBOR backed by `go-ipld-cbor` and `refmt` packages.
//
// A Struct must first be registered by calling CBORRefmtMarshaler.Register(Struct{}) before being
// able to be marshaled/unmarshaled.
//
// Marshaling/unmarshaling can be customized by `refmt` tags:
//
//   type Model struct {
//   	Foo string `refmt:"bar,omitempty"`
//   }
//
type CBORRefmtMarshaler struct{}

var _ marsha.Marshaler = (*CBORRefmtMarshaler)(nil)

// New creates a CBORRefmtMarshaler.
func New() *CBORRefmtMarshaler {
	return &CBORRefmtMarshaler{}
}

// Register a Struct type by passing empty a Struct.
func (m *CBORRefmtMarshaler) Register(i interface{}) {
	cbor.RegisterCborType(i)
}

func (m *CBORRefmtMarshaler) MarshalPrimitive(p interface{}) ([]byte, error) {
	return cbor.DumpObject(p)
}

func (m *CBORRefmtMarshaler) UnmarshalPrimitive(bin []byte, p interface{}) error {
	return cbor.DecodeReader(bytes.NewReader(bin), p)
}

func (m *CBORRefmtMarshaler) MarshalStruct(p marsha.StructPtr) ([]byte, error) {
	bin, err := cbor.DumpObject(p)
	if err != nil {
		return nil, ErrNotRegistered
	}
	return bin, nil
}

func (m *CBORRefmtMarshaler) UnmarshalStruct(bin []byte, p marsha.StructPtr) error {
	return cbor.DecodeReader(bytes.NewReader(bin), p)
}

func (m *CBORRefmtMarshaler) MarshalStructSlice(p marsha.StructSlicePtr) ([]byte, error) {
	return cbor.DumpObject(p)
}

func (m *CBORRefmtMarshaler) UnmarshalStructSlice(bin []byte, p marsha.StructSlicePtr) error {
	return cbor.DecodeReader(bytes.NewReader(bin), p)
}
