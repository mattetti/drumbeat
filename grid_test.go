package drumbeat

import "testing"

func TestPattern_StepSize(t *testing.T) {
	tests := []struct {
		name string
		grid GridRes
		want uint64
	}{
		{name: "1/4th", grid: One4, want: uint64(DefaultPPQN)},
		{name: "1/8th", grid: One8, want: uint64(DefaultPPQN) / 2},
		{name: "1/16th", grid: One16, want: uint64(DefaultPPQN) / 4},
		{name: "1/32th", grid: One32, want: uint64(DefaultPPQN) / 8},
		{name: "1/64th", grid: One64, want: uint64(DefaultPPQN) / 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pattern{
				Name: tt.name,
				PPQN: DefaultPPQN,
				Grid: tt.grid,
			}
			if got := p.StepSize(); got != tt.want {
				t.Errorf("Pattern.StepSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
