package drumbeat

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/mattetti/audio/midi"
)

func TestPulses_Offset(t *testing.T) {
	tests := []struct {
		name   string
		pulses Pulses
		n      int
		want   Pulses
	}{
		{name: "shift 2 to the right", pulses: []float64{0.1, 0.2, 0.3, 0.4}, n: 2, want: []float64{0.3, 0.4, 0.1, 0.2}},
		{name: "shift 2 to the right once again", pulses: []float64{1.0, 0.0, 0.1, 0.0}, n: 2, want: []float64{0.1, 0.0, 1.0, 0.0}},
		{name: "shift by the length of the slice", pulses: []float64{1.0, 0.0}, n: 2, want: []float64{1.0, 0.0}},
		{name: "shift by more than the length of the slice", pulses: []float64{1.0, 0.0}, n: 3, want: []float64{0.0, 1.0}},
		{name: "shift by more than the length of the slice once again", pulses: []float64{0.1, 0.2, 0.3}, n: 5, want: []float64{0.2, 0.3, 0.1}},
		{name: "shift by huge number", pulses: []float64{0.1, 0.2, 0.3}, n: 42, want: []float64{0.1, 0.2, 0.3}},
		{name: "shift using a negative value to go the other way around", pulses: []float64{0.0, 0.1, 0.2, 0.3}, n: -2, want: []float64{0.2, 0.3, 0.0, 0.1}},
		{name: "shift negatively by more than the length of the slice", pulses: []float64{0.1, 0.2, 0.3}, n: -4, want: []float64{0.2, 0.3, 0.1}},
		{name: "shift negatively by a huge number", pulses: []float64{0.1, 0.2, 0.3}, n: -46, want: []float64{0.2, 0.3, 0.1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pulses.Offset(tt.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pulses.Offset() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestPulses_String(t *testing.T) {
	tests := []struct {
		pulses Pulses
		want   string
	}{
		{[]float64{0, 1, 0, 1}, ".X.X"},
		{[]float64{0, 0, 0, 0}, "...."},
		{[]float64{1, 1, 1, 1}, "XXXX"},
		{[]float64{1, 0.2, 1, 0.5}, "XXXX"},
		{[]float64{0.0, 0.01, 0.2, 0}, ".XX."},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			if got := tt.pulses.String(); got != tt.want {
				t.Errorf("Pulses.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFromString(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want *Pattern
	}{
		{name: "basic", str: "x...x...x...", want: &Pattern{
			StepDuration: midi.Dur8th,
			Steps:        []float64{1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0},
			Velocity:     []float64{0.9, 0.0, 0.0, 0.0, 0.9, 0.0, 0.0, 0.0, 0.9, 0.0, 0.0, 0.0}},
		},
		{name: "with uppercase X", str: "X...x...X...", want: &Pattern{
			StepDuration: midi.Dur8th,
			Steps:        []float64{1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0},
			Velocity:     []float64{0.9, 0.0, 0.0, 0.0, 0.9, 0.0, 0.0, 0.0, 0.9, 0.0, 0.0, 0.0}},
		},
		{name: "without dots", str: "X___x   X~~~", want: &Pattern{
			StepDuration: midi.Dur8th,
			Steps:        []float64{1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0},
			Velocity:     []float64{0.9, 0.0, 0.0, 0.0, 0.9, 0.0, 0.0, 0.0, 0.9, 0.0, 0.0, 0.0}},
		},
		{name: "blank", str: "blank", want: &Pattern{
			StepDuration: midi.Dur8th,
			Steps:        []float64{0.0, 0.0, 0.0, 0.0, 0.0},
			Velocity:     []float64{0.0, 0.0, 0.0, 0.0, 0.0}},
		},
		{name: "empty", str: "", want: &Pattern{
			StepDuration: midi.Dur8th,
			Steps:        []float64{},
			Velocity:     []float64{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFromString(tt.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}
