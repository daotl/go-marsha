package cborgen_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/daotl/go-marsha/cborgen"
	"github.com/daotl/go-marsha/test"
)

func TestMarsha(t *testing.T) {
	mer := cborgen.New()
	test.SubTestAll(t, mer)
}

func TestSpecial(t *testing.T) {
	req := require.New(t)
	asrt := assert.New(t)
	mer := cborgen.New()
	s := &test.TestStruct{Data: "test"}

	t.Run("UnmarshalStruct error: model type does not match", func(t *testing.T) {
		bin, err := mer.MarshalStruct(s)
		req.NoError(err)
		s2 := &test.TestStruct2{}
		_, err = mer.UnmarshalStruct(bin, s2)
		asrt.Equal(cborgen.ErrTypeNotMatch, err)
	})
}
