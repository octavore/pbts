package pbts

import (
	"bytes"
	"testing"

	"github.com/octavore/pbts/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestProto3Message(t *testing.T) {
	buf := &bytes.Buffer{}
	g := NewGenerator(buf)
	g.RegisterMany(&test.TestProto3Message{})
	g.Write()

	expected := loadFixture("TestProto3Message.ts")
	assert.Equal(t, expected, buf.String())
}
