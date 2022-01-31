// Package cborgen provides a fast Marsha implementation for CBOR backed by `go-ipld-cbor` package
// and marshaling/unmarshaling code generated by github.com/daotl/cbor-gen package.

package cborgen

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/daotl/go-marsha"
	"github.com/daotl/go-marsha/internal/refmt"
)

var (
	ErrNotCBORStructPtr      = errors.New("not a cbor.StructPtr")
	ErrNotCBORStructSlicePtr = errors.New("not a cbor.StructSlicePtr")
	ErrTypeNotMatch          = errors.New("model type does not match")
)

// StructPtr implements `github.com/ipfs-ipld-cbor/encoding.cborMarshaler`
type StructPtr interface {
	marsha.StructPtr

	// MarshalCBOR should be generated by `github.com/daotl/cbor-gen` package.
	// This implements `github.com/daotl/cbor-gen.CBORMarshaler`
	MarshalCBOR(w io.Writer) (int, error)

	// UnmarshalCBOR should be generated by `github.com/daotl/cbor-gen` package.
	// This implements `github.com/daotl/cbor-gen.CBORUnmarshaler`
	UnmarshalCBOR(r io.Reader) (int, error)
}

type StructSlicePtr interface {
	marsha.StructSlicePtr

	// NewStruct should return an empty marsha.Struct.
	NewStructPtr() marsha.StructPtr

	// Append should append `p.Val()` to the `StructSlice` this `StructSlicePtr` points to.
	Append(p StructPtr)
}

// Marsha is a fast Marsha implementation for CBOR backed by `go-ipld-cbor` package
// and marshaling/unmarshaling code generated by github.com/daotl/cbor-gen package.
type Marsha struct {
	refmt *refmt.Refmt
}

var _ marsha.Marsha = (*Marsha)(nil)

// New creates a Marsha.
func New() *Marsha {
	return &Marsha{
		refmt: refmt.New(),
	}
}

// Register a Struct type by passing empty a Struct.
func (m *Marsha) Register(i interface{}) {
	m.refmt.RegisterCborType(i)
}

func (m *Marsha) MarshalPrimitive(p interface{}) ([]byte, error) {
	return m.refmt.Marshaller.Marshal(p)
}

// This implementation does not support returning the count of bytes read.
func (m *Marsha) UnmarshalPrimitive(bin []byte, p interface{}) (int, error) {
	return -1, m.refmt.Unmarshaller.Unmarshal(bin, p)
}

func (m *Marsha) MarshalStruct(p marsha.StructPtr) ([]byte, error) {
	cbp, ok := p.(StructPtr)
	if !ok {
		return nil, ErrNotCBORStructPtr
	}
	var buf bytes.Buffer
	if _, err := cbp.MarshalCBOR(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *Marsha) UnmarshalStruct(bin []byte, p marsha.StructPtr) (int, error) {
	cbp, ok := p.(StructPtr)
	if !ok {
		return 0, ErrNotCBORStructPtr
	}
	return unmarshal(bytes.NewReader(bin), cbp)
}

func (m *Marsha) MarshalStructSlice(p marsha.StructSlicePtr) (bin []byte, err error) {
	cbp, ok := p.(StructSlicePtr)
	if !ok {
		return nil, ErrNotCBORStructSlicePtr
	}

	ptrs := cbp.Val()
	bins := make([][]byte, len(ptrs))
	l := 0
	for i, s := range ptrs {
		if bins[i], err = m.MarshalStruct(s); err != nil {
			return nil, err
		}
		l += len(bins[i])
	}

	//binH := cbg.CborEncodeMajorType(cbg.MajArray, uint64(len(p.Val())))
	bin = make([]byte, 0 /*len(binH)+*/, l)
	//bin = append(bin, binH...)
	for _, b := range bins {
		bin = append(bin, b...)
	}
	return bin, nil
}

func (m *Marsha) UnmarshalStructSlice(bin []byte, p marsha.StructSlicePtr) (int, error) {
	cbp, ok := p.(StructSlicePtr)
	if !ok {
		return 0, ErrNotCBORStructSlicePtr
	}

	r := bytes.NewReader(bin)
	bytesRead := 0
	for {
		s, ok := cbp.NewStructPtr().(StructPtr)
		if !ok {
			return 0, ErrNotCBORStructPtr
		}
		if read, err := unmarshal(r, s); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return bytesRead, err
		} else {
			bytesRead += read
		}
		cbp.Append(s)
	}
	return bytesRead, nil
}

