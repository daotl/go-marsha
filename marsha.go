// Package marsha provides a standard data marshaling and unmarshaling interface which can
// be implemented by different encodings and implementations such as CBOR and Protocol Buffers.
package marsha

import (
	"errors"
)

var (
	ErrUnimplemented = errors.New("unimplemented")
)

// Struct should be implemented by structs you want to marshal/unmarshal.
type Struct interface {
	// Ptr returns the pointer to Struct.
	Ptr() StructPtr
}

// StructPtr should be implemented by pointers to structs you want to marshal/unmarshal.
type StructPtr interface {
	// Val returns the struct that StructPtr points to.
	Val() Struct
}

// StructSlicePtr should be implemented by pointers to struct slices you want to marshal/unmarshal.
type StructSlicePtr interface {
	// Val returns the struct slice that StructSlicePtr points to.
	Val() []Struct
}

// Marshaler is a standard data marshaling and unmarshaling interface which can
// be implemented by different encodings and implementations such as CBOR and Protocol Buffers.
type Marshaler interface {
	// MarshalPrimitive marshals the primitive value/slice `ptr` points to into bytes.
	MarshalPrimitive(ptr interface{}) ([]byte, error)

	// UnmarshalPrimitive unmarshals bytes `bs` into the primitive value/slice `ptr` points to.
	UnmarshalPrimitive(bs []byte, ptr interface{}) error

	// MarshalStruct marshals the struct `ptr` points to into bytes.
	MarshalStruct(ptr StructPtr) ([]byte, error)

	// UnmarshalStruct unmarshals bytes `bs` into the struct `ptr` points to.
	UnmarshalStruct(bs []byte, ptr StructPtr) error

	// MarshalStructSlice marshals the struct slice `ptr` points to into bytes.
	MarshalStructSlice(ptr StructSlicePtr) ([]byte, error)

	// UnmarshalStructSlice unmarshals bytes `bs` into the struct slice `ptr` points to.
	UnmarshalStructSlice([]byte, StructSlicePtr) error
}
