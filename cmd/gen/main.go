package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"

	"github.com/mattetti/drumbeat"
)

func main() {
	grid := drumbeat.One16
	pattern := &drumbeat.Pattern{Grid: grid, PPQN: drumbeat.DefaultPPQN}
	pattern.Pulses = []*drumbeat.Pulse{}
	if err := saveMIDI("Hiphop", hipHopPatterns); err != nil {
		panic(err)
	}
	if err := savePNG("HipHop", hipHopPatterns); err != nil {
		panic(err)
	}
	if err := saveGIF("HipHop", 96, hipHopPatterns); err != nil {
		panic(err)
	}
	if err := saveMIDI("dubStep", dubStepPatterns); err != nil {
		panic(err)
	}
}

var (
	hipHopPatterns = drumbeat.NewFromString(drumbeat.One16, `
[kick]	{C1}	x.x.......xx...x	x.x.....x......x;
[snare]	{D1}	....x.......x...	....x.......x...;
[hihat]	{F#1}	x.x.x.x.x.x.x.x.	x.x.x.x.x.x.x.x.
	`)
	// kick; snare; hhc; hho
	dubStepPatterns = drumbeat.NewFromString(drumbeat.One16, `
[kick]		{C1}	x.........x.....	x..x..x...x.....;
[snare]		{D1}	........x.......	........x.......;
[hihat]		{F#1}	.xx...x....x..x.	.xx...x....x..x.;
[hihat open]{F#1}	....x........x..	....x........x..
	`)
)

func savePNG(name string, patterns []*drumbeat.Pattern) error {
	f, err := os.Create(fmt.Sprintf("%s.png", name))
	if err != nil {
		return fmt.Errorf("Failed to create img file - %v", err)
	}
	defer f.Close()
	return drumbeat.SaveAsPNG(f, patterns)
}

func saveGIF(name string, bpm float64, patterns []*drumbeat.Pattern) error {
	if len(patterns) < 1 {
		return nil
	}
	for _, pat := range patterns {
		pat.ReAlign()
	}
	nbImg := len(patterns[0].Pulses)
	width := 120
	height := 120

	// kick
	f, err := os.Open("circle-white.gif")
	if err != nil {
		panic(err)
	}
	kickImg, err := gif.Decode(f)
	if err != nil {
		panic(err)
	}
	f.Close()

	// snare
	f, err = os.Open("splatter.gif")
	if err != nil {
		panic(err)
	}
	snareImg, err := gif.Decode(f)
	if err != nil {
		panic(err)
	}
	f.Close()

	// hihat
	f, err = os.Open("center-splat-white.gif")
	if err != nil {
		panic(err)
	}
	hhImg, err := gif.Decode(f)
	if err != nil {
		panic(err)
	}
	f.Close()

	pal := make(color.Palette, 256)
	copy(pal, palette.Plan9)
	pal[0] = image.Black
	g := &gif.GIF{
		Image:           make([]*image.Paletted, 0, nbImg),
		Delay:           make([]int, 0, nbImg),
		Disposal:        make([]byte, 0, nbImg),
		Config:          image.Config{ColorModel: pal, Width: width, Height: height},
		BackgroundIndex: 0,
	}

	frameTransition := int((60 / bpm) / float64(patterns[0].Grid.StepsInBeat()) * 100)
	for i := 0; i < nbImg; i++ {

		g.Delay = append(g.Delay, frameTransition)
		g.Disposal = append(g.Disposal, gif.DisposalPrevious)
		g.Image = append(g.Image, image.NewPaletted(image.Rect(0, 0, width, height), pal))
		img := g.Image[len(g.Image)-1]

		// let's pretend 0 is kick
		if patterns[0].Pulses[i] != nil {
			draw.DrawMask(img, img.Bounds(), kickImg, image.ZP, kickImg, image.ZP, draw.Over)
		}
		// and 1 is snare
		if len(patterns) > 1 && patterns[1].Pulses[i] != nil {
			draw.DrawMask(img, img.Bounds(), snareImg, image.ZP, snareImg, image.ZP, draw.Over)
		}
		// and 2 is the hihat
		if len(patterns) > 2 && patterns[2].Pulses[i] != nil {
			draw.DrawMask(img, img.Bounds(), hhImg, image.ZP, hhImg, image.ZP, draw.Over)
		}
	}

	f, err = os.Create("drumbeat.gif")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	return gif.EncodeAll(f, g)
}

func saveMIDI(name string, patterns []*drumbeat.Pattern) error {
	fmt.Println(name)
	for i, pattern := range patterns {
		fmt.Println(i, pattern.Pulses)
	}
	f, err := os.Create(fmt.Sprintf("%s.mid", name))
	if err != nil {
		return fmt.Errorf("something wrong happened when creating the MIDI file - %v", err)
	}
	defer f.Close()
	return drumbeat.ToMIDI(f, patterns...)
}

// HLine draws a horizontal line
func HLine(img *image.Paletted, col color.Color, x1, y, x2 int) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

// VLine draws a vertical line
func VLine(img *image.Paletted, col color.Color, x, y1, y2 int) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

// Rect draws a rectangle utilizing HLine() and VLine()
func Rect(img *image.Paletted, col color.Color, x1, y1, x2, y2 int) {
	HLine(img, col, x1, y1, x2)
	HLine(img, col, x1, y2, x2)
	VLine(img, col, x1, y1, y2)
	VLine(img, col, x2, y1, y2)
}

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.NRGBA{122, 90, 30, 255}
	}
	return color.NRGBA{255, 255, 255, 255}
}
