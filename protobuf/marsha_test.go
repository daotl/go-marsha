package protobuf_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	"github.com/daotl/go-marsha"
	"github.com/daotl/go-marsha/protobuf"
)

type TestStruct struct {
	*protobuf.Test
	Data string
}

func (s TestStruct) Ptr() marsha.StructPtr { return &s }
func (s *TestStruct) Val() marsha.Struct   { return *s }

func (s *TestStruct) EmptyPB() proto.Message {
	return &protobuf.Test{}
}

func (s *TestStruct) LoadPB(pb proto.Message) error {
	tpb := pb.(*protobuf.Test)
	s.Test = tpb
	s.Data = tpb.Data
	return nil
}

func (s *TestStruct) PB() proto.Message {
	if s.Test == nil {
		s.Test = &protobuf.Test{}
	}
	s.Test.Data = s.Data
	return s.Test
}

type TestStruct2 struct {
	*protobuf.Test2
	Data2 int32
}

func (s TestStruct2) Ptr() marsha.StructPtr { return &s }
func (s *TestStruct2) Val() marsha.Struct   { return *s }

func (s *TestStruct2) EmptyPB() proto.Message {
	// Deliberate wrong type
	return &protobuf.Test{}
}

func (s *TestStruct2) LoadPB(pb proto.Message) error {
	tpb := pb.(*protobuf.Test2)
	s.Test2 = tpb
	s.Data2 = tpb.Data2
	return nil
}

func (s *TestStruct2) PB() proto.Message {
	if s.Test2 == nil {
		s.Test2 = &protobuf.Test2{}
	}
	s.Test2.Data2 = s.Data2
	return s.Test2
}

func TestPBMarshaler(t *testing.T) {
	assert := assert.New(t)
	mer := protobuf.New()
	s := &TestStruct{&protobuf.Test{}, "test"}

	t.Run("Success", func(t *testing.T) {
		bin, err := mer.MarshalStruct(s)
		assert.NoError(err)
		n := &TestStruct{}
		err = mer.UnmarshalStruct(bin, n)
		assert.NoError(err)
		assert.Equal(s.Data, n.Data)
	})

	t.Run("Error: wrong protocol buffers type", func(t *testing.T) {
		bin, err := mer.MarshalStruct(s)
		assert.NoError(err)
		n := &TestStruct2{}
		err = mer.UnmarshalStruct(bin, n)
		assert.Equal(protobuf.ErrWrongPBType, err)
	})
}
