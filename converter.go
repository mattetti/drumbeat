package drumbeat

import (
	"fmt"
	"io"
	"math"

	"github.com/go-audio/midi"
)

// absolute representation of a pulse the duration of the event indicates
// how long it lasts but the pulse will be represented as a single hit at
// the grid resolution
type absEv struct {
	start    uint64
	duration uint32
	vel      uint8
}

// ToMIDI converts the passed patterns to a single MIDI file.
func ToMIDI(w io.WriteSeeker, patterns ...*Pattern) error {
	if len(patterns) < 1 || patterns[0] == nil {
		return nil
	}

	nbrSteps := len(patterns[0].Pulses)
	ppq := patterns[0].PPQN
	e := midi.NewEncoder(w, 0, ppq)

	trackState := map[int]bool{}

	tr := e.NewTrack()
	var delta uint32
	// var pushedDelta bool

	currentStepDuration := uint32(ppq) / 2
	// loop through all the steps, one step at a time and inject
	// all track states inside the same channel.
	for i := 0; i < nbrSteps; i++ {
		if i > 0 {
			delta += currentStepDuration
		}
		for _, t := range patterns {
			notePitch := t.Key
			var stepVal *Pulse
			// guard
			if len(t.Pulses) > i {
				stepVal = t.Pulses[i]
			}

			// empty step: stop playing note if needed
			if stepVal == nil {
				if on, ok := trackState[notePitch]; ok && on {
					// note is playing,let's stop it
					tr.AddAfterDelta(delta, midi.NoteOff(0, notePitch))
					trackState[notePitch] = false
					delta = 0.0
				}
				continue
			}

			// we have a pulse!

			// TODO: maybe offer a difference between x an X for velocity
			vel := 90
			tr.AddAfterDelta(delta, midi.NoteOn(0, notePitch, vel))
			// mark note as playing
			trackState[notePitch] = true
			delta = 0.0
		}
	}

	var wroteLastStep bool
	for pitch, on := range trackState {
		if on {
			tr.AddAfterDelta(delta+currentStepDuration, midi.NoteOff(0, pitch))
			delta = 0
			wroteLastStep = true
		}
	}

	// end the track after the last step
	if wroteLastStep {
		tr.Add(0, midi.EndOfTrack())
	} else {
		tr.AddAfterDelta(delta+currentStepDuration, midi.EndOfTrack())
	}

	return e.Write()
}

// FromMIDI converts the content of a MIDI file into drum beat patterns. Note
// that this is for drum patterns only, expect the unexpected if you use non
// drum sequences.
func FromMIDI(r io.Reader) ([]*Pattern, error) {
	dec := midi.NewDecoder(r)
	if err := dec.Parse(); err != nil {
		return nil, err
	}
	totalDuration := uint32(0) // in ticks
	patterns := []*Pattern{}

	absEvs := map[int][]absEv{}
	curEvsStart := map[string]int{}

	// We expect to only have 1 track with the patterns being transcribed across
	// notes where a note is a specific drum sample/instrument.
	for _, t := range dec.Tracks {
		for _, ev := range t.Events {
			totalDuration += ev.TimeDelta
			pitch := int(ev.Note)
			n := midi.NoteToName(pitch)
			// fmt.Printf("%s %s @ %.2f beats\n", n, midi.EventMap[ev.MsgType], float64(totalDuration))

			if _, ok := absEvs[pitch]; !ok {
				absEvs[pitch] = []absEv{}
			}
			if _, ok := curEvsStart[n]; !ok {
				curEvsStart[n] = -1
			}
			switch ev.MsgType {
			// TODO: check for a time signature
			// case midi.EventByteMap["Meta"]:
			// 	if midi.MetaCmdMap[ev.Cmd] == "Time Signature" {
			// 		// latest Time signature
			// 		timeSignature = ev.TimeSignature
			// 	}
			case midi.EventByteMap["NoteOn"]:
				if curEvsStart[n] >= 0 {
					// end previous note
					start := uint32(curEvsStart[n])
					absEvs[pitch] = append(absEvs[pitch], absEv{
						start:    uint64(curEvsStart[n]),
						duration: totalDuration - start,
						vel:      ev.Velocity},
					)
				}
				curEvsStart[n] = int(ev.AbsTicks)
			case midi.EventByteMap["NoteOff"]:
				absEvs[pitch] = append(absEvs[pitch],
					absEv{
						start:    uint64(curEvsStart[n]),
						duration: totalDuration - uint32(ev.AbsTicks),
					})
				curEvsStart[n] = -1
			}
		}
	}

	// 1/16th
	gridRes := uint32(dec.TicksPerQuarterNote) / 2

	for pitch, events := range absEvs {
		if len(events) < 1 {
			continue
		}

		pat := &Pattern{
			Name: midi.NoteToName(pitch),
			Key:  pitch,
			PPQN: dec.TicksPerQuarterNote,
		}

		nbrSteps := math.Ceil(float64(totalDuration) / float64(gridRes))
		pat.Pulses = make(Pulses, int(nbrSteps))

		// TODO: quantize
		for i := range pat.Pulses {
			start := uint64(i) * uint64(gridRes)
			var busy bool
			for _, e := range events {
				if e.start >= start && e.start < start+uint64(gridRes) {
					if busy {
						fmt.Println("already busy")
					}
					pat.Pulses[i] = &Pulse{
						Ticks:    start,
						Duration: uint16(gridRes),
						Velocity: e.vel,
					}
					busy = true
				}
			}
		}

		patterns = append(patterns, pat)
	}

	return patterns, nil
}
