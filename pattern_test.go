package drumbeat

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/go-audio/midi"
)

func TestPulses_Offset(t *testing.T) {
	tests := []struct {
		name   string
		pulses Pulses
		n      int
		want   Pulses
	}{
		{name: "shift 2 to the right",
			pulses: []*Pulse{{Ticks: 0, Velocity: 90}, nil, nil, nil},
			n:      2,
			want:   []*Pulse{nil, nil, {Ticks: 48, Velocity: 90}, nil}},
		{name: "shift 2 to the right once again",
			pulses: []*Pulse{nil, {Ticks: 24, Velocity: 90}, {Ticks: 48, Velocity: 90}, nil},
			n:      2,
			want:   []*Pulse{{Ticks: 0, Velocity: 90}, nil, nil, {Ticks: 72, Velocity: 90}}},
		{name: "shift by the length of the slice",
			pulses: []*Pulse{{Ticks: 0, Velocity: 90}, {Ticks: 24, Velocity: 90}},
			n:      2,
			want:   []*Pulse{{Ticks: 0, Velocity: 90}, {Ticks: 24, Velocity: 90}}},
		{name: "shift by more than the length of the slice",
			pulses: []*Pulse{{Ticks: 0, Velocity: 90}, nil, {Ticks: 24, Velocity: 90}},
			n:      4,
			want:   []*Pulse{{Ticks: 0, Velocity: 90}, {Ticks: 24, Velocity: 90}, nil}},
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
				for i, p := range pat.Pulses {
					if !reflect.DeepEqual(p, tt.want[i]) {
						t.Logf("[%d] got: %+v vs want: %+v\n", i, p, tt.want[i])
					}
				}
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
			{Ticks: 0, Velocity: 99},
			nil,
			{Ticks: 0, Velocity: 99},
		}, ".X.X"},
		{[]*Pulse{nil, &Pulse{}, nil, nil, nil, nil, nil, nil}, "........"},
		{[]*Pulse{
			{Ticks: 0, Velocity: 90}, {Ticks: 24, Velocity: 90}, {Ticks: 48, Velocity: 90}, {Ticks: 72, Velocity: 90},
			{Ticks: 96, Velocity: 90}, {Ticks: 120, Velocity: 90}, {Ticks: 144, Velocity: 90}, {Ticks: 168, Velocity: 90},
			{Ticks: 192, Velocity: 90}, {Ticks: 216, Velocity: 90}, {Ticks: 240, Velocity: 90}, {Ticks: 264, Velocity: 90},
			{Ticks: 288, Velocity: 90}, {Ticks: 312, Velocity: 90}, {Ticks: 336, Velocity: 90}, {Ticks: 360, Velocity: 90}},
			"XXXXXXXXXXXXXXXX"},
		{[]*Pulse{{Ticks: 0, Velocity: 90}, {Ticks: 24, Velocity: 90}, {Ticks: 48, Velocity: 90}, {Ticks: 72, Velocity: 90}}, "XXXX"},
		{[]*Pulse{nil, {Ticks: 24, Velocity: 90}, {Ticks: 48, Velocity: 90}, nil}, ".XX."},
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
		want []*Pattern
	}{
		{name: "basic", str: "x...x...x...x...", want: []*Pattern{{
			Grid: One16,
			PPQN: 96,
			Pulses: []*Pulse{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: 96, Velocity: 90}, nil, nil, nil,
				{Ticks: 192, Velocity: 90}, nil, nil, nil,
				{Ticks: 288, Velocity: 90}, nil, nil, nil}}},
		},
		{name: "basic with return carriage and tabs", str: `x...x...
			x...x...`, want: []*Pattern{{
			Grid: One16,
			PPQN: 96,
			Pulses: []*Pulse{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: 96, Velocity: 90}, nil, nil, nil,
				{Ticks: 192, Velocity: 90}, nil, nil, nil,
				{Ticks: 288, Velocity: 90}, nil, nil, nil}}},
		},
		{name: "basic with name at the beginning", str: "[Kick]x...x...x...x...", want: []*Pattern{{
			Grid: One16,
			PPQN: 96,
			Name: "Kick",
			Pulses: []*Pulse{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: 96, Velocity: 90}, nil, nil, nil,
				{Ticks: 192, Velocity: 90}, nil, nil, nil,
				{Ticks: 288, Velocity: 90}, nil, nil, nil}}},
		},
		{name: "basic with name at the end", str: "x...x...x...x...[Kick]", want: []*Pattern{{
			Grid: One16,
			PPQN: 96,
			Name: "Kick",
			Pulses: []*Pulse{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: 96, Velocity: 90}, nil, nil, nil,
				{Ticks: 192, Velocity: 90}, nil, nil, nil,
				{Ticks: 288, Velocity: 90}, nil, nil, nil}}},
		},
		{name: "basic with name in the middle", str: "x...x...	[Kick]	x...x...", want: []*Pattern{{
			Grid: One16,
			PPQN: 96,
			Name: "Kick",
			Pulses: []*Pulse{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: 96, Velocity: 90}, nil, nil, nil,
				{Ticks: 192, Velocity: 90}, nil, nil, nil,
				{Ticks: 288, Velocity: 90}, nil, nil, nil}}},
		},
		{name: "basic with a bracket but no names", str: "x[..x...	x...x...", want: []*Pattern{{
			Grid: One16,
			PPQN: 96,
			Pulses: []*Pulse{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: 96, Velocity: 90}, nil, nil, nil,
				{Ticks: 192, Velocity: 90}, nil, nil, nil,
				{Ticks: 288, Velocity: 90}, nil, nil, nil}}},
		},
		{name: "basic with key at the beginning", str: "{D2}x...x...x...x...", want: []*Pattern{{
			Grid: One16,
			PPQN: 96,
			Key:  midi.KeyInt("D", 2),
			Pulses: []*Pulse{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: 96, Velocity: 90}, nil, nil, nil,
				{Ticks: 192, Velocity: 90}, nil, nil, nil,
				{Ticks: 288, Velocity: 90}, nil, nil, nil}}},
		},
		{name: "basic with key at the end", str: "x...x...x...x...{C#3}", want: []*Pattern{{
			Grid: One16,
			PPQN: 96,
			Key:  midi.KeyInt("C#", 3),
			Pulses: []*Pulse{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: 96, Velocity: 90}, nil, nil, nil,
				{Ticks: 192, Velocity: 90}, nil, nil, nil,
				{Ticks: 288, Velocity: 90}, nil, nil, nil}}},
		},
		{name: "basic with bad key data", str: "{noteAKey}x...x...x...x...", want: []*Pattern{{
			Grid: One16,
			PPQN: 96,
			Pulses: []*Pulse{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: 96, Velocity: 90}, nil, nil, nil,
				{Ticks: 192, Velocity: 90}, nil, nil, nil,
				{Ticks: 288, Velocity: 90}, nil, nil, nil}}},
		},
		{name: "with uppercase X", str: "X...x...X...X...", want: []*Pattern{{
			Grid: One16,
			PPQN: 96,
			Pulses: []*Pulse{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: 96, Velocity: 90}, nil, nil, nil,
				{Ticks: 192, Velocity: 90}, nil, nil, nil,
				{Ticks: 288, Velocity: 90}, nil, nil, nil}}},
		},
		{name: "without dots", str: "X___x   X~~~*...", want: []*Pattern{{
			Grid: One16,
			PPQN: 96,
			Pulses: []*Pulse{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: 96, Velocity: 90}, nil, nil, nil,
				{Ticks: 192, Velocity: 90}, nil, nil, nil,
				nil, nil, nil, nil}}},
		},
		{name: "blank", str: "blank", want: []*Pattern{{Grid: One8, PPQN: 96, Pulses: []*Pulse{nil, nil, nil, nil, nil}}}},
		{name: "empty", str: "", want: []*Pattern{{Grid: One8, PPQN: 96}}},
		{name: "basic 2 patterns", str: "x.......x.......;....x.......x...", want: []*Pattern{
			{Grid: One16,
				PPQN: 96,
				Pulses: []*Pulse{
					{Ticks: 0, Velocity: 90}, nil, nil, nil,
					nil, nil, nil, nil,
					{Ticks: 192, Velocity: 90}, nil, nil, nil,
					nil, nil, nil, nil,
				},
			},
			{Grid: One16,
				PPQN: 96,
				Pulses: []*Pulse{
					nil, nil, nil, nil,
					{Ticks: 96, Velocity: 90}, nil, nil, nil,
					nil, nil, nil, nil,
					{Ticks: 288, Velocity: 90}, nil, nil, nil,
				},
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFromString(tt.want[0].Grid, tt.str)
			for i, want := range tt.want {
				if !reflect.DeepEqual(got[i].PPQN, want.PPQN) {
					t.Errorf("[%d] NewFromString().PPQN = %d, want %d", i, got[i].PPQN, want.PPQN)
				}
				if !reflect.DeepEqual(got[i].Grid, want.Grid) {
					t.Errorf("[%d] NewFromString().grid = %s, want %s", i, got[i].Grid, want.Grid)
				}
				if !reflect.DeepEqual(got[i].Name, want.Name) {
					t.Errorf("[%d] NewFromString().Name = %s, want %s", i, got[i].Name, want.Name)
				}
				if !reflect.DeepEqual(got[i].Key, want.Key) {
					t.Errorf("[%d] NewFromString().Key = %d, want %d", i, got[i].Key, want.Key)
				}
				if !reflect.DeepEqual(got[i].Pulses.String(), want.Pulses.String()) {
					t.Errorf("[%d] NewFromString() = %s, want %s", i, got[i].Pulses.String(), want.Pulses)
				}
			}
		})
	}
	// making sure a nil pattern getting offset doesn't crash
	var nilPat *Pattern
	nilPat.Offset(42)
}

