package pbts

import (
	"bytes"
	"testing"

	"github.com/octavore/pbts/v2/internal/test"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestProto3Message(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(
		&test.TestProto3Message{},
		&test.TestProto3NestedMessage{},
	)
	g.Write()

	expected := loadFixture("TestProto3Message.ts")
	assert.Equal(t, expected, buf.String())
}

func TestProtoJSON(t *testing.T) {
	j := &protojson.MarshalOptions{
		EmitDefaultValues: true,
	}
	data, _ := j.Marshal(&test.TestProto3Message{})
	// unfortunately EmitDefaultValues will output [] and {} for
	// empty lists and maps respectively and `optional repeated`/`optional map` is not supported
	expected := loadFixture("proto3message.json")
	assert.JSONEq(t, expected, string(data))
}
