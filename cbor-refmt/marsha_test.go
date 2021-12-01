package cbor_refmt_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestSuite(t *testing.T) {
	mer := cbor_refmt.New()
	mer.Register(test.TestStruct{})
	test.SubTestAll(t, mer)
}

func TestNoGenBasic(t *testing.T) {
	asrt := assert.New(t)
	mer := cbor_refmt.New()
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
		err = mer.UnmarshalStruct(bin, s2)
		asrt.NoError(err)
		asrt.Equal(s.Data, s2.Data)
	})

	t.Run("UnmarshalStruct error: model type does not match", func(t *testing.T) {
		mer.Register(TestStruct2NoGen{})
		bin, err := mer.MarshalStruct(s)
		asrt.NoError(err)
		s2 := &TestStruct2NoGen{}
		err = mer.UnmarshalStruct(bin, s2)
		asrt.Error(err)
	})

	t.Run("MarshalStructSlice/UnmarshalStructSlice", func(t *testing.T) {
		bin, err := mer.MarshalStructSlice(ss)
		asrt.NoError(err)
		ss2 := &TestStructsNoGen{}
		err = mer.UnmarshalStructSlice(bin, ss2)
		asrt.NoError(err)
		asrt.Equal(ss, ss2)
	})

	t.Run("Encoder_decoder NoGen", func(t *testing.T) {
		var network bytes.Buffer
		enc := mer.NewEncoder(&network)
		dec := mer.NewDecoder(&network)

		t.Run("MarshalStruct/Unmarshal", func(t *testing.T) {
			err := enc.EncodeStruct(s)
			asrt.NoError(err)
			s2 := &TestStructNoGen{}
			err = dec.DecodeStruct(s2)
			asrt.NoError(err)
			asrt.Equal(s.Data, s2.Data)
		})

		t.Run("MarshalStructSlice/Unmarshal", func(t *testing.T) {
			err := enc.EncodeStructSlice(ss)
			asrt.NoError(err)
			ss2 := &TestStructsNoGen{}
			err = dec.DecodeStructSlice(ss2)
			asrt.NoError(err)
			asrt.Equal(ss, ss2)
		})
	})
}
