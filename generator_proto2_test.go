package pbts

import (
	"bytes"
	"embed"
	"testing"

	"github.com/octavore/pbts/v2/internal/test"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

//go:embed fixtures
var fixtures embed.FS

func loadFixture(name string) string {
	data, err := fixtures.ReadFile("fixtures/" + name)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func TestOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(&test.TestMessage{})
	g.Write()

	expected := loadFixture("TestMessage.ts")
	assert.Equal(t, expected, buf.String())
}

func TestEnumOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(&test.TestEnumStruct{})
	g.Write()

	expected := loadFixture("TestEnumStruct.ts")
	assert.Equal(t, expected, buf.String())
}

func TestNativeEnumOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.NativeEnums = true
	g.RegisterMany(&test.TestEnumStruct{})
	g.Write()

	expected := loadFixture("TestNativeEnumStruct.ts")
	assert.Equal(t, expected, buf.String())
}

func TestNewProtoStructOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.NativeEnums = true
	g.RegisterMany(&structpb.Struct{})
	g.Write()

	expected := loadFixture("TestProtoStruct.ts")
	assert.Equal(t, expected, buf.String())
}

func TestOneofOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(&test.TestOneofStruct{})
	g.Write()

	expected := loadFixture("TestOneOfStruct.ts")
	assert.Equal(t, expected, buf.String())
}
