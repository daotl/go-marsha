// Package cbor_refmt provides a `Marshaler` implementation for CBOR backed by go-ipld-cbor and
// refmt packages.
package cbor_refmt

import (
	"bytes"
	"errors"
	"strings"

	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/ipfs/go-ipld-cbor/encoding"
	"github.com/polydawn/refmt/obj/atlas"

	"github.com/daotl/go-marsha"
)

var (
	ErrNotRegistered = errors.New("model type not registered")
	ErrTypeNotMatch  = errors.New("model type does not match")
)

var (
	emptyAtlas  = atlas.MustBuild()
	marshaler   = encoding.NewPooledMarshaller(emptyAtlas)
	unmarshaler = encoding.NewPooledUnmarshaller(emptyAtlas)
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
	return marshaler.Marshal(p)
}

func (m *CBORRefmtMarshaler) UnmarshalPrimitive(bin []byte, p interface{}) error {
	r := bytes.NewReader(bin)
	if err := unmarshaler.Decode(r, p); err != nil {
		estr := err.Error()
		if strings.Contains(estr, "wrong type") {
			return ErrTypeNotMatch
		} else if strings.Contains(estr, "EOF") {
			return err
		} else {
			panic("unexpected error in unmarshalling")
		}
	}
	return nil
}

func (m *CBORRefmtMarshaler) MarshalStruct(p marsha.StructPtr) ([]byte, error) {
	bin, err := cbor.DumpObject(p)
	if err != nil {
		return nil, ErrNotRegistered
	}
	return bin, nil
}

func (m *CBORRefmtMarshaler) UnmarshalStruct(bin []byte, p marsha.StructPtr) error {
	if err := cbor.DecodeReader(bytes.NewReader(bin), p); err != nil {
		return ErrTypeNotMatch
	}
	return nil
}

func (m *CBORRefmtMarshaler) MarshalStructSlice(p marsha.StructSlicePtr) ([]byte, error) {
	return cbor.DumpObject(p)
}

func (m *CBORRefmtMarshaler) UnmarshalStructSlice(bin []byte, p marsha.StructSlicePtr) error {
	if err := cbor.DecodeReader(bytes.NewReader(bin), p); err != nil {
		return ErrTypeNotMatch
	}
	return nil
}
