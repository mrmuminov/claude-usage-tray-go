package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// GenerateIconPNG renders the percentage number onto a 22x22 PNG icon.
// Color: green ≤50%, yellow ≤80%, red >80%.
func GenerateIconPNG(pct int) []byte {
	const size = 22

	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Transparent background
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	// Pick text color based on usage
	var fg color.Color
	switch {
	case pct > 80:
		fg = color.RGBA{255, 80, 80, 255} // red
	case pct > 50:
		fg = color.RGBA{255, 200, 60, 255} // yellow
	default:
		fg = color.RGBA{100, 220, 100, 255} // green
	}

	text := fmt.Sprintf("%d", pct)
	face := basicfont.Face7x13

	// Center text
	textW := font.MeasureString(face, text).Ceil()
	asc := face.Metrics().Ascent.Ceil()
	x := (size - textW) / 2
	y := (size+asc)/2 - 1

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fg),
		Face: face,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(text)

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}
