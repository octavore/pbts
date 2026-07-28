package pbts

import (
	"bytes"
	"strings"
	"testing"

	"github.com/octavore/pbts/v2/internal/test"
	"github.com/stretchr/testify/assert"
)

// GenerateAll registers every message in protoregistry.GlobalTypes, which includes nested message
// types (they're globally registered under their fully-qualified name).
// RegisterDescriptor also walks a message's own nested type declarations and adds them to
// indirectModels.
// This means nested messages end up in both g.models (once GenerateAll registers it directly) and
// g.indirectModels (once via its parent) and must be deduped.
func TestDedupNestedMessages(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(
		&test.TestOneofStruct{},
		// simulates GenerateAll also registering the nested type directly,
		// since it's independently present in the global registry.
		&test.TestOneofStruct_Stock{},
		&test.TestOneofStruct_Currency{},
	)
	g.Write()

	out := buf.String()
	assert.Equal(t, 1, strings.Count(out, "export interface TestOneofStruct_Stock {"))
	assert.Equal(t, 1, strings.Count(out, "export interface TestOneofStruct_Currency {"))
}

// A single enum type referenced from more than one field must only be
// emitted once.
func TestDedupEnums(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(&test.TestProto3DupEnumStruct{})
	g.Write()

	out := buf.String()
	assert.Equal(t, 1, strings.Count(out, "export type TestProto3EnumStruct_TestProto3Enum ="))
}
