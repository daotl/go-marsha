// Package marsha provides a standard data marshaling and unmarshaling interface which can
// be implemented by different encodings and implementations such as CBOR and Protocol Buffers.
package marsha

import (
	"errors"
	"io"
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
	Val() []StructPtr
}

// Marsha is a standard data marshaling and unmarshaling interface which can
// be implemented by different encodings and implementations such as CBOR and Protocol Buffers.
type Marsha interface {
	// MarshalPrimitive marshals the primitive value/slice `p` points to into bytes.
	MarshalPrimitive(p interface{}) ([]byte, error)

	// UnmarshalPrimitive unmarshals bytes `bin` into the primitive value/slice `p` points to.
	UnmarshalPrimitive(bin []byte, p interface{}) error

	// MarshalStruct marshals the struct `p` points to into bytes.
	MarshalStruct(p StructPtr) ([]byte, error)

	// UnmarshalStruct unmarshals bytes `bin` into the struct `p` points to.
	UnmarshalStruct(bin []byte, p StructPtr) error

	// MarshalStructSlice marshals the struct slice `p` points to into bytes.
	MarshalStructSlice(p StructSlicePtr) ([]byte, error)

	// UnmarshalStructSlice unmarshals bytes `bin` into the struct slice `p` points to.
	UnmarshalStructSlice(bin []byte, p StructSlicePtr) error

	// NewEncoder returns a new encoder that will transmit on the io.Writer.
	NewEncoder(w io.Writer) Encoder

	// NewDecoder returns a new decoder that reads from the io.Reader.
	NewDecoder(r io.Reader) Decoder
}

type Encoder interface {
	// EncodePrimitive marshals and transmits the primitive value/slice `p` points to,
	// guaranteeing that all necessary type information has been transmitted first.
	EncodePrimitive(p interface{}) error

	// EncodeStruct marshals and transmits the struct `p` points to,
	// Guaranteeing that all necessary type information has been transmitted first.
	EncodeStruct(p StructPtr) error

	// EncodeStructSlice marshals and transmits the struct slice `p` points to,
	// Guaranteeing that all necessary type information has been transmitted first.
	EncodeStructSlice(p StructSlicePtr) error
}

// A Decoder manages the receipt of type and data information read from the remote side of a connection.
// It is safe for concurrent use by multiple goroutines.
//
// The Decoder doesn't do sanity checking on decoded input sizes.
// Take caution when decoding data from untrusted sources.
type Decoder interface {
	// DecodePrimitive reads the next value from the input stream and stores it in the value/slice
	// `p` points to. If `p` is nil, the value will be discarded. Otherwise, the value
	// underlying `p` must be a pointer to the correct type for the next data item received.
	// If the input is at EOF, Decode returns io.EOF and does not modify p.
	DecodePrimitive(p interface{}) error

	// DecodeStruct reads the next value from the input stream and stores it in the struct `p`
	// points to. If `p` is nil, the value will be discarded. Otherwise, the value
	// underlying `p` must be a pointer to the correct type for the next data item received.
	// If the input is at EOF, Decode returns io.EOF and does not modify p.
	DecodeStruct(p StructPtr) error

	// DecodeStructSlice reads the next value from the input stream and stores it in the struct slice
	// `p` points to. If `p` is nil, the value will be discarded. Otherwise, the value
	// underlying `p` must be a pointer to the correct type for the next data item received.
	// If the input is at EOF, Decode returns io.EOF and does not modify p.
	DecodeStructSlice(p StructSlicePtr) error
}
