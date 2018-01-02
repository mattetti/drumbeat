package drumbeat

import (
	"math"
	"strconv"
	"strings"

	"github.com/go-audio/midi"
)

var (
	patStrReplacer = strings.NewReplacer("\t", "", "\n", "")
)

// NewFromString converts a string where `x` are converted into active pulses.
// The first argument is the resolution of the grid so we can define how many
// steps fit in a bar. Default velocity is 0.9
//
// Multiple patterns can be provided if separated by a semi colon: `;`.
func NewFromString(grid GridRes, str string) []*Pattern {
	// support multiplexing of patterns by separating them by a `;`
	patStrs := strings.Split(str, ";")

	patterns := []*Pattern{}
	ppqn := uint64(DefaultPPQN)
	for _, patStr := range patStrs {
		pat := &Pattern{PPQN: DefaultPPQN, Grid: grid}
		gridRes := ppqn / grid.StepsInBeat()

		// Name
		nameStartIDX := strings.IndexByte(patStr, '[')
		nameEndIDX := strings.IndexByte(patStr, ']')
		if nameStartIDX != -1 && nameEndIDX != -1 {
			pat.Name = patStr[nameStartIDX+1 : nameEndIDX]
			patStr = patStr[:nameStartIDX] + patStr[nameEndIDX+1:]
		}

		// Key
		keyStartIDX := strings.IndexByte(patStr, '{')
		keyEndIDX := strings.IndexByte(patStr, '}')
		if keyStartIDX != -1 && keyEndIDX != -1 {
			keyStr := patStr[keyStartIDX+1 : keyEndIDX]
			if l := len(keyStr); l > 1 {
				oct, err := strconv.Atoi(string(keyStr[l-1]))
				if err == nil {
					keyLetter := string(keyStr[:l-1])
					pat.Key = midi.KeyInt(keyLetter, oct)
				}
			}
			patStr = patStr[:keyStartIDX] + patStr[keyEndIDX+1:]
		}

		patStr = patStrReplacer.Replace(patStr)

		pat.Pulses = make(Pulses, len(patStr))
		for i, r := range strings.ToLower(patStr) {
			if r == 'x' {
				pat.Pulses[i] = &Pulse{
					Ticks:    gridRes * uint64(i),
					Velocity: 90,
					Duration: uint16(gridRes),
				}
			}
		}
		patterns = append(patterns, pat)
	}

	return patterns
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

// compact removes the nil pulses.
func (p *Pattern) compact() {
	if p == nil {
		return
	}
	pulses := []*Pulse{}
	for _, pulse := range p.Pulses {
		if pulse == nil {
			continue
		}
		pulses = append(pulses, pulse)
	}
	p.Pulses = pulses
}
