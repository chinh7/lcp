package core

import (
	"testing"
)

func TestParseEvent(t *testing.T) {
	attribute := []byte("AB")
	attributeBytes := append(make([][]byte, 0), attribute)
	event := parseEvent(attributeBytes)
	if string(event.Attributes[0].Value) != string(attribute) {
		t.Errorf("Expect event to be %s, got %s", string(attribute), event.Attributes[0].Value)
	}
}
