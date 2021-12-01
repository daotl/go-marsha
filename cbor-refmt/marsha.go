// Package cbor_refmt provides a `Marsha` implementation for CBOR backed by go-ipld-cbor and
// refmt packages.
package cbor_refmt

import (
	"io"
	"sync"

	"github.com/daotl/go-marsha"
)

// Marsha is a Marsha implementation for CBOR backed by `go-ipld-cbor` and `refmt` packages.
//
// A Struct must first be registered by calling Marsha.Register(Struct{}) before being
// able to be marshaled/unmarshaled.
//
// Marshaling/unmarshaling can be customized by `refmt` tags:
//
//   type Model struct {
//   	Foo string `refmt:"bar,omitempty"`
//   }
//
type Marsha struct{}

var _ marsha.Marsha = (*Marsha)(nil)

// New creates a Marsha.
func New() *Marsha {
	return &Marsha{}
}

// Register a Struct type by passing empty a Struct.
func (m *Marsha) Register(i interface{}) {
	registerCborType(i)
}

func (m *Marsha) MarshalPrimitive(p interface{}) ([]byte, error) {
	return marshaller.Marshal(p)
}

func (m *Marsha) MarshalStruct(p marsha.StructPtr) ([]byte, error) {
	return marshaller.Marshal(p)
}

func (m *Marsha) MarshalStructSlice(p marsha.StructSlicePtr) ([]byte, error) {
	return marshaller.Marshal(p)
}

func (m *Marsha) UnmarshalStruct(bin []byte, p marsha.StructPtr) error {
	return unmarshaller.Unmarshal(bin, p)
}

func (m *Marsha) UnmarshalPrimitive(bin []byte, p interface{}) error {
	return unmarshaller.Unmarshal(bin, p)
}

func (m *Marsha) UnmarshalStructSlice(bin []byte, p marsha.StructSlicePtr) error {
	return unmarshaller.Unmarshal(bin, p)
}

// Not implemented
func (m *Marsha) NewEncoder(w io.Writer) marsha.Encoder {
	enc := new(encoder)
	enc.w = w
	return enc
}

// Not implemented
func (m *Marsha) NewDecoder(r io.Reader) marsha.Decoder {
	return &decoder{
		r: r,
	}
}

// An encoder manages the transmission of type and data information to the
// other side of a connection.  It is safe for concurrent use by multiple
// goroutines.
type encoder struct {
	sync.Mutex           // each item must be sent atomically
	w          io.Writer // where to send the data
}

func (enc *encoder) EncodePrimitive(p interface{}) error {
	return enc.encode(p)
}
func (enc *encoder) EncodeStruct(p marsha.StructPtr) error {
	return enc.encode(p)
}
func (enc *encoder) EncodeStructSlice(p marsha.StructSlicePtr) error {
	return enc.encode(p)
}

func (enc *encoder) encode(p interface{}) error {
	// Make sure we're single-threaded through here, so multiple
	// goroutines can share an encoder.
	enc.Lock()
	defer enc.Unlock()

	return marshaller.Encode(p, enc.w)
}

type decoder struct {
	sync.Mutex           // each item must be sent atomically
	r          io.Reader // where to receive the data
}

func (d *decoder) DecodePrimitive(p interface{}) error {
	return d.decode(p)
}

func (d *decoder) DecodeStruct(p marsha.StructPtr) error {
	return d.decode(p)
}

func (d *decoder) DecodeStructSlice(p marsha.StructSlicePtr) error {
	return d.decode(p)
}

func (d *decoder) decode(p interface{}) error {
	// Make sure we're single-threaded through here, so multiple
	// goroutines can share an encoder.
	d.Lock()
	defer d.Unlock()

	return unmarshaller.Decode(d.r, p)
}
