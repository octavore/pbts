package pbts

import (
	"bytes"
	"strings"
	"testing"

	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/octavore/pbts/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(
		test.TestMessage{},
	)
	g.Write()

	expected := filePreamble + `
export abstract class TestMessage {
  strField?: string;
  int32Field?: number;
  int64Field?: string;
  strList?: string[];
  metadata?: {[key: string]: string};
  static copy(from: TestMessage, to?: TestMessage): TestMessage {
    to = to || {};
    to.strField = from.strField;
    to.int32Field = from.int32Field;
    to.int64Field = from.int64Field;
    to.strList = from.strList;
    to.metadata = from.metadata;
    return to;
  }
}

`
	assert.Equal(t, expected, buf.String())
}

func TestEnumOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(
		test.TestEnumStruct{},
	)
	g.Write()

	expected := filePreamble + `
export abstract class TestEnumStruct {
  enumField?: TestEnumStruct_TestEnum;
  static copy(from: TestEnumStruct, to?: TestEnumStruct): TestEnumStruct {
    to = to || {};
    to.enumField = from.enumField;
    return to;
  }
}

export type TestEnumStruct_TestEnum = 'bar' | 'foo' | 'unknown';
`

	assert.Equal(t, expected, buf.String())
}

func TestNativeEnumOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.NativeEnums = true
	g.RegisterMany(
		test.TestEnumStruct{},
	)
	g.Write()

	expected := filePreamble + `
export abstract class TestEnumStruct {
  enumField?: TestEnumStruct_TestEnum;
  static copy(from: TestEnumStruct, to?: TestEnumStruct): TestEnumStruct {
    to = to || {};
    to.enumField = from.enumField;
    return to;
  }
}

export enum TestEnumStruct_TestEnum {
  Bar = "bar",
  Foo = "foo",
  Unknown = "unknown",
}
`

	assert.Equal(t, expected, buf.String())
}

func TestProtoStructOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.NativeEnums = true
	g.RegisterMany(
		_struct.Struct{},
	)
	g.Write()

	expected := filePreamble + `
export abstract class Struct {
  fields?: {[key: string]: any};
  static copy(from: Struct, to?: Struct): Struct {
    to = to || {};
    to.fields = from.fields;
    return to;
  }
}`

	if strings.TrimSpace(buf.String()) != strings.TrimSpace(expected) {
		t.Error(buf.String())
		t.Error(expected)
	}
}

func TestOneofOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(
		test.TestOneofStruct{},
	)
	g.Write()

	expected := filePreamble + `
export abstract class TestOneofStruct {
  // skipped field: instrument

  // oneof types:
  currency?: TestOneofStruct_Currency;
  stock?: TestOneofStruct_Stock;
  static copy(from: TestOneofStruct, to?: TestOneofStruct): TestOneofStruct {
    to = to || {};
    if ('currency' in from) {
      to.currency = TestOneofStruct_Currency.copy(from.currency || {}, to.currency || {});
    }
    if ('stock' in from) {
      to.stock = TestOneofStruct_Stock.copy(from.stock || {}, to.stock || {});
    }
    return to;
  }
}


// oneof types
export enum TestOneofStruct_InstrumentOneOf {
  Currency = 'currency',
  Stock = 'stock',
}
`

	assert.Equal(t, expected, buf.String())
}
