package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Place your custom Ubuntu font file as "ubuntu.ttf" in this directory.
// If it doesn't exist, falls back to Go Bold font.
//
//go:embed fonts/Ubuntu-Regular.ttf
var ubuntuTTF []byte

const iconSize = 128

// GenerateIconPNG renders the percentage number onto a 128x128 PNG icon.
// Color: green ≤50%, yellow ≤80%, red >80%.
func GenerateIconPNG(pct int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	var fg color.Color
	switch {
	case pct > 80:
		fg = color.RGBA{255, 80, 80, 255}
	case pct > 50:
		fg = color.RGBA{255, 200, 60, 255}
	default:
		fg = color.RGBA{100, 220, 100, 255}
	}

	text := fmt.Sprintf("%d", pct)

	fontData := ubuntuTTF
	if len(fontData) == 0 {
		fontData = gobold.TTF
	}

	ttf, err := opentype.Parse(fontData)
	if err != nil {
		return generateIconFallback(text, fg)
	}

	// Adjust font size for 3-digit numbers (100%)
	fontSize := 96.0
	if len(text) == 3 {
		fontSize = 66.0
	}

	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return generateIconFallback(text, fg)
	}
	defer face.Close()

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(fg),
		Face: face,
	}

	textW := d.MeasureString(text).Ceil()
	metrics := face.Metrics()
	asc := metrics.Ascent.Ceil()
	desc := metrics.Descent.Ceil()

	x := (iconSize - textW) / 2
	y := (iconSize + asc - desc) / 2

	d.Dot = fixed.P(x, y)
	d.DrawString(text)

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func generateIconFallback(text string, fg color.Color) []byte {
	const size = 22
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	face := basicfont.Face7x13
	textW := font.MeasureString(face, text).Ceil()
	asc := face.Metrics().Ascent.Ceil()
	x := (size - textW) / 2
	y := (size + asc) / 2

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
