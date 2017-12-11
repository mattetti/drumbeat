package drumbeat

import (
	"bytes"
)

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
			buf.WriteString(`X`)
		} else {
			buf.WriteString(`.`)
		}
	}
	return buf.String()
}
