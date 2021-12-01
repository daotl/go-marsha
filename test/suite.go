package test

import (
	"bytes"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/daotl/go-marsha"
)

var Subtests = []func(t *testing.T, mer marsha.Marsha){
	SubTestBasic,
	SubTestEncoderDecoder,
}

func SubTestAll(t *testing.T, mer marsha.Marsha) {
	for _, f := range Subtests {
		t.Run(getFunctionName(f), func(t *testing.T) {
			f(t, mer)
		})
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func SubTestBasic(t *testing.T, mer marsha.Marsha) {
	asrt := assert.New(t)
	s := &TestStruct{"test"}
	ss := &TestStructs{TestStruct{"test"}, TestStruct{"test2"}}

	t.Run("MarshalPrimitive/Unmarshal primitives", func(t *testing.T) {
		v1 := 52
		bin, err := mer.MarshalPrimitive(&v1)
		asrt.NoError(err)
		v2 := 0
		err = mer.UnmarshalPrimitive(bin, &v2)
		asrt.NoError(err)
		asrt.Equal(v1, v2)
	})

	t.Run("MarshalPrimitive/Unmarshal primitive slices", func(t *testing.T) {
		s1 := []int{4, 13}
		bin, err := mer.MarshalPrimitive(&s1)
		asrt.NoError(err)
		var s2 []int
		err = mer.UnmarshalPrimitive(bin, &s2)
		asrt.NoError(err)
		asrt.Equal(s1, s2)
	})

	t.Run("MarshalStruct/Unmarshal", func(t *testing.T) {
		bin, err := mer.MarshalStruct(s)
		asrt.NoError(err)
		s2 := &TestStruct{}
		err = mer.UnmarshalStruct(bin, s2)
		asrt.NoError(err)
		asrt.Equal(s.Data, s2.Data)
	})

	t.Run("MarshalStructSlice/Unmarshal", func(t *testing.T) {
		bin, err := mer.MarshalStructSlice(ss)
		asrt.NoError(err)
		ss2 := &TestStructs{}
		err = mer.UnmarshalStructSlice(bin, ss2)
		asrt.NoError(err)
		asrt.Equal(ss, ss2)
	})
}

func SubTestEncoderDecoder(t *testing.T, m marsha.Marsha) {
	asrt := assert.New(t)
	var network bytes.Buffer
	enc := m.NewEncoder(&network)
	dec := m.NewDecoder(&network)
	s := &TestStruct{"test"}
	ss := &TestStructs{TestStruct{"test"}, TestStruct{"test2"}}

	t.Run("MarshalPrimitive/Unmarshal primitives", func(t *testing.T) {
		v1 := 52
		err := enc.EncodePrimitive(&v1)
		asrt.NoError(err)
		v2 := 0
		err = dec.DecodePrimitive(&v2)
		asrt.NoError(err)
		asrt.Equal(v1, v2)
	})

	t.Run("MarshalPrimitive/Unmarshal primitive slices", func(t *testing.T) {
		s1 := []int{4, 13}
		err := enc.EncodePrimitive(&s1)
		asrt.NoError(err)
		var s2 []int
		err = dec.DecodePrimitive(&s2)
		asrt.NoError(err)
		asrt.Equal(s1, s2)
	})

	t.Run("MarshalStruct/Unmarshal", func(t *testing.T) {
		err := enc.EncodeStruct(s)
		asrt.NoError(err)
		s2 := &TestStruct{}
		err = dec.DecodeStruct(s2)
		asrt.NoError(err)
		asrt.Equal(s.Data, s2.Data)
	})

	t.Run("MarshalStructSlice/Unmarshal", func(t *testing.T) {
		err := enc.EncodeStructSlice(ss)
		asrt.NoError(err)
		ss2 := &TestStructs{}
		err = dec.DecodeStructSlice(ss2)
		asrt.NoError(err)
		asrt.Equal(ss, ss2)
	})
}
