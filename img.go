package drumbeat

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"

	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
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
	labelWidth := 7 * stepWidth

	hitFill := color.NRGBA{124, 178, 227, 255}
	hitStroke := color.NRGBA{30, 30, 30, 255}
	labelBgColor := color.NRGBA{135, 135, 135, 255}

	// backgrounds/strokes alternate per row
	// other is lighter
	bgColor := color.NRGBA{165, 165, 165, 255}
	bgColorOther := color.NRGBA{149, 149, 149, 255}
	gridStrokeColor := color.NRGBA{153, 153, 153, 255}
	gridStrokeColorOther := color.NRGBA{133, 133, 133, 255}

	altBgColor := color.NRGBA{158, 158, 158, 255}
	altBgColorOther := color.NRGBA{143, 143, 143, 255}
	altGridStrokeColorOther := color.NRGBA{138, 138, 138, 255}
	altGridStrokeColor := color.NRGBA{147, 147, 147, 255}

	width := labelWidth + (nbrSteps * stepWidth)
	height := len(patterns) * stepHeight

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// white background
	draw.Draw(img, image.Rect(0, 0, labelWidth, height), image.NewUniform(labelBgColor), image.ZP, draw.Over)
	// grid background
	draw.Draw(img, image.Rect(labelWidth, 0, width, height), image.NewUniform(bgColor), image.ZP, draw.Over)

	var isOtherRow bool

	// draw the underlying grid
	for patternIDX, pattern := range patterns {
		stepsInBeat := int(pattern.Grid.StepsInBeat())
		patternY := patternIDX * stepHeight
		// line separating each label
		for x := 0; x < labelWidth; x++ {
			img.Set(x, patternY, gridStrokeColor)
		}

		var (
			isAltBeat    bool
			strokeColor  color.NRGBA
			bgPaintColor color.NRGBA
		)
		isOtherRow = false

		// background alternates per row
		if patternIDX%2 != 0 {
			isOtherRow = true
		}

		// vertical grid lines
		for pulseIDX := range pattern.Pulses {
			x := labelWidth + (pulseIDX * stepWidth)

			// detect beats and change the color
			if pulseIDX%stepsInBeat == 0 {
				// paint the steps for the entire beat
				if isAltBeat {
					if isOtherRow {
						bgPaintColor = altBgColorOther
						strokeColor = altGridStrokeColorOther
					} else {
						bgPaintColor = altBgColor
						strokeColor = altGridStrokeColor
					}
				} else {
					if isOtherRow {
						bgPaintColor = bgColorOther
						strokeColor = gridStrokeColorOther
					} else {
						bgPaintColor = bgColor
						strokeColor = gridStrokeColor
					}
				}

				// draw the row/beat background
				draw.Draw(img,
					// top lef, bottom right
					image.Rect(x, patternY, x+(stepWidth*stepsInBeat), patternY+stepHeight),
					image.NewUniform(bgPaintColor), image.ZP, draw.Over)
				isAltBeat = !isAltBeat
			}
			for h := 0; h < stepHeight; h++ {
				y := patternY + h
				// left
				img.Set(x, y, strokeColor)
				// right
				img.Set(x+stepWidth, y, strokeColor)
			}
		}
	}
	// bottom grid line line
	for x := 0; x < width; x++ {
		img.Set(x, stepHeight*(len(patterns)), gridStrokeColor)
	}

	for patternIDX, pattern := range patterns {
		patternY := patternIDX * stepHeight

		isOtherRow = false
		var bottomY int
		var heightToPaint int
		if patternIDX%2 != 0 {
			isOtherRow = true
		}

		if isOtherRow {
			bottomY = patternY + stepHeight - 2
			heightToPaint = stepHeight - 1
		} else {
			heightToPaint = stepHeight
			bottomY = patternY + stepHeight - 1
		}

		for pulseIDX, pulse := range pattern.Pulses {
			if pulse != nil {
				for w := 0; w < stepWidth+1; w++ {
					x := labelWidth + (pulseIDX * stepWidth) + w
					if w == 0 || w == stepWidth {
						// vertical stokes at the beginning and end of the pulse
						for h := 0; h < heightToPaint; h++ {
							y := patternY + h
							img.Set(x, y, hitStroke)
						}
					} else {
						for h := 0; h < stepHeight-1; h++ {
							y := patternY + h
							img.Set(x, y, hitFill)
						}
						// horizontal stokes
						img.Set(x, patternY, hitStroke)
						img.Set(x, bottomY, hitStroke)
					}

				}
			}
		}
		addLabel(img, 5, patternY+15, pattern.Name)
	}
	return png.Encode(w, img)
}

func addLabel(img *image.RGBA, x, y int, label string) {
	// truncate the labels to fit
	if len(label) > 16 {
		label = label[:16]
	}
	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: inconsolata.Regular8x16,
		Dot:  point,
	}
	d.DrawString(label)
}
