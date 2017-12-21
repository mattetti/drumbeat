package drumbeat_test

import (
	"fmt"
	"log"
	"os"

	"github.com/mattetti/drumbeat"
)

func Example() {
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
