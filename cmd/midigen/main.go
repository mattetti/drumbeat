package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mattetti/drumbeat"
)

func main() {
	flag.Parse()
	patternStr := flag.String("pattern", "x...x...", "The pattern to convert to MIDI")
	flag.Parse()
	patterns := drumbeat.NewFromString(*patternStr)
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
}
