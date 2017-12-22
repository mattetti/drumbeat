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

func ExamplePulse_Offset() {
	patternStr := "x..xx..."
	patterns := drumbeat.NewFromString(patternStr)
	patterns[0].Steps = patterns[0].Steps.Offset(2)
	fmt.Println(patterns[0].Steps)
	patterns[0].Steps = patterns[0].Steps.Offset(-2)
	fmt.Println(patterns[0].Steps)
	patterns[0].Steps = patterns[0].Steps.Offset(-2)
	fmt.Println(patterns[0].Steps)
	// Output: ..x..xx.
	// x..xx...
	// .xx...x.
}
