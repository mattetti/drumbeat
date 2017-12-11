package drumbeat

import (
	"io"

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
