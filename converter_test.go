package drumbeat

import (
	"testing"

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
			w := filebuffer.New(nil)
			patterns := make([]*Pattern, len(tt.patterns))
			for i, strPat := range tt.patterns {
				patterns[i] = NewFromString(strPat)
			}
			if err := ToMIDI(w, patterns...); (err != nil) != tt.wantErr {
				t.Errorf("ToMIDI() error = %v, wantErr %v", err, tt.wantErr)
			}
			// TODO: parse the MIDI and compare the patterns
		})
	}
}
