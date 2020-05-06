package pbts

import (
	"bytes"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
)

type TestStruct struct {
	Field      *string          `protobuf:"bytes,1,opt,name=field" json:"field,omitempty"`
	FieldInt   *int64           `protobuf:"varint,6,opt,name=field_int,json=fieldInt" json:"field_int,omitempty"`
	OtherField string           `json:"other_field"`
	Metadata   map[string]int32 `protobuf:"bytes,11,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	Tags       []string         `protobuf:"bytes,12,rep,name=tags" json:"tags,omitempty"`
}

func TestOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(
		TestStruct{},
	)
	g.Write()

	expected := filePreamble + `
export abstract class TestStruct {
  field?: string;
  fieldInt?: string;
  other_field?: string;
  metadata?: { [key: string]: number; };
  tags?: string[];
  static copy(from: TestStruct, to?: TestStruct): TestStruct {
    to = to || {};
    to.field = from.field;
    to.fieldInt = from.fieldInt;
    to.other_field = from.other_field;
    to.metadata = from.metadata;
    to.tags = from.tags;
    return to;
  }
}`

	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expected) {
		t.Error(buf.String())
		t.Error(expected)
	}
}

type TestEnum int32

func (TestEnum) EnumDescriptor() ([]byte, []int) {
	return nil, nil
}

const (
	TestEnum_unknown TestEnum = 0
	TestEnum_foo     TestEnum = 1
	TestEnum_bar     TestEnum = 2
)

var TestEnum_name = map[int32]string{
	0: "unknown",
	1: "foo",
	2: "bar",
}
var TestEnum_value = map[string]int32{
	"unknown": 0,
	"foo":     1,
	"bar":     2,
}

type TestEnumStruct struct {
	EnumField *TestEnum `protobuf:"varint,1,opt,name=test_enum,json=testEnum,enum=pbts.TestEnum" json:"test_enum,omitempty"`
}

func init() {
	proto.RegisterEnum("pbts.TestEnum", TestEnum_name, TestEnum_value)
}

func TestEnumOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(
		TestEnumStruct{},
	)
	g.Write()

	expected := filePreamble + `
export abstract class TestEnumStruct {
  testEnum?: TestEnum;
  static copy(from: TestEnumStruct, to?: TestEnumStruct): TestEnumStruct {
    to = to || {};
    to.testEnum = from.testEnum;
    return to;
  }
}

export type TestEnum = 'bar' | 'foo' | 'unknown';
`

	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expected) {
		t.Error(buf.String())
		t.Error(expected)
	}
}

func TestNativeEnumOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.NativeEnums = true
	g.RegisterMany(
		TestEnumStruct{},
	)
	g.Write()

	expected := filePreamble + `
export abstract class TestEnumStruct {
  testEnum?: TestEnum;
  static copy(from: TestEnumStruct, to?: TestEnumStruct): TestEnumStruct {
    to = to || {};
    to.testEnum = from.testEnum;
    return to;
  }
}

export enum TestEnum {
  Bar = "bar",
  Foo = "foo",
  Unknown = "unknown",
}`

	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expected) {
		t.Error(buf.String())
		t.Error(expected)
	}
}
