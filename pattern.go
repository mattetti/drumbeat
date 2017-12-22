package drumbeat

import (
	"bytes"
	"strings"

	"github.com/mattetti/audio/midi"
)

// NewFromString converts a string where `x` are converted into active pulses.
// The default step duration is 1/8th and the pulse is on for the entire step.
// Default velocity is 0.9
func NewFromString(str string) []*Pattern {
	// TODO(mattetti): support multiplexing patterns when separated by a `;`
	pat := &Pattern{StepDuration: midi.Dur8th}
	pat.Steps = make(Pulses, len(str))
	pat.Velocity = make([]float64, len(str))
	for i, r := range strings.ToLower(str) {
		if r == 'x' {
			pat.Steps[i] = 1.0
			pat.Velocity[i] = 0.9
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
	Steps Pulses
	// Velocity indicates the velocity of each step between 0 and 1
	Velocity []float64
	// StepDuration indicates the length of a full step
	StepDuration float64
	// Key indicates the MIDI key this pattern should be triggering. Useful when
	// converting to MIDI
	Key int
}

// Pulses is a collection of pulses and their durations
// Each pulse starts on the corresponding step.
type Pulses []float64

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
	pulses = append(pulses[cutoff:], pulses[:cutoff]...)
	return pulses
}

// String implements the stringer interface
func (pulses Pulses) String() string {
	buf := bytes.Buffer{}
	for _, s := range pulses {
		if s > 0.0 {
			buf.WriteString(`x`)
		} else {
			buf.WriteString(`.`)
		}
	}
	return buf.String()
}
