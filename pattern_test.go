package drumbeat

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestPulses_Offset(t *testing.T) {
	onTheOne := &Pulse{Ticks: 0, Duration: 96 / 2, Velocity: 90}
	onTheOneDiv2 := &Pulse{Ticks: 96 / 2, Duration: 96 / 2, Velocity: 90}
	// onTheTwo := &Pulse{Ticks: 96, Duration: 96 / 2, Velocity: 90}

	// TODO: check that thechange the start time of the events have been updated

	tests := []struct {
		name   string
		pulses Pulses
		n      int
		want   Pulses
	}{
		{name: "shift 2 to the right", pulses: []*Pulse{onTheOne, nil, nil, nil}, n: 2, want: []*Pulse{nil, nil, onTheOne, nil}},
		{name: "shift 2 to the right once again", pulses: []*Pulse{nil, onTheOne, onTheOneDiv2, nil}, n: 2, want: []*Pulse{onTheOneDiv2, nil, nil, onTheOne}},
		{name: "shift by the length of the slice", pulses: []*Pulse{onTheOne, onTheOneDiv2}, n: 2, want: []*Pulse{onTheOne, onTheOneDiv2}},
		// {name: "shift by more than the length of the slice", pulses: []*Pulse{}, n: 3, want: []*Pulse{}},
		// {name: "shift by more than the length of the slice once again", pulses: []*Pulse{}, n: 5, want: []*Pulse{}},
		// {name: "shift by huge number", pulses: []*Pulse{}, n: 42, want: []*Pulse{}},
		// {name: "shift using a negative value to go the other way around", pulses: []*Pulse{}, n: -2, want: []*Pulse{}},
		// {name: "shift negatively by more than the length of the slice", pulses: []*Pulse{}, n: -4, want: []*Pulse{}},
		// {name: "shift negatively by a huge number", pulses: []*Pulse{}, n: -46, want: []*Pulse{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pat := &Pattern{PPQN: DefaultPPQN, Grid: One16, Pulses: tt.pulses}
			pat.Offset(tt.n)
			if !reflect.DeepEqual(pat.Pulses, tt.want) {
				t.Errorf("Pulses.Offset() = %#v, want %#v", pat.Pulses, tt.want)
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
		{[]*Pulse{nil, &Pulse{}, nil, nil, nil, nil, nil, nil}, "........"},
		{[]*Pulse{
			&Pulse{Ticks: 0, Velocity: 90}, &Pulse{Ticks: 24, Velocity: 90}, &Pulse{Ticks: 48, Velocity: 90}, &Pulse{Ticks: 72, Velocity: 90},
			&Pulse{Ticks: 96, Velocity: 90}, &Pulse{Ticks: 120, Velocity: 90}, &Pulse{Ticks: 144, Velocity: 90}, &Pulse{Ticks: 168, Velocity: 90},
			&Pulse{Ticks: 192, Velocity: 90}, &Pulse{Ticks: 216, Velocity: 90}, &Pulse{Ticks: 240, Velocity: 90}, &Pulse{Ticks: 264, Velocity: 90},
			&Pulse{Ticks: 288, Velocity: 90}, &Pulse{Ticks: 312, Velocity: 90}, &Pulse{Ticks: 336, Velocity: 90}, &Pulse{Ticks: 360, Velocity: 90}},
			"XXXXXXXXXXXXXXXX"},
		{[]*Pulse{&Pulse{Ticks: 0, Velocity: 90}, &Pulse{Ticks: 24, Velocity: 90}, &Pulse{Ticks: 48, Velocity: 90}, &Pulse{Ticks: 72, Velocity: 90}}, "XXXX"},
		{[]*Pulse{nil, &Pulse{Ticks: 24, Velocity: 90}, &Pulse{Ticks: 48, Velocity: 90}, nil}, ".XX."},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s", tt.want), func(t *testing.T) {
			want := strings.ToLower(tt.want)
			if got := tt.pulses.String(); got != want {
				t.Errorf("Pulses.String() = %v, want %v", got, want)
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
		{name: "basic", str: "x...x...x...x...", want: &Pattern{
			Grid: One16,
			PPQN: 96,
			Pulses: []*Pulse{
				&Pulse{Ticks: 0, Velocity: 90}, nil, nil, nil,
				&Pulse{Ticks: 96, Velocity: 90}, nil, nil, nil,
				&Pulse{Ticks: 192, Velocity: 90}, nil, nil, nil,
				&Pulse{Ticks: 288, Velocity: 90}, nil, nil, nil}},
		},
		{name: "with uppercase X", str: "X...x...X...X...", want: &Pattern{
			Grid: One16,
			PPQN: 96,
			Pulses: []*Pulse{
				&Pulse{Ticks: 0, Velocity: 90}, nil, nil, nil,
				&Pulse{Ticks: 96, Velocity: 90}, nil, nil, nil,
				&Pulse{Ticks: 192, Velocity: 90}, nil, nil, nil,
				&Pulse{Ticks: 288, Velocity: 90}, nil, nil, nil}},
		},
		{name: "without dots", str: "X___x   X~~~*...", want: &Pattern{
			Grid: One16,
			PPQN: 96,
			Pulses: []*Pulse{
				&Pulse{Ticks: 0, Velocity: 90}, nil, nil, nil,
				&Pulse{Ticks: 96, Velocity: 90}, nil, nil, nil,
				&Pulse{Ticks: 192, Velocity: 90}, nil, nil, nil,
				nil, nil, nil, nil}},
		},
		{name: "blank", str: "blank", want: &Pattern{Grid: One16, PPQN: 96, Pulses: []*Pulse{nil, nil, nil, nil, nil}}},
		{name: "empty", str: "", want: &Pattern{Grid: One16, PPQN: 96}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFromString(tt.str)[0].PPQN; !reflect.DeepEqual(got, tt.want.PPQN) {
				t.Errorf("NewFromString().PPQN = %d, want %d", got, tt.want.PPQN)
			}
			if got := NewFromString(tt.str)[0].Grid; !reflect.DeepEqual(got, tt.want.Grid) {
				t.Errorf("NewFromString().grid = %s, want %s", got, tt.want.Grid)
			}
			if got := NewFromString(tt.str)[0].Pulses.String(); !reflect.DeepEqual(got, tt.want.Pulses.String()) {
				t.Errorf("NewFromString() = %s, want %s", got, tt.want.Pulses)
			}
		})
	}
}
