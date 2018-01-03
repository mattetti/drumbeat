package drumbeat

import (
	"image"
	"image/color"
	"image/png"
	"io"
)

// SaveAsPNG converts the patterns into an image.
func SaveAsPNG(w io.Writer, patterns []*Pattern) error {
	if len(patterns) < 1 {
		return nil
	}
	for _, pat := range patterns {
		pat.ReAlign()
	}
	nbrSteps := len(patterns[0].Pulses)
	stepHeight := 20
	stepWidth := 20

	gridColor := color.NRGBA{0, 0, 0, 255}
	missColor := color.NRGBA{255, 255, 255, 255}
	hitColor := color.NRGBA{255, 0, 0, 255}

	width := nbrSteps * stepWidth
	height := len(patterns) * stepHeight

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for patternIDX, pattern := range patterns {
		patternY := patternIDX * stepHeight
		for pulseIDX, pulse := range pattern.Pulses {
			c := missColor
			if pulse != nil {
				c = hitColor
			}
			for w := 0; w < stepWidth; w++ {
				x := (pulseIDX * stepWidth) + w
				for h := 0; h < stepHeight; h++ {
					y := patternY + h
					img.Set(x, y, c)
				}
			}
		}
		// horizontal grid
		for x := 0; x < width; x++ {
			img.Set(x, patternY, gridColor)
		}
	}
	// bottom line
	for x := 0; x < width; x++ {
		if x%stepWidth == 0 {
			for y := 0; y < height; y++ {
				img.Set(x, y, gridColor)
			}
		}
		img.Set(x, height-1, gridColor)
	}

	return png.Encode(w, img)
}
