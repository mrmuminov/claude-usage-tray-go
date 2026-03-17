package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed fonts/Ubuntu-Regular.ttf
var ubuntuTTF []byte

const iconSize = 128

// pctColor returns the status color for a given percentage.
func pctColor(pct int) color.RGBA {
	switch {
	case pct > 80:
		return color.RGBA{255, 80, 80, 255}
	case pct > 50:
		return color.RGBA{255, 200, 60, 255}
	default:
		return color.RGBA{100, 220, 100, 255}
	}
}

// GenerateIconPNG renders the percentage number onto a 128x128 PNG icon
// with a rounded dark background for visibility on both light and dark taskbars.
func GenerateIconPNG(pct int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	// Rounded dark background
	bg := color.RGBA{30, 30, 30, 220}
	drawRoundedRect(img, 0, 0, iconSize, iconSize, 24, bg)

	fg := pctColor(pct)
	text := fmt.Sprintf("%d", pct)

	fontData := ubuntuTTF
	if len(fontData) == 0 {
		fontData = gobold.TTF
	}

	ttf, err := opentype.Parse(fontData)
	if err != nil {
		return generateIconFallback(text, fg)
	}

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

// GenerateMenuDotIcon creates a 16x16 colored circle icon for menu items.
func GenerateMenuDotIcon(pct int) []byte {
	const size = 16
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	c := pctColor(pct)
	cx, cy, r := float64(size)/2, float64(size)/2, float64(size)/2-1
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - cx + 0.5
			dy := float64(y) - cy + 0.5
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= r {
				img.SetRGBA(x, y, c)
			} else if dist <= r+1 {
				a := uint8(float64(c.A) * (r + 1 - dist))
				img.SetRGBA(x, y, color.RGBA{c.R, c.G, c.B, a})
			}
		}
	}

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

// GenerateMenuActionIcon creates a 16x16 icon with a symbol.
// Symbols: "refresh", "github", "quit"
func GenerateMenuActionIcon(symbol string) []byte {
	const size = 16
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	c := color.RGBA{160, 160, 160, 255}

	switch symbol {
	case "refresh":
		cx, cy, r := 8.0, 8.0, 5.0
		for angle := 0.4; angle < 5.6; angle += 0.04 {
			px := cx + r*math.Cos(angle)
			py := cy + r*math.Sin(angle)
			setPixelAA(img, px, py, c)
			px2 := cx + (r-1)*math.Cos(angle)
			py2 := cy + (r-1)*math.Sin(angle)
			setPixelAA(img, px2, py2, c)
		}
		// arrowhead
		ex := cx + r*math.Cos(0.4)
		ey := cy + r*math.Sin(0.4)
		for i := 0.0; i < 4; i++ {
			setPixelAA(img, ex+i, ey-1, c)
			setPixelAA(img, ex, ey-1-i, c)
		}

	case "github":
		// Simple star/asterisk shape
		for i := 3; i <= 12; i++ {
			setPixelAA(img, 8, float64(i), c)  // vertical
			setPixelAA(img, float64(i), 8, c)   // horizontal
		}
		for i := 0; i < 8; i++ {
			setPixelAA(img, float64(4+i), float64(4+i), c)   // diagonal \
			setPixelAA(img, float64(11-i), float64(4+i), c)  // diagonal /
		}

	case "quit":
		for i := 0; i < 10; i++ {
			fi := float64(i)
			setPixelAA(img, 3+fi, 3+fi, c)
			setPixelAA(img, 4+fi, 3+fi, c)
			setPixelAA(img, 12-fi, 3+fi, c)
			setPixelAA(img, 11-fi, 3+fi, c)
		}
	}

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func setPixelAA(img *image.RGBA, x, y float64, c color.RGBA) {
	ix, iy := int(x), int(y)
	if ix >= 0 && ix < img.Bounds().Dx() && iy >= 0 && iy < img.Bounds().Dy() {
		img.SetRGBA(ix, iy, c)
	}
}

func drawRoundedRect(img *image.RGBA, x0, y0, x1, y1, radius int, c color.RGBA) {
	for y := y0; y < y1; y++ {
		for x := x0; x < x1; x++ {
			if isInRoundedRect(x, y, x0, y0, x1, y1, radius) {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

func isInRoundedRect(x, y, x0, y0, x1, y1, r int) bool {
	if x >= x0+r && x < x1-r {
		return true
	}
	if y >= y0+r && y < y1-r {
		return true
	}
	corners := [][2]int{
		{x0 + r, y0 + r},
		{x1 - r - 1, y0 + r},
		{x0 + r, y1 - r - 1},
		{x1 - r - 1, y1 - r - 1},
	}
	for _, corner := range corners {
		dx := x - corner[0]
		dy := y - corner[1]
		if dx*dx+dy*dy <= r*r {
			return true
		}
	}
	return false
}

func generateIconFallback(text string, fg color.RGBA) []byte {
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
