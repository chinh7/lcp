package core

import (
	"testing"
)

func TestParseEvents(t *testing.T) {
	event := []byte("AB")
	eventBytes := append(make([][]byte, 0), event)
	parseEvents(eventBytes)
	if string(events[0].Attributes[0].Value) != string(event) {
		t.Errorf("Expect event to be %s, got %s", string(event), events[0].Attributes[0].Value)
	}
}
