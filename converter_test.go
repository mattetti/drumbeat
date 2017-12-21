package drumbeat

import (
	"io"
	"os"
	"testing"

	"github.com/mattetti/audio/midi"
	"github.com/mattetti/filebuffer"
)

func TestFromMIDI(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns []string
	}{
		{name: "single pattern", path: "fixtures/singlePattern.mid", patterns: []string{"x...x..."}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := os.Open(tt.path)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			patterns, err := FromMIDI(f)
			if err != nil {
				t.Fatal(err)
			}
			for i, p := range patterns {
				if tt.patterns[i] != p.Steps.String() {
					t.Errorf("Expected %s, got %s", tt.patterns[i], p.Steps)
				}
			}
		})
	}
}

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

			// debugging
			// of, err := os.Create("test.mid")
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// defer of.Close()
			// of.Write(buf.Bytes())
			// buf.Seek(0, io.SeekStart)

			extractedPatterns, err := FromMIDI(buf)
			if err != nil {
				t.Fatalf("FromMIDI failed to decode - %s", err)
			}
			if len(extractedPatterns) != len(patterns) {
				t.Errorf("Expected %d patterns, but got %d", len(patterns), len(extractedPatterns))
			}
			for i, extr := range extractedPatterns {
				t.Logf("Got: %#v\n", extr)
				if extr.Steps.String() != tt.patterns[i] {
					t.Errorf("Expected pattern %d to look like %s but got %s", i, tt.patterns[i], extr.Steps.String())
				}
			}
		})
	}
}
