package test

import (
	"github.com/daotl/go-marsha"
	"github.com/daotl/go-marsha/cborgen"
)

var _ cborgen.StructPtr = (*TestStruct)(nil)

type TestStruct struct {
	Data string
}

func (s TestStruct) Ptr() marsha.StructPtr { return &s }
func (s *TestStruct) Val() marsha.Struct   { return *s }

type TestStructs []TestStruct

func (s *TestStructs) Val() []marsha.Struct {
	models := make([]marsha.Struct, 0, len(*s))
	for _, m := range *s {
		models = append(models, m)
	}
	return models
}
func (*TestStructs) NewStruct() marsha.Struct     { return TestStruct{} }
func (s *TestStructs) Append(m cborgen.StructPtr) { *s = append(*s, *(m.(*TestStruct))) }

type TestStruct2 struct {
	Data2 int64
}

func (s TestStruct2) Ptr() marsha.StructPtr { return &s }
func (s *TestStruct2) Val() marsha.Struct   { return *s }