func (m *Marsha) NewEncoder(w io.Writer) marsha.Encoder {
	return &encoder{
		refmt: m.refmt,
		w:     w,
	}
}

func (m *Marsha) NewDecoder(r io.Reader) marsha.Decoder {
	return &decoder{
		refmt: m.refmt,
		r:     r,
	}
}

type encoder struct {
	sync.Mutex // each item must be sent atomically
	refmt      *refmt.Refmt
	w          io.Writer
}

func (e *encoder) EncodePrimitive(p interface{}) (int, error) {
	// Make sure we're single-threaded through here, so multiple
	// goroutines can share an encoder.
	e.Lock()
	defer e.Unlock()
	//return e.refmt.Marshaller.Encode(p, e.w)
	if bin, err := e.refmt.Marshaller.Marshal(p); err != nil {
		return 0, err
	} else {
		return len(bin), e.refmt.Marshaller.Encode(p, e.w)
	}
}

func (e *encoder) EncodeStruct(p marsha.StructPtr) (int, error) {
	cbp, ok := p.(StructPtr)
	if !ok {
		return 0, ErrNotCBORStructPtr
	}
	e.Lock()
	defer e.Unlock()
	return cbp.MarshalCBOR(e.w)
}

func (e *encoder) EncodeStructSlice(p marsha.StructSlicePtr) (n int, err error) {
	cbp, ok := p.(StructSlicePtr)
	if !ok {
		return 0, ErrNotCBORStructSlicePtr
	}
	e.Lock()
	defer e.Unlock()

	for _, s := range cbp.Val() {
		elem, ok := s.(StructPtr)
		if !ok {
			return n, ErrNotCBORStructPtr
		}
		n_, err := elem.MarshalCBOR(e.w)
		n += n_
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

type decoder struct {
	sync.Mutex // each item must be sent atomically
	refmt      *refmt.Refmt
	r          io.Reader
}

// This implementation does not support returning the count of bytes read.
func (d *decoder) DecodePrimitive(p interface{}) (int, error) {
	// Make sure we're single-threaded through here, so multiple
	// goroutines can share a decoder.
	d.Lock()
	defer d.Unlock()
	return -1, d.refmt.Unmarshaller.Decode(d.r, p)
}

func (d *decoder) DecodeStruct(p marsha.StructPtr) (int, error) {
	cbp, ok := p.(StructPtr)
	if !ok {
		return 0, ErrNotCBORStructPtr
	}
	d.Lock()
	defer d.Unlock()
	return unmarshal(d.r, cbp)
}

func (d *decoder) DecodeStructSlice(p marsha.StructSlicePtr) (int, error) {
	cbp, ok := p.(StructSlicePtr)
	if !ok {
		return 0, ErrNotCBORStructSlicePtr
	}
	d.Lock()
	defer d.Unlock()

	bytesRead := 0
	for {
		s, ok := cbp.NewStructPtr().(StructPtr)
		if !ok {
			return 0, ErrNotCBORStructPtr
		}
		if read, err := unmarshal(d.r, s); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return bytesRead, err
		} else {
			bytesRead += read
		}
		cbp.Append(s)
	}
	return bytesRead, nil
}

func unmarshal(r io.Reader, p StructPtr) (read int, err error) {
	if read, err = p.UnmarshalCBOR(r); err != nil {
		if strings.Contains(err.Error(), "wrong type") {
			err = ErrTypeNotMatch
		}
	}
	return read, err
}
