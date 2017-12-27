package drumbeat

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestPulses_Offset(t *testing.T) {
	onTheOne := &Pulse{Ticks: 0, Duration: 96 / 2, Velocity: 90}
	onTheOneDev2 := &Pulse{Ticks: 96 / 2, Duration: 96 / 2, Velocity: 90}
	// onTheTwo := &Pulse{Ticks: 96, Duration: 96 / 2, Velocity: 90}

	// TODO: check that thechange the start time of the events have been updated

	tests := []struct {
		name   string
		pulses Pulses
		n      int
		want   Pulses
	}{
		{name: "shift 2 to the right", pulses: []*Pulse{onTheOne, nil, nil, nil}, n: 2, want: []*Pulse{nil, nil, onTheOne, nil}},
		{name: "shift 2 to the right once again", pulses: []*Pulse{nil, onTheOne, onTheOneDev2, nil}, n: 2, want: []*Pulse{onTheOneDev2, nil, nil, onTheOne}},
		// {name: "shift by the length of the slice", pulses: []*Pulse{}, n: 2, want: []*Pulse{}},
		// {name: "shift by more than the length of the slice", pulses: []*Pulse{}, n: 3, want: []*Pulse{}},
		// {name: "shift by more than the length of the slice once again", pulses: []*Pulse{}, n: 5, want: []*Pulse{}},
		// {name: "shift by huge number", pulses: []*Pulse{}, n: 42, want: []*Pulse{}},
		// {name: "shift using a negative value to go the other way around", pulses: []*Pulse{}, n: -2, want: []*Pulse{}},
		// {name: "shift negatively by more than the length of the slice", pulses: []*Pulse{}, n: -4, want: []*Pulse{}},
		// {name: "shift negatively by a huge number", pulses: []*Pulse{}, n: -46, want: []*Pulse{}},
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
		{[]*Pulse{
			nil,
			&Pulse{Ticks: 0, Velocity: 99},
			nil,
			&Pulse{Ticks: 0, Velocity: 99},
		}, ".X.X"},
		{[]*Pulse{}, "...."},
		{[]*Pulse{}, "XXXX"},
		{[]*Pulse{}, "XXXX"},
		{[]*Pulse{}, ".XX."},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			want := strings.ToLower(tt.want)
			if got := tt.pulses.String(); got != want {
				t.Errorf("Pulses.String() = %v, want %v", got, want)
			}
		})
	}
}

func TestNewFromString(t *testing.T) {
	t.Skip()
	/*
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
				if got := NewFromString(tt.str)[0]; !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewFromString() = %v, want %v", got, tt.want)
				}
			})
		}
	*/
}
