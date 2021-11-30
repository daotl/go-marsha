// Package cbor_refmt provides a `Marsha` implementation for CBOR backed by go-ipld-cbor and
// refmt packages.
package cbor_refmt

import (
	"bytes"
	"encoding/gob"
	"io"
	"reflect"
	"sync"
	_ "unsafe"

	"github.com/daotl/go-marsha"
	"github.com/ipfs/go-ipld-cbor/encoding"
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
	RegisterCborType(i)
}

func (m *Marsha) MarshalPrimitive(p interface{}) ([]byte, error) {
	return marshaller.Marshal(p)
}

func (m *Marsha) UnmarshalPrimitive(bin []byte, p interface{}) error {
	return unmarshaller.Decode(bytes.NewReader(bin), p)
}

func (m *Marsha) MarshalStruct(p marsha.StructPtr) ([]byte, error) {
	return marshaller.Marshal(p)
}

func (m *Marsha) UnmarshalStruct(bin []byte, p marsha.StructPtr) error {
	return unmarshaller.Decode(bytes.NewReader(bin), p)
}

func (m *Marsha) MarshalStructSlice(p marsha.StructSlicePtr) ([]byte, error) {
	return marshaller.Marshal(p)
}

func (m *Marsha) UnmarshalStructSlice(bin []byte, p marsha.StructSlicePtr) error {
	return unmarshaller.Decode(bytes.NewReader(bin), p)
}

// Not implemented
func (m *Marsha) NewEncoder(w io.Writer) *marsha.Encoder {
	enc := new(Encoder)
	enc.w = w
	enc.countState = enc.newEncoderState(new(encBuffer))
	enc.marshaller =
	return enc
}

// Not implemented
func (m *Marsha) NewDecoder(_ io.Reader) *marsha.Decoder {
	panic(marsha.ErrUnimplemented)
}

// An Encoder manages the transmission of type and data information to the
// other side of a connection.  It is safe for concurrent use by multiple
// goroutines.
type Encoder struct {
	mutex      sync.Mutex              // each item must be sent atomically
	w          io.Writer             // where to send the data
	countState *encoderState           // stage for writing counts
	freeList   *encoderState           // list of free encoderStates; avoids reallocation
	byteBuf    encBuffer               // buffer for top-level encoderState
	err        error
}

func (enc *Encoder) EncodeEncodePrimitive(p interface{}) error {
	return enc.encode(p)
}
func (enc *Encoder) EncodeStruct(p marsha.StructPtr) error {
	return enc.encode(p)
}
func (enc *Encoder) EncodeStructSlice(p marsha.StructSlicePtr) error {
	return enc.encode(p)
}

func (enc *Encoder) encode(p interface{}) error {
	// Make sure we're single-threaded through here, so multiple
	// goroutines can share an encoder.
	enc.mutex.Lock()
	defer enc.mutex.Unlock()

	gob.NewEncoder(nil).

	m := marshaller.pool.Get().(*encoding.Marshaller)
	err := m.Encode(obj, w)
	p.pool.Put(m)
	return err

	// Encode the object.
	enc.encode(state.b, value, ut)
	if enc.err == nil {
		enc.writeMessage(enc.writer(), state.b)
	}

	enc.freeEncoderState(state)
	return enc.err
}

