package cborgen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/daotl/go-marsha/cborgen"
	"github.com/daotl/go-marsha/test"
)

func TestMarsha(t *testing.T) {
	assert := assert.New(t)
	mer := cborgen.New()
	s := &test.TestStruct{"test"}
	ss := &test.TestStructs{test.TestStruct{"test"}, test.TestStruct{"test2"}}

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
		s2 := &test.TestStruct{}
		err = mer.UnmarshalStruct(bin, s2)
		assert.NoError(err)
		assert.Equal(s.Data, s2.Data)
	})

	t.Run("UnmarshalStruct error: model type does not match", func(t *testing.T) {
		bin, err := mer.MarshalStruct(s)
		assert.NoError(err)
		s2 := &test.TestStruct2{}
		err = mer.UnmarshalStruct(bin, s2)
		assert.Equal(cborgen.ErrTypeNotMatch, err)
	})

	t.Run("MarshalStructSlice/UnmarshalStructSlice", func(t *testing.T) {
		bin, err := mer.MarshalStructSlice(ss)
		assert.NoError(err)
		ss2 := &test.TestStructs{}
		err = mer.UnmarshalStructSlice(bin, ss2)
		assert.NoError(err)
		assert.Equal(ss, ss2)
	})
}
