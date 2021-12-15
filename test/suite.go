package test

import (
	"bytes"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func SubTestBasic(t *testing.T, mrsh marsha.Marsha) {
	req := require.New(t)
	asrt := assert.New(t)
	s := &TestStruct{"test"}
	ss := &TestStructs{TestStruct{"test"}, TestStruct{"test2"}}

	t.Run("MarshalPrimitive/Unmarshal primitives", func(t *testing.T) {
		v1 := 52
		bin, err := mrsh.MarshalPrimitive(&v1)
		req.NoError(err)
		v2 := 0
		read, err := mrsh.UnmarshalPrimitive(bin, &v2)
		req.NoError(err)
		asrt.Equal(v1, v2)
		if read != -1 {
			asrt.Equal(len(bin), read)
		}
	})

	t.Run("MarshalPrimitive/Unmarshal primitive slices", func(t *testing.T) {
		s1 := []int{4, 13}
		bin, err := mrsh.MarshalPrimitive(&s1)
		req.NoError(err)
		var s2 []int
		read, err := mrsh.UnmarshalPrimitive(bin, &s2)
		req.NoError(err)
		asrt.Equal(s1, s2)
		if read != -1 {
			asrt.Equal(len(bin), read)
		}
	})

	t.Run("MarshalStruct/Unmarshal", func(t *testing.T) {
		bin, err := mrsh.MarshalStruct(s)
		req.NoError(err)
		s2 := &TestStruct{}
		read, err := mrsh.UnmarshalStruct(bin, s2)
		req.NoError(err)
		asrt.Equal(s.Data, s2.Data)
		if read != -1 {
			asrt.Equal(len(bin), read)
		}
	})

	t.Run("MarshalStructSlice/Unmarshal", func(t *testing.T) {
		bin, err := mrsh.MarshalStructSlice(ss)
		req.NoError(err)
		ss2 := &TestStructs{}
		read, err := mrsh.UnmarshalStructSlice(bin, ss2)
		req.NoError(err)
		asrt.Equal(ss, ss2)
		if read != -1 {
			asrt.Equal(len(bin), read)
		}
	})
}

func SubTestEncoderDecoder(t *testing.T, m marsha.Marsha) {
	req := require.New(t)
	asrt := assert.New(t)
	s := &TestStruct{"test"}
	ss := &TestStructs{TestStruct{"test"}, TestStruct{"test2"}}

	t.Run("MarshalPrimitive/Unmarshal primitives", func(t *testing.T) {
		var buf bytes.Buffer
		enc := m.NewEncoder(&buf)
		dec := m.NewDecoder(&buf)

		v1 := 52
		err := enc.EncodePrimitive(&v1)
		req.NoError(err)
		bin := buf.Bytes()
		v2 := 0
		read, err := dec.DecodePrimitive(&v2)
		req.NoError(err)
		asrt.Equal(v1, v2)
		if read != -1 {
			asrt.Equal(len(bin), read)
		}
	})

	t.Run("MarshalPrimitive/Unmarshal primitive slices", func(t *testing.T) {
		var buf bytes.Buffer
		enc := m.NewEncoder(&buf)
		dec := m.NewDecoder(&buf)

		s1 := []int{4, 13}
		err := enc.EncodePrimitive(&s1)
		req.NoError(err)
		bin := buf.Bytes()
		var s2 []int
		read, err := dec.DecodePrimitive(&s2)
		req.NoError(err)
		asrt.Equal(s1, s2)
		if read != -1 {
			asrt.Equal(len(bin), read)
		}
	})

	t.Run("MarshalStruct/Unmarshal", func(t *testing.T) {
		var buf bytes.Buffer
		enc := m.NewEncoder(&buf)
		dec := m.NewDecoder(&buf)

		err := enc.EncodeStruct(s)
		req.NoError(err)
		bin := buf.Bytes()
		s2 := &TestStruct{}
		read, err := dec.DecodeStruct(s2)
		req.NoError(err)
		asrt.Equal(s.Data, s2.Data)
		if read != -1 {
			asrt.Equal(len(bin), read)
		}
	})

	t.Run("MarshalStructSlice/Unmarshal", func(t *testing.T) {
		var buf bytes.Buffer
		enc := m.NewEncoder(&buf)
		dec := m.NewDecoder(&buf)

		err := enc.EncodeStructSlice(ss)
		req.NoError(err)
		bin := buf.Bytes()
		ss2 := &TestStructs{}
		read, err := dec.DecodeStructSlice(ss2)
		req.NoError(err)
		asrt.Equal(ss, ss2)
		if read != -1 {
			asrt.Equal(len(bin), read)
		}
	})
}
