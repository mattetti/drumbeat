package drumbeat

import (
	"fmt"
	"io"
	"math"

	"github.com/mattetti/audio/midi"
)

// ToMIDI converts the passed patterns to a single MIDI file.
func ToMIDI(w io.WriteSeeker, patterns ...*Pattern) error {
	if len(patterns) < 1 || patterns[0] == nil {
		return nil
	}
	e := midi.NewEncoder(w, 0, 96)

	nbrSteps := len(patterns[0].Steps)
	//TODO: Mix and matching step duration is currently broken

	trackState := map[int]bool{}

	tr := e.NewTrack()
	var delta float64
	// var pushedDelta bool
	var currentStepDuration float64
	// loop through all the steps, one step at a time and inject
	// all track states inside the same channel.
	for i := 0; i < nbrSteps; i++ {
		if i > 0 {
			delta += currentStepDuration
		}
		for _, t := range patterns {
			currentStepDuration = t.StepDuration
			notePitch := t.Key
			var stepVal float64
			// guard
			if len(t.Steps) > i {
				stepVal = t.Steps[i]
			}

			// empty step: stop playing note if needed
			if stepVal == 0.0 {
				if on, ok := trackState[notePitch]; ok && on {
					// note is playing,let's stop it
					tr.Add(delta, midi.NoteOff(0, notePitch))
					trackState[notePitch] = false
					delta = 0.0
				}
				continue
			}

			// we have a pulse!

			// TODO: maybe offer a difference between x an X for velocity
			vel := 90
			tr.Add(delta, midi.NoteOn(0, notePitch, vel))
			// mark note as playing
			trackState[notePitch] = true
			delta = 0.0
		}
	}

	var wroteLastStep bool
	for pitch, on := range trackState {
		if on {
			tr.Add(delta+currentStepDuration, midi.NoteOff(0, pitch))
			delta = 0
			wroteLastStep = true
		}
	}

	// end the track after the last step
	if wroteLastStep {
		tr.Add(0, midi.EndOfTrack())
	} else {
		tr.Add(delta+currentStepDuration, midi.EndOfTrack())
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
	// fmt.Println("PPQN: ", dec.TicksPerQuarterNote)
	lastNoteOffOffset := map[string]uint32{}
	totalDuration := uint32(0) // in ticks
	patterns := []*Pattern{}

	// absolute representation of a pulse the duration of the event indicates
	// how long it lasts but the pulse will be represented as a single hit at
	// the grid resolution
	type absEv struct {
		start    uint32
		duration uint32
		vel      uint8
	}

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
						start:    start,
						duration: totalDuration - start,
						vel:      ev.Velocity},
					)
				}
				curEvsStart[n] = int(totalDuration)
				lastNoteOffOffset[n] = 0
			case midi.EventByteMap["NoteOff"]:
				if lastNoteOffOffset[n] != 0 {
					// we have many notes off following each other
					lastNoteOffOffset[n] += ev.TimeDelta
				} else {
					lastNoteOffOffset[n] = ev.TimeDelta
				}
				start := uint32(curEvsStart[n])
				absEvs[pitch] = append(absEvs[pitch], absEv{start: start, duration: totalDuration - start})
				curEvsStart[n] = -1
			}
		}
	}

	gridRes := uint32(dec.TicksPerQuarterNote)
	for _, events := range absEvs {
		if len(events) < 1 {
			continue
		}

		for _, ev := range events {
			if ev.duration > 0 && ev.duration < gridRes {
				gridRes = ev.duration
			}
		}
	}

	if gridRes < uint32(dec.TicksPerQuarterNote) {
		// limit to 1/32 grid
		if min := (uint32(dec.TicksPerQuarterNote) / 8); min > gridRes {
			gridRes = min
		}
	}

	for pitch, events := range absEvs {
		fmt.Printf("%s %#v\n", midi.NoteToName(pitch), events)
		if len(events) < 1 {
			continue
		}
		pat := &Pattern{Name: midi.NoteToName(pitch), Key: pitch}
		if gridRes < uint32(dec.TicksPerQuarterNote) {
			pat.StepDuration = float64(gridRes) / float64(dec.TicksPerQuarterNote)
		} else {
			pat.StepDuration = float64(dec.TicksPerQuarterNote) / float64(gridRes)
		}

		nbrSteps := math.Ceil(float64(totalDuration) / float64(gridRes))
		pat.Steps = make(Pulses, int(nbrSteps))
		pat.Velocity = make([]float64, int(nbrSteps))

		// TODO(mattetti): set velocity
		for i := range pat.Steps {
			start := uint32(i) * gridRes
			for _, e := range events {
				if e.start == start {
					pat.Steps[i] = pat.StepDuration
					pat.Velocity[i] = float64(e.vel) / 127.0
					break
				}
			}
		}

		patterns = append(patterns, pat)

	}

	return patterns, nil
}
