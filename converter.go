package drumbeat

import (
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
	return nil
	/*
		e := midi.NewEncoder(w, 0, 96)

		nbrSteps := len(patterns[0].Pulses)
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
				notePitch := t.Key
				var stepVal float64
				// guard
				if len(t.Pulses) > i {
					stepVal = t.Pulses[i]
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
	*/
}

// FromMIDI converts the content of a MIDI file into drum beat patterns. Note
// that this is for drum patterns only, expect the unexpected if you use non
// drum sequences.
func FromMIDI(r io.Reader) ([]*Pattern, error) {
	dec := midi.NewDecoder(r)
	if err := dec.Parse(); err != nil {
		return nil, err
	}
	ppq := uint32(dec.TicksPerQuarterNote)
	// fmt.Println("PPQN: ", dec.TicksPerQuarterNote)
	lastNoteOffOffset := map[string]uint32{}
	totalDuration := uint32(0) // in ticks
	patterns := []*Pattern{}

	absEvs := map[int][]absEv{}
	curEvsStart := map[string]int{}

	// We expect to only have 1 track with the patterns being transcribed across
	// notes where a note is a specific drum sample/instrument.
	for _, t := range dec.Tracks {
		// fmt.Println(t.Name())
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
						start:    ev.AbsTicks,
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
				absEvs[pitch] = append(absEvs[pitch], absEv{start: ev.AbsTicks, duration: totalDuration - uint32(ev.AbsTicks)})
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

	// we only have 3 grids: 1/8, 1/16, 1/32
	gridRes = adjustGrid(gridRes, uint32(dec.TicksPerQuarterNote))

	for pitch, events := range absEvs {
		// fmt.Printf("%s %#v\n", midi.NoteToName(pitch), events)
		if len(events) < 1 {
			continue
		}

		pat := &Pattern{Name: midi.NoteToName(pitch), Key: pitch, PPQN: dec.TicksPerQuarterNote}

		nbrSteps := math.Ceil(float64(totalDuration) / float64(gridRes))
		pat.Pulses = make(Pulses, int(nbrSteps))
		pat.Velocity = make([]float64, int(nbrSteps))

		// TickPosition(val, ppq)
		gridRes := uint64(ppq / 2)
		// fmt.Println(nbrSteps, gridRes, dec.TicksPerQuarterNote)

		// TODO(mattetti): set velocity
		// TODO: quantize
		for i := range pat.Pulses {

			start := uint64(i) * gridRes
			var busy bool
			for _, e := range events {
				if e.start >= start && e.start < start+gridRes {
					if busy {
						// fmt.Println("already busy")
					}
					pat.Pulses[i] = &Pulse{
						Ticks:    start,
						Duration: uint16(gridRes),
						Velocity: e.vel,
					}
					busy = true
					// break
				}
			}
		}

		patterns = append(patterns, pat)

	}

	return patterns, nil
}

func adjustGrid(res, quarter uint32) uint32 {

	// 1/8th
	if res > quarter {
		return quarter
	}
	// 1/16th
	if res >= quarter/2 {
		return quarter / 2
	}
	// 1/32th
	if res >= quarter/4 {
		return quarter / 4
	}
	// 1/64th
	return quarter / 8
	// 	fmt.Println("more than 1/16", gridRes)
	// 	gridRes = uint32(dec.TicksPerQuarterNote) / 4
	// }
	// 	} else {
	// 		// only 1/32th if we have a note that is at that resolution or lower
	// 		if gridRes <= uint32(dec.TicksPerQuarterNote)/8 {
	// 			// gridRes = uint32(dec.TicksPerQuarterNote) / 8
	// 		} else {
	// 			gridRes = uint32(dec.TicksPerQuarterNote) / 4
	// 		}
	// 	}
}
