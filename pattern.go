package drumbeat

import (
	"math"
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

// ReAlign adds the nil steps if the pulses are unbalanced and reorder the steps
// if needed. This also makes sure we have the right number of pulses to fill
// full bars.
func (p *Pattern) ReAlign() {
	if p == nil {
		return
	}
	var max uint64
	for _, pulse := range p.Pulses {
		if pulse == nil {
			continue
		}
		if pulse.Ticks > max {
			max = pulse.Ticks
		}
	}
	stepSize := p.StepSize()
	gridSteps := (max + stepSize) / stepSize

	// set a minimum of steps to have
	minSteps := int(p.Grid.StepsInBeat() * 4)
	if max > (uint64(minSteps) * stepSize) {
		// we have a step that is further than the last step at minimum length
		steps := int(math.Ceil(float64(max+1) / float64(stepSize)))
		for (steps % minSteps) != 0 {
			steps++
		}
		minSteps = steps
	}

	// make sure we fullfill the minimum quota of steps
	if len(p.Pulses) < minSteps {
		gridSteps = uint64(minSteps)
		p.Pulses = append(p.Pulses, make([]*Pulse, int(gridSteps)-len(p.Pulses))...)
	}

	// make sure we fill full bars
	for i := 0; gridSteps < uint64(len(p.Pulses)) || (int(gridSteps)%minSteps != 0); i++ {
		if i == 0 {
			gridSteps = uint64(minSteps)
			continue
		}
		gridSteps += uint64(minSteps)
	}

	newPulses := make([]*Pulse, gridSteps)
	for i := uint64(0); i < gridSteps; i++ {
		start := i * stepSize
		end := start + stepSize
		for _, pulse := range p.Pulses {
			if pulse == nil {
				continue
			}
			if pulse.Ticks >= start && pulse.Ticks < end {
				// we only keep 1 pulse per step, the earliest
				if exPulse := newPulses[i]; exPulse != nil {
					if exPulse.Ticks > pulse.Ticks {
						newPulses[i] = pulse
					}
					continue
				}
				newPulses[i] = pulse
			}
		}
	}

	p.Pulses = newPulses
}

// Offset offsets the slice of pulses by moving the pulses to the right by n
// positions.
func (p *Pattern) Offset(n int) {
	if p == nil {
		return
	}
	total := len(p.Pulses)
	for n > total {
		n -= total
	}
	cutoffIDX := total - (n % total)
	if cutoffIDX == total {
		return
	}
	if cutoffIDX > total {
		cutoffIDX -= total
	}
	stepSize := p.StepSize()
	for i, pulse := range p.Pulses {
		if pulse == nil {
			continue
		}
		stepTick := (uint64(i) * stepSize)
		if stepTick > pulse.Ticks {
			pulse.Ticks = 0
			continue
		}
		pulse.Ticks -= stepTick
	}
	p.Pulses = append(p.Pulses[cutoffIDX:], p.Pulses[:cutoffIDX]...)
	for i, pulse := range p.Pulses {
		if pulse == nil {
			continue
		}
		pulse.Ticks += (uint64(i) * stepSize)
	}
}
