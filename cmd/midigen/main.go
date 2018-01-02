package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-audio/generator/euclidean"
	"github.com/go-audio/midi"

	"github.com/mattetti/drumbeat"
)

func main() {
	patternStr := flag.String("pattern", "", "The pattern to convert to MIDI")
	// 2 beats at 1/16th.
	genSteps := flag.Int("steps", 32, "Number of steps to use for generation.")
	genPulses := flag.Int("pulses", 0, "Number of pulses for the steps")
	genOffset := flag.Int("offset", 0, "Offset for the first pulse")

	flag.Parse()
	if *genSteps < 8 {
		*genSteps = 8
	}
	if *genPulses < 1 {
		*genPulses = 1 + *genSteps/8
	}

	if *patternStr != "" {
		patterns := drumbeat.NewFromString(drumbeat.One16, *patternStr)
		f, err := os.Create("drumbeat.mid")
		if err != nil {
			log.Println("something wrong happened when creating the MIDI file", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := drumbeat.ToMIDI(f, patterns...); err != nil {
			log.Fatal(err)
		}
		fmt.Println("drumbeat.mid generated off of", *patternStr)
		return
	}

	// TODO: support controls for each channel, not just the kick
	// TODO: velocity distribution

	// split in two to give it more swag
	kickSeq := euclidean.Rhythm(*genPulses/2, *genSteps/2)
	// The second time, we through in an extra kick, for free
	kickSeq = append(kickSeq, euclidean.Rhythm((*genPulses/2)+1, *genSteps/2)...)
	kickBeat := drumbeat.NewFromString(boolsToSeq(kickSeq))[0]
	if *genOffset != 0 {
		kickBeat.Offset(*genOffset)
	}
	kickBeat.Key = midi.KeyInt("C", 1)
	kickBeat.Name = "Kick"

	snareSeq := euclidean.Rhythm((*genPulses/2)+1, *genSteps)
	snareBeat := drumbeat.NewFromString(boolsToSeq(snareSeq))[0]
	snareBeat.Offset(4)
	snareBeat.Key = midi.KeyInt("D", 1)
	snareBeat.Name = "Snare"

	// add some randomness to those hats
	total := *genSteps
	chunkSize := 3
	groupSize := total / chunkSize
	hatSeq := []bool{}
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < chunkSize; i++ {
		pulses := (*genSteps / chunkSize) / 2
		x := rnd.Intn(100)
		if x%2 == 0 {
			pulses++
		}
		hatSeq = append(hatSeq, euclidean.Rhythm(pulses, groupSize)...)
	}
	leftOver := total % 3
	if leftOver > 0 {
		hatSeq = append(hatSeq, hatSeq[len(hatSeq)-leftOver:]...)
	}
	hatBeat := drumbeat.NewFromString(boolsToSeq(hatSeq))[0]
	hatBeat.Key = midi.KeyInt("F#", 1)
	hatBeat.Name = "HiHat"

	fmt.Println(kickBeat.Name, "\t", kickBeat.Pulses)
	fmt.Println(snareBeat.Name, "\t", snareBeat.Pulses)
	fmt.Println(hatBeat.Name, "\t", hatBeat.Pulses)
	f, err := os.Create("gen_drumbeat.mid")
	if err != nil {
		log.Println("something wrong happened when creating the MIDI file", err)
		os.Exit(1)
	}
	defer f.Close()
	drumbeat.ToMIDI(f, kickBeat, snareBeat, hatBeat)
	fmt.Println("generated MIDI pattern available at gen_drumbeat.mid")
}

func boolsToSeq(bools []bool) string {
	str := make([]byte, len(bools))
	for i, b := range bools {
		if b {
			str[i] = 'x'
		} else {
			str[i] = '.'
		}
	}
	return string(str)
}
