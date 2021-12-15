package cbor_refmt_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daotl/go-marsha"
	cbor_refmt "github.com/daotl/go-marsha/cbor-refmt"
	"github.com/daotl/go-marsha/test"
)

type TestStructNoGen struct {
	Data string
}

func (s TestStructNoGen) Ptr() marsha.StructPtr { return &s }
func (s *TestStructNoGen) Val() marsha.Struct   { return *s }

type TestStructsNoGen []TestStructNoGen

func (s *TestStructsNoGen) Val() []marsha.StructPtr {
	models := make([]marsha.StructPtr, 0, len(*s))
	for i := range *s {
		models = append(models, &(*s)[i])
	}
	return models
}

type TestStruct2NoGen struct {
	Data2 int64
}

func (s TestStruct2NoGen) Ptr() marsha.StructPtr { return &s }
func (s *TestStruct2NoGen) Val() marsha.Struct   { return *s }

// Some tests expected to fail for now, because `test.TestStruct` implements `cborMarshaler` interface
// from `github.com/ipfs/go-ipld-cbor/encoding` package, so marshaling will call `MarshalCBOR` which
// uses tuple to encode the outmost struct. `test.TestStruct` also implements github.com/daotl/cbor-gen.Unmarshaler`,
// which is modified and different from `cborUnmarshaler`, so unmarshaling won't call `UnmarshalCBOR`
// and try to use default map encoding for the outmost struct.
func TestSuite(t *testing.T) {
	mer := cbor_refmt.New()
	mer.Register(test.TestStruct{})
	test.SubTestAll(t, mer)
}

func TestNoGenBasic(t *testing.T) {
	req := require.New(t)
	asrt := assert.New(t)
	mer := cbor_refmt.New()
	mer.Register(test.TestStruct{})
	s := &TestStructNoGen{"test"}
	ss := &TestStructsNoGen{TestStructNoGen{"test"}, TestStructNoGen{"test2"}}

	t.Run("Marhsal error: model type not registered", func(t *testing.T) {
		_, err := mer.MarshalStruct(s)
		asrt.Error(err)
	})

	mer.Register(TestStructNoGen{})

	t.Run("MarshalStruct/UnmarshalStruct", func(t *testing.T) {
		bin, err := mer.MarshalStruct(s)
		asrt.NoError(err)
		s2 := &TestStructNoGen{}
		read, err := mer.UnmarshalStruct(bin, s2)
		req.NoError(err)
		asrt.Equal(s.Data, s2.Data)
		asrt.Equal(len(bin), read)
	})

	t.Run("UnmarshalStruct error: model type does not match", func(t *testing.T) {
		mer.Register(TestStruct2NoGen{})
		bin, err := mer.MarshalStruct(s)
		asrt.NoError(err)
		s2 := &TestStruct2NoGen{}
		read, err := mer.UnmarshalStruct(bin, s2)
		req.NoError(err)
		asrt.Equal(s.Data, s2.Data2)
		asrt.Equal(len(bin), read)
	})

	t.Run("MarshalStructSlice/UnmarshalStructSlice", func(t *testing.T) {
		bin, err := mer.MarshalStructSlice(ss)
		asrt.NoError(err)
		ss2 := &TestStructsNoGen{}
		read, err := mer.UnmarshalStructSlice(bin, ss2)
		req.NoError(err)
		asrt.Equal(ss, ss2)
		asrt.Equal(len(bin), read)
	})

	t.Run("Encoder_decoder NoGen", func(t *testing.T) {
		var network bytes.Buffer
		enc := mer.NewEncoder(&network)
		dec := mer.NewDecoder(&network)

		t.Run("MarshalStruct/Unmarshal", func(t *testing.T) {
			err := enc.EncodeStruct(s)
			asrt.NoError(err)
			s2 := &TestStructNoGen{}
			_, err = dec.DecodeStruct(s2)
			asrt.NoError(err)
			asrt.Equal(s.Data, s2.Data)
		})

		t.Run("MarshalStructSlice/Unmarshal", func(t *testing.T) {
			err := enc.EncodeStructSlice(ss)
			asrt.NoError(err)
			ss2 := &TestStructsNoGen{}
			_, err = dec.DecodeStructSlice(ss2)
			asrt.NoError(err)
			asrt.Equal(ss, ss2)
		})
	})
}
