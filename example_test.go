package drumbeat_test

import (
	"fmt"
	"log"
	"os"

	"github.com/mattetti/drumbeat"
)

func ExampleNewFromString() {
	patternStr := "x.xxx...x...x.x.x..x.xx.x..x.xxx"
	patterns := drumbeat.NewFromString(patternStr)
	f, err := os.Create("drumbeat.mid")
	if err != nil {
		log.Println("something wrong happened when creating the MIDI file", err)
		os.Exit(1)
	}
	if err := drumbeat.ToMIDI(f, patterns...); err != nil {
		log.Fatal(err)
	}
	fmt.Println("drumbeat.mid generated off of", patternStr)
	// Output: drumbeat.mid generated off of x.xxx...x...x.x.x..x.xx.x..x.xxx
	f.Close()
	os.Remove(f.Name())
}

func ExampleFromMIDI() {
	f, err := os.Open("fixtures/singlePattern.mid")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	patterns, err := drumbeat.FromMIDI(f)
	if err != nil {
		log.Fatalf("Failed to parse the MIDI file - %v", err)
	}
	// Default to 1/16th grid
	fmt.Printf("%s: %s", patterns[0].Name, patterns[0].Pulses)
	// Output: C1: x.......x.......
}

func ExamplePulses_Offset() {
	patternStr := "x..xx..."
	patterns := drumbeat.NewFromString(patternStr)
	patterns[0].Offset(2)
	fmt.Println(patterns[0].Pulses)
	patterns[0].Offset(-2)
	fmt.Println(patterns[0].Pulses)
	patterns[0].Offset(-2)
	fmt.Println(patterns[0].Pulses)
	// Output: ..x..xx.
	// x..xx...
	// .xx...x.
}
