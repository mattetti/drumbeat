package drumbeat

import (
	"bytes"
	"strings"
)

// NewFromString converts a string where `x` are converted into active pulses.
// The default step duration is 1/8th and the pulse is on for the entire step.
// Default velocity is 0.9
func NewFromString(str string) []*Pattern {
	// TODO(mattetti): support multiplexing patterns when separated by a `;`
	ppqn := uint64(DefaultPPQN)
	pat := &Pattern{PPQN: DefaultPPQN, Grid: One16}
	gridRes := ppqn / 2

	pat.Pulses = make(Pulses, len(str))
	for i, r := range strings.ToLower(str) {
		if r == 'x' {
			pat.Pulses[i] = &Pulse{
				Ticks:    gridRes * uint64(i),
				Velocity: 90,
				Duration: uint16(gridRes),
			}
		}
	}

	return []*Pattern{pat}
}

// Pattern represent the content of a drum pattern/beat.
type Pattern struct {
	// Name of the pattern or instrument
	Name string
	// Steps are the values for each step 0.0 means no pulse, a pulse greater
	// than 0 indicates the duration in beats of the pulse
	Pulses Pulses
	// Key indicates the MIDI key this pattern should be triggering. Useful when
	// converting to MIDI
	Key int
	// PPQN is the amount of ticks per quarter note.
	PPQN uint16
	// Grid is the resolution of the pattern
	Grid GridRes
}

// Offset offsets the slice of pulses by moving the pulses to the right by n
// positions.
func (p *Pattern) Offset(n int) {
	if p == nil {
		return
	}
	total := len(p.Pulses)
	cutoff := total - (n % total)
	if cutoff == total {
		return
	}
	if cutoff > total {
		cutoff -= total
	}
	// cutoff is where we are starting now
	// we need to remove cutoff * step in ticks to the entries after the cutoff
	// and we need to add `cutoff * step in ticks` to the other steps

	// TODO: offset the Ticks positions
	p.Pulses = append(p.Pulses[cutoff:], p.Pulses[:cutoff]...)
}

// Pulses is a collection of ordered pulses
type Pulses []*Pulse

// Pulse indicates a drum hit
type Pulse struct {
	Ticks    uint64
	Duration uint16
	Velocity uint8
}

// String implements the stringer interface
func (pulses Pulses) String() string {
	buf := bytes.Buffer{}
	for _, s := range pulses {
		if s != nil && s.Velocity > 0 {
			// TODO: vel > 100 == A
			buf.WriteString(`x`)
		} else {
			buf.WriteString(`.`)
		}
	}
	return buf.String()
}
