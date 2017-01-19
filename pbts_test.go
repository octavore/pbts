package pbts

import (
	"bytes"
	"strings"
	"testing"
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
