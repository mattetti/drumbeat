# drumbeat

Drumbeat is a Go library to create/parse drum beat patterns.

[![GoDoc](https://godoc.org/github.com/mattetti/drumbeat?status.svg)](https://godoc.org/github.com/mattetti/drumbeat)
[![Go Report Card](https://goreportcard.com/badge/github.com/mattetti/drumbeat)](https://goreportcard.com/report/github.com/mattetti/drumbeat)
[![Coverage Status](https://codecov.io/gh/mattetti/drumbeat/graph/badge.svg)](https://codecov.io/gh/mattetti/drumbeat)
[![Build Status](https://travis-ci.org/mattetti/drumbeat.svg)](https://travis-ci.org/mattetti/drumbeat)

## Example

More information and documentation available at [![GoDoc](https://godoc.org/github.com/mattetti/drumbeat?status.svg)](https://godoc.org/github.com/mattetti/drumbeat).

```go
// define a pattern using text:
patterns := drumbeat.NewFromString(drumbeat.One16, `
	[kick]	{C1}	x.x.......xx...x	x.x.....x......x;
	[snare]	{D1}	....x.......x...	....x.......x...;
	[hihat]	{F#1}	x.x.x.x.x.x.x.x.	x.x.x.x.x.x.x.x.
`)

// convert to MIDI
f, err := os.Create("drumbeat.mid")
if err != nil {
    log.Println("something wrong happened when creating the MIDI file", err)
    os.Exit(1)
}
if err := drumbeat.ToMIDI(f, patterns...); err != nil {
    log.Fatal(err)
}
f.Close()

// generate a PNG visualization
imgf, err := os.Create("drumbeat.png")
if err != nil {
    log.Println("something wrong happened when creating the image file", err)
    os.Exit(1)
}
if err := drumbeat.SaveAsPNG(imgf, patterns); err != nil {
    log.Fatal(err)
}
imgf.Close()
```

![png output](https://github.com/mattetti/drumbeat/blob/master/example.png?raw=true)
