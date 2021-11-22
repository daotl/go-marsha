// Package protobuf provides a Marshaler implementation for Protocol Buffers backed by `*.pb.go`
// files pre-generated by `protoc`.
package protobuf

import (
	"errors"
	"runtime"

	"google.golang.org/protobuf/proto"

	"github.com/daotl/go-marsha"
)

var (
	ErrNotPBStructPtr = errors.New("not a protobuf.StructPtr")
	ErrWrongPBType    = errors.New("wrong protocol buffers type")
)

// Struct implementations should embed their corresponding `proto.Message`s by pointer.
type Struct interface {
	proto.Message
}

type StructPtr interface {
	marsha.Struct

	// EmptyPB should return an empty corresponding `proto.Message`.
	EmptyPB() proto.Message

	// LoadPB should load data from `m` into the Struct this StructPtr points to.
	LoadPB(m proto.Message) error

	// PB should return a corresponding `proto.Message` filled with data of the Struct this StructPtr points to.
	PB() proto.Message
}

// PBMarshaler is a `marsha.Marshaler` implementation for Protocol Buffers backed by `*.pb.go` files
// pre-generated by `protoc`.
// Only `MarshalStruct` and `UnmarshalStruct` are supported by this implementation because of the
// limitation of Protocol Buffers.
type PBMarshaler struct{}

var _ marsha.Marshaler = (*PBMarshaler)(nil)

// New creates a PBMarshaler.
func New() *PBMarshaler {
	return &PBMarshaler{}
}

// Not implemented
func (m *PBMarshaler) MarshalPrimitive(_ interface{}) ([]byte, error) {
	return nil, marsha.ErrUnimplemented
}

// Not implemented
func (m *PBMarshaler) UnmarshalPrimitive(_ []byte, _ interface{}) error {
	return marsha.ErrUnimplemented
}

func (m *PBMarshaler) MarshalStruct(p marsha.StructPtr) ([]byte, error) {
	pbp, ok := p.(StructPtr)
	if !ok {
		return nil, ErrNotPBStructPtr
	}
	return proto.Marshal(pbp.PB())
}

func (m *PBMarshaler) UnmarshalStruct(bin []byte, p marsha.StructPtr) (err error) {
	// Recover if type assertion in LoadPB fails
	defer func() {
		if r := recover(); r != nil {
			switch e := r.(type) {
			case *runtime.TypeAssertionError:
				err = ErrWrongPBType
			case error:
				err = e
			case string:
				err = errors.New(e)
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	pbp, ok := p.(StructPtr)
	if !ok {
		return ErrNotPBStructPtr
	}
	pb := pbp.EmptyPB()
	err = proto.Unmarshal(bin, pb)
	if err == nil {
		err = pbp.LoadPB(pb)
	}
	return
}

// Not implemented
func (m *PBMarshaler) MarshalStructSlice(_ marsha.StructSlicePtr) ([]byte, error) {
	return nil, marsha.ErrUnimplemented
}

// Not implemented
func (m *PBMarshaler) UnmarshalStructSlice(_ []byte, _ marsha.StructSlicePtr) error {
	return marsha.ErrUnimplemented
}
