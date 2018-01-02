package drumbeat

import (
	"bytes"
	"reflect"
	"testing"
)

func TestWriteTo(t *testing.T) {
	tests := []struct {
		name     string
		patterns []*Pattern
		wantErr  bool
	}{
		{name: "nil pattern", patterns: nil},
		{name: "1 pattern", patterns: NewFromString(One8, "x...x...")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := WriteTo(w, tt.patterns...); (err != nil) != tt.wantErr {
				t.Errorf("WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			patterns, err := ReadFrom(w)
			if err != nil {
				t.Errorf("ReadFrom() error = %v", err)
				return
			}
			if !reflect.DeepEqual(patterns, tt.patterns) {
				for i, p := range patterns {
					if p != tt.patterns[i] {
						t.Logf("[%d] %v != %v", i, p, tt.patterns[i])
					}
				}
				t.Fatalf("Expected %#v to be %#v", patterns, tt.patterns)
			}
		})
	}
}
