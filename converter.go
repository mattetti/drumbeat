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
	// Mix and matching step duration is currently broken

	trackState := map[int]bool{}

	tr := e.NewTrack()
	delta := 0.0
	pushedDelta := false
	var currentStepDuration float64
	// loop through all the steps, one step at a time and inject
	// all track states inside the same channel.
	for i := 0; i < nbrSteps; i++ {
		if i > 0 {
			if pushedDelta {
				delta = currentStepDuration
			} else {
				delta += currentStepDuration
			}
		}
		pushedDelta = false
		for _, t := range patterns {
			currentStepDuration = t.StepDuration
			noteVal := t.Key
			var stepVal float64
			// guard
			if len(t.Steps) > i {
				stepVal = t.Steps[i]
			}
			// stop previously played noted
			if stepVal == 0.0 {
				on, ok := trackState[noteVal]
				if ok && on {
					tr.Add(delta, midi.NoteOff(0, noteVal))
					trackState[noteVal] = false
					delta = 0.0
					pushedDelta = true
				}
				continue
			}
			vel := 90
			// stop notes that are already playing
			if on, ok := trackState[noteVal]; ok && on {
				tr.Add(delta, midi.NoteOff(0, noteVal))
				delta = 0.0
			}
			tr.Add(delta, midi.NoteOn(0, noteVal, vel))
			trackState[noteVal] = true
			delta = 0.0
			pushedDelta = true
		}
	}
	lastStepSet := false
	delta = currentStepDuration
	for n, on := range trackState {
		if on {
			tr.Add(delta, midi.NoteOff(0, n))
			if !lastStepSet {
				lastStepSet = true
			}
		}
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
	var lastNoteOffOffset uint32
	runningTime := uint32(0) // in ticks
	notePatterns := map[string][]*midi.Event{}
	patterns := []*Pattern{}

	// We expect to only have 1 track with the patterns being transcribed across
	// notes where a note is a specific drum sample/instrument.
	for _, t := range dec.Tracks {
		for _, ev := range t.Events {
			runningTime += ev.TimeDelta
			// fmt.Println(midi.EventMap[ev.MsgType], float64(runningTime), "beats")

			switch ev.MsgType {
			// TODO: check for a time signature
			// case midi.EventByteMap["Meta"]:
			// 	if midi.MetaCmdMap[ev.Cmd] == "Time Signature" {
			// 		// latest Time signature
			// 		timeSignature = ev.TimeSignature
			// 	}
			case midi.EventByteMap["NoteOn"]:
				n := midi.NoteToName(int(ev.Note))
				if _, ok := notePatterns[n]; !ok {
					notePatterns[n] = []*midi.Event{}
				}
				notePatterns[n] = append(notePatterns[n], ev)
				lastNoteOffOffset = 0
			case midi.EventByteMap["NoteOff"]:
				if lastNoteOffOffset != 0 {
					// we have many notes off following each other
					lastNoteOffOffset += ev.TimeDelta
				} else {
					lastNoteOffOffset = ev.TimeDelta
				}
				n := midi.NoteToName(int(ev.Note))
				if _, ok := notePatterns[n]; !ok {
					notePatterns[n] = []*midi.Event{}
				}
				notePatterns[n] = append(notePatterns[n], ev)
			}
		}
	}

	// look at the note patterns as drum patterns and extract their sequencing
	for note, events := range notePatterns {
		if len(events) < 1 || events[0] == nil {
			continue
		}
		pat := &Pattern{Name: note, Key: int(events[0].Note)}
		// TODO: convert events into steps.
		// check the shortest note which is the delta between on and off
		// we can simply look at the off events to see the duration of the note
		gridRes := uint32(dec.TicksPerQuarterNote)
		for _, ev := range events {
			if ev.MsgType == midi.EventByteMap["NoteOff"] {
				if ev.TimeDelta > 0 && ev.TimeDelta < gridRes {
					gridRes = ev.TimeDelta
				}
			}
		}
		if gridRes < uint32(dec.TicksPerQuarterNote) {
			// limit to 1/32 grid
			if min := (uint32(dec.TicksPerQuarterNote) / 8); min > gridRes {
				gridRes = min
			}
			pat.StepDuration = float64(gridRes) / float64(dec.TicksPerQuarterNote)
		} else {
			pat.StepDuration = float64(dec.TicksPerQuarterNote) / float64(gridRes)
		}

		type absEv struct {
			start uint32
			end   uint32
		}

		absEvs := []*absEv{}
		var curEvStart uint32
		var runningTime uint32
		for _, ev := range events {
			runningTime += ev.TimeDelta
			if ev.MsgType == midi.EventByteMap["NoteOn"] {
				curEvStart = runningTime
			}
			if ev.MsgType == midi.EventByteMap["NoteOff"] {
				absEvs = append(absEvs,
					&absEv{start: curEvStart, end: runningTime})
				curEvStart = 0
			}
		}
		// fmt.Println("grid length note", gridRes, "; duration", runningTime)
		nbrSteps := math.Ceil(float64(runningTime) / float64(gridRes))
		pat.Steps = make(Pulses, int(nbrSteps))

		// TODO(mattetti): set velocity
		for i, _ := range pat.Steps {
			start := uint32(i) * gridRes
			end := (uint32(i) + 1) * gridRes
			for _, e := range absEvs {
				if e.start >= start && e.end <= end {
					pat.Steps[i] = pat.StepDuration
					break
				}
			}
		}

		for _, e := range absEvs {
			fmt.Printf("From %d to %d\n", e.start, e.end)
		}

		// fmt.Println(events)
		patterns = append(patterns, pat)
	}

	return patterns, nil
}
