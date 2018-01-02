package drumbeat

import (
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/mattetti/audio/midi"
	"github.com/mattetti/filebuffer"
)

func TestFromMIDI(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns map[string]string
	}{
		{name: "single pattern", path: "fixtures/singlePattern.mid", patterns: map[string]string{"C1": "x.......x......."}},
		{name: "kick pattern", path: "fixtures/kick.mid", patterns: map[string]string{"C1": "x...x...x...x..."}},
		// TODO: this should automatically switch to 1/32th
		{name: "full beat", path: "fixtures/beat.mid", patterns: map[string]string{
			"C1":  "x...x.x.x...x...",
			"E1":  "............x...",
			"G#1": "xxx.x.xxxxxxxxxx", // 1/32th
			"D#1": "....x.......x...",
			"F#1": "x.x...x.x.x.x.x.",
		}},
		{name: "kick snare unquantized", path: "fixtures/kickSnare.mid", patterns: map[string]string{
			"C1": "x......x.......x.......x.......x.......x.......x.......x........",
			"D1": "...x.x.....x.x.....x.x.....x.x.....x.x.....x.x.....x.x.....x.x..",
		}},
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
			if len(patterns) != len(tt.patterns) {
				t.Fatalf("Expected %d patterns; got %d patterns", len(tt.patterns), len(patterns))
			}
			for _, p := range patterns {
				if tt.patterns[p.Name] != p.Pulses.String() {
					if len(tt.patterns[p.Name]) != len(p.Pulses) {
						t.Errorf("%s - expected %d Pulses, got %d", p.Name, len(tt.patterns[p.Name]), len(p.Pulses))
					}
					t.Errorf("Expected %s:\n%s, got\n%s", p.Name, tt.patterns[p.Name], p.Pulses)
				}
			}
		})
	}
}

func TestToMIDI(t *testing.T) {
	tests := []struct {
		name     string
		patterns map[string]string
		wantErr  bool
	}{
		// {name: "no patterns"},
		{name: "single pattern", patterns: map[string]string{"C1": "x...x...x...x..."}},
		{name: "last step is a pulse", patterns: map[string]string{"C1": "x...x..x"}},
		{name: "following pulses", patterns: map[string]string{"C1": "xxx.x..x"}},
		{name: "multiple patterns", patterns: map[string]string{"C1": "x...x...", "C#1": "..x...x."}},
		{name: "multiple similar patterns", patterns: map[string]string{"C1": "x...x...", "C#1": "x...x..."}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := filebuffer.New(nil)
			patterns := make([]*Pattern, len(tt.patterns))
			// startingKey := midi.KeyInt("C", 1)
			var i int
			for strKey, strPat := range tt.patterns {
				patterns[i] = NewFromString(One8, strPat)[0]
				oct, _ := strconv.Atoi(strKey[len(strKey)-1:])
				patterns[i].Key = midi.KeyInt(strKey[:len(strKey)-1], oct)
				i++
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
				// t.Logf("Got: %#v\n", extr)
				if extr.Pulses.String() != tt.patterns[extr.Name] {
					t.Errorf("Expected pattern %d to look like %s but got %s", i, tt.patterns[extr.Name], extr.Pulses.String())
				}
			}
		})
	}

	t.Run("autoset key", func(t *testing.T) {
		patterns := NewFromString(One8, "x.......;....x...")
		buf := filebuffer.New(nil)
		if err := ToMIDI(buf, patterns...); err != nil {
			t.Fatalf("ToMIDI() error = %v", err)
		}
		if len(patterns) != 2 {
			t.Fatalf("expected to have 2 patterns from string")
		}
		buf.Seek(0, io.SeekStart)
		extractedPatterns, err := FromMIDI(buf)
		if err != nil {
			t.Fatalf("FromMIDI failed to decode - %s", err)
		}
		if len(extractedPatterns) != 2 {
			t.Fatalf("expected to have extracted 2 patterns")
		}
		if extractedPatterns[0].Key != 36 {
			t.Fatalf("expected to the first extracted pattern to be set to C1 but was %v", extractedPatterns[0].Key)
		}
		if extractedPatterns[1].Key != 37 {
			t.Fatalf("expected to the first extracted pattern to be set to C#1 but was %v", extractedPatterns[1].Key)
		}
	})

}
