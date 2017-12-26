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
	ppqn := uint64(96)
	gridRes := ppqn / 2
	pat := &Pattern{PPQN: uint16(ppqn)}
	pat.Pulses = make(Pulses, len(str))
	pat.Velocity = make([]float64, len(str))
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
	// Velocity indicates the velocity of each step between 0 and 1
	Velocity []float64
	// Key indicates the MIDI key this pattern should be triggering. Useful when
	// converting to MIDI
	Key int
	// PPQN is the amount of ticks per quarter note.
	PPQN uint16
}

// Pulses is a collection of ordered pulses
type Pulses []*Pulse

// Pulse indicates a drum hit
type Pulse struct {
	Ticks    uint64
	Duration uint16
	Velocity uint8
}

// Offset offsets the slice of pulses by moving the pulses to the right by n
// positions.
func (pulses Pulses) Offset(n int) Pulses {
	total := len(pulses)
	cutoff := total - (n % total)
	if cutoff == total {
		return pulses
	}
	if cutoff > total {
		cutoff -= total
	}
	// TODO: offset the Ticks positions
	pulses = append(pulses[cutoff:], pulses[:cutoff]...)
	return pulses
}

// String implements the stringer interface
func (pulses Pulses) String() string {
	buf := bytes.Buffer{}
	for _, s := range pulses {
		if s != nil {
			// TODO: vel > 100 == A
			buf.WriteString(`x`)
		} else {
			buf.WriteString(`.`)
		}
	}
	return buf.String()
}
