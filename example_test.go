package drumbeat_test

import (
	"fmt"
	"log"
	"os"

	"github.com/mattetti/drumbeat"
)

func ExampleNewFromString() {
	patterns := drumbeat.NewFromString(drumbeat.One16, `
		[kick]	{C1}	x.x.......xx...x	x.x.....x......x;
		[snare]	{D1}	....x.......x...	....x.......x...;
		[hihat]	{F#1}	x.x.x.x.x.x.x.x.	x.x.x.x.x.x.x.x.
	`)
	f, err := os.Create("drumbeat.mid")
	if err != nil {
		log.Println("something wrong happened when creating the MIDI file", err)
		os.Exit(1)
	}
	if err := drumbeat.ToMIDI(f, patterns...); err != nil {
		log.Fatal(err)
	}
	f.Close()
	fmt.Println("drumbeat.mid generated")

	imgf, err := os.Create("drumbeat.png")
	if err != nil {
		log.Println("something wrong happened when creating the image file", err)
		os.Exit(1)
	}
	if err := drumbeat.SaveAsPNG(imgf, patterns); err != nil {
		log.Fatal(err)
	}
	imgf.Close()
	fmt.Println("drumbeat.png generated")
	// Output: drumbeat.mid generated
	// drumbeat.png generated

	os.Remove(f.Name())
	os.Remove(imgf.Name())
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
	patterns := drumbeat.NewFromString(drumbeat.One8, patternStr)
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
