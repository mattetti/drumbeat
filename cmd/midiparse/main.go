package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mattetti/drumbeat"
)

var (
	flagSrc = flag.String("src", "", "MIDI file to parse and convert")
)

func main() {
	flag.Parse()
	if *flagSrc == "" {
		fmt.Println("")
	}
	f, err := os.Open(*flagSrc)
	if err != nil {
		log.Fatalf("Failed to read the MIDI file - %v", err)
	}
	defer f.Close()
	patterns, err := drumbeat.FromMIDI(f)
	if err != nil {
		log.Fatalf("Failed to parse the MIDI file - %v", err)
	}
	for _, p := range patterns {
		fmt.Printf("%s: %s\n", p.Name, p.Steps)
	}

}