func TestPattern_ReAlign(t *testing.T) {
	tests := []struct {
		name   string
		pulses Pulses
		PPQN   uint16
		grid   GridRes
		want   Pulses
	}{
		{name: "blank",
			pulses: Pulses{},
			PPQN:   DefaultPPQN,
			grid:   One8,
			want:   Pulses{nil, nil, nil, nil, nil, nil, nil, nil}},
		{name: "1 step",
			pulses: Pulses{{Ticks: 0, Velocity: 90}},
			PPQN:   DefaultPPQN,
			grid:   One8,
			want:   Pulses{{Ticks: 0, Velocity: 90}, nil, nil, nil, nil, nil, nil, nil}},
		{name: "1st step in 2nd div",
			pulses: Pulses{{Ticks: 48, Velocity: 90}},
			PPQN:   DefaultPPQN,
			grid:   One8,
			want:   Pulses{nil, {Ticks: 48, Velocity: 90}, nil, nil, nil, nil, nil, nil}},
		{name: "1st step not quantized",
			pulses: Pulses{{Ticks: 24, Velocity: 90}},
			PPQN:   DefaultPPQN,
			grid:   One8,
			want:   Pulses{{Ticks: 24, Velocity: 90}, nil, nil, nil, nil, nil, nil, nil}},
		{name: "1st step not quantized - missing nil steps",
			pulses: Pulses{{Ticks: 24, Velocity: 90}, nil, nil, nil},
			PPQN:   DefaultPPQN,
			grid:   One8,
			want:   Pulses{{Ticks: 24, Velocity: 90}, nil, nil, nil, nil, nil, nil, nil}},
		{name: "2 pulses in 1 step",
			pulses: Pulses{{Ticks: 0, Velocity: 90}, {Ticks: 24, Velocity: 90}},
			PPQN:   DefaultPPQN,
			grid:   One8,
			want:   Pulses{{Ticks: 0, Velocity: 90}, nil, nil, nil, nil, nil, nil, nil}},
		{name: "2 unordered pulses in 1 step",
			pulses: Pulses{{Ticks: 24, Velocity: 90}, {Ticks: 0, Velocity: 90}},
			PPQN:   DefaultPPQN,
			grid:   One8,
			want:   Pulses{{Ticks: 0, Velocity: 90}, nil, nil, nil, nil, nil, nil, nil}},
		{name: "2 steps - 2 bars",
			pulses: Pulses{{Ticks: 0, Velocity: 90}, {Ticks: uint64(DefaultPPQN), Velocity: 90}, nil, nil, nil, nil, nil, nil, nil},
			PPQN:   DefaultPPQN,
			grid:   One16,
			want: Pulses{
				{Ticks: 0, Velocity: 90}, nil, nil, nil,
				{Ticks: uint64(DefaultPPQN), Velocity: 90}, nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil}},
		{name: "2 steps - fill 2 bars",
			pulses: Pulses{
				{Ticks: 0, Velocity: 90}, {Ticks: uint64(DefaultPPQN), Velocity: 90}, nil, nil, nil, nil, nil, nil, nil},
			PPQN: DefaultPPQN,
			grid: One8,
			want: Pulses{
				{Ticks: 0, Velocity: 90}, nil,
				{Ticks: uint64(DefaultPPQN), Velocity: 90}, nil,
				nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil}},
		{name: "1 step - fill 1 1/32 bar",
			pulses: Pulses{
				{Ticks: 0, Velocity: 90}},
			PPQN: DefaultPPQN,
			grid: One32,
			want: Pulses{
				{Ticks: 0, Velocity: 90}, nil, nil, nil, nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil}},
		{name: "1 step - fill 1 1/64 bar",
			pulses: Pulses{
				{Ticks: uint64(DefaultPPQN), Velocity: 90}},
			PPQN: DefaultPPQN,
			grid: One64,
			want: Pulses{
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				{Ticks: uint64(DefaultPPQN), Velocity: 90}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			}},
		{name: "spread out pulses",
			// over 4 bars
			pulses: Pulses{{Ticks: 0, Velocity: 90}, {Ticks: 384, Velocity: 90}, {Ticks: 768, Velocity: 90}, {Ticks: 1152, Velocity: 90}},
			PPQN:   DefaultPPQN,
			grid:   One16,
			want: Pulses{
				{Ticks: 0, Velocity: 90}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				{Ticks: 384, Velocity: 90}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				{Ticks: 768, Velocity: 90}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				{Ticks: 1152, Velocity: 90}, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			}},
	}
	for _, tt := range tests {
		// Make sure we don't crash on a nil pattern
		var p *Pattern
		p.ReAlign()

		t.Run(tt.name, func(t *testing.T) {
			p := &Pattern{
				Name:   tt.name,
				Pulses: tt.pulses,
				PPQN:   tt.PPQN,
				Grid:   tt.grid,
			}
			p.ReAlign()
			if len(p.Pulses) != len(tt.want) {
				t.Fatalf("expected %d pulses, got %d", len(tt.want), len(p.Pulses))
			}
			if !reflect.DeepEqual(tt.want, p.Pulses) {
				for i, pulse := range p.Pulses {
					if pulse != tt.want[i] {
						t.Logf("[%d] wanted %+v got %+v\n", i, tt.want[i], pulse)
					}
				}
				t.Errorf("expected %#v, got %#v", tt.pulses, p.Pulses)
			}
		})
	}
}
