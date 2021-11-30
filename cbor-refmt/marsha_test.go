package cbor_refmt_test

import (
	"testing"

	"github.com/daotl/go-marsha"
	"github.com/stretchr/testify/assert"

	cbor_refmt "github.com/daotl/go-marsha/cbor-refmt"
)

type TestStructNoGen struct {
	Data string
}

func (s TestStructNoGen) Ptr() marsha.StructPtr { return &s }
func (s *TestStructNoGen) Val() marsha.Struct   { return *s }

type TestStructsNoGen []TestStructNoGen

func (s *TestStructsNoGen) Val() []marsha.Struct {
	models := make([]marsha.Struct, 0, len(*s))
	for _, m := range *s {
		models = append(models, m)
	}
	return models
}

type TestStruct2NoGen struct {
	Data2 int64
}

func (s TestStruct2NoGen) Ptr() marsha.StructPtr { return &s }
func (s *TestStruct2NoGen) Val() marsha.Struct   { return *s }

func TestMarsha(t *testing.T) {
	assert := assert.New(t)
	mer := cbor_refmt.New()
	s := &TestStructNoGen{"test"}
	ss := &TestStructsNoGen{TestStructNoGen{"test"}, TestStructNoGen{"test2"}}

	t.Run("Marhsal error: model type not registered", func(t *testing.T) {
		_, err := mer.MarshalStruct(s)
		assert.Error(err)
	})

	mer.Register(TestStructNoGen{})

	t.Run("MarshalPrimitive/MarshalPrimitive primitives", func(t *testing.T) {
		v1 := 52
		bin, err := mer.MarshalPrimitive(&v1)
		assert.NoError(err)
		v2 := 0
		err = mer.UnmarshalPrimitive(bin, &v2)
		assert.NoError(err)
		assert.Equal(v1, v2)
	})

	t.Run("MarshalPrimitive/MarshalPrimitive primitive slices", func(t *testing.T) {
		s1 := []int{4, 13}
		bin, err := mer.MarshalPrimitive(&s1)
		assert.NoError(err)
		var s2 []int
		err = mer.UnmarshalPrimitive(bin, &s2)
		assert.NoError(err)
		assert.Equal(s1, s2)
	})

	t.Run("MarshalStruct/UnmarshalStruct", func(t *testing.T) {
		bin, err := mer.MarshalStruct(s)
		assert.NoError(err)
		s2 := &TestStructNoGen{}
		err = mer.UnmarshalStruct(bin, s2)
		assert.NoError(err)
		assert.Equal(s.Data, s2.Data)
	})

	t.Run("UnmarshalStruct error: model type does not match", func(t *testing.T) {
		mer.Register(TestStruct2NoGen{})
		bin, err := mer.MarshalStruct(s)
		assert.NoError(err)
		s2 := &TestStruct2NoGen{}
		err = mer.UnmarshalStruct(bin, s2)
		assert.Error(err)
	})

	t.Run("MarshalStructSlice/UnmarshalStructSlice", func(t *testing.T) {
		bin, err := mer.MarshalStructSlice(ss)
		assert.NoError(err)
		ss2 := &TestStructsNoGen{}
		err = mer.UnmarshalStructSlice(bin, ss2)
		assert.NoError(err)
		assert.Equal(ss, ss2)
	})
}
