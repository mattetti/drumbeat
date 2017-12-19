package drumbeat

import (
	"fmt"
	"io"
	"testing"

	"github.com/mattetti/audio/midi"
	"github.com/mattetti/filebuffer"
)

func TestToMIDI(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		wantErr  bool
	}{
		{name: "no patterns"},
		{name: "single pattern", patterns: []string{"x...x..."}},
		{name: "last step is a pulse", patterns: []string{"x...x..x"}},
		{name: "multiple patterns", patterns: []string{"x...x...", "..x...x."}},
		{name: "multiple similar patterns", patterns: []string{"x...x...", "x...x..."}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := filebuffer.New(nil)
			patterns := make([]*Pattern, len(tt.patterns))
			startingKey := midi.KeyInt("C", 1)
			for i, strPat := range tt.patterns {
				patterns[i] = NewFromString(strPat)[0]
				patterns[i].Key = startingKey + i
			}
			if err := ToMIDI(buf, patterns...); (err != nil) != tt.wantErr {
				t.Errorf("ToMIDI() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(patterns) < 1 {
				return
			}
			// Verify the generated MIDI
			// Rewind the buffer
			buf.Seek(0, io.SeekStart)
			extractedPatterns, err := FromMIDI(buf)
			if err != nil {
				t.Fatalf("FromMIDI failed to decode - %s", err)
			}
			if len(extractedPatterns) != len(patterns) {
				t.Errorf("Expected %d patterns, but got %d", len(patterns), len(extractedPatterns))
			}
			fmt.Printf("%#v\n", extractedPatterns[0])
			for i, extr := range extractedPatterns {
				if extr.Steps.String() != tt.patterns[i] {
					t.Errorf("Expected pattern %d to look like %s but got %s", i, tt.patterns[i], extr.Steps.String())
				}
			}
		})
	}
}
