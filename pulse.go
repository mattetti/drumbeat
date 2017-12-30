package drumbeat

import "bytes"

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
