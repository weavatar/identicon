// Package identicon generates identicons based on a hash.
package identicon

import (
	"hash"
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	"math"
)

// IdentIcon represents an identicon generator
type IdentIcon struct {
	sqSize int
	rows   int
	cols   int
	h      hash.Hash64
	maxX   int
	maxY   int
}

// New creates a new identicon renderer
func New(size, rows, cols int) *IdentIcon {
	return &IdentIcon{
		sqSize: size / max(rows, cols),
		rows:   rows,
		cols:   cols,
		h:      fnv.New64a(),
		maxX:   size,
		maxY:   size,
	}
}

// Make creates an identicon image based on the input hash
func (icon *IdentIcon) Make(hash []byte) image.Image {
	icon.h.Reset()
	if _, err := icon.h.Write(hash); err != nil {
		panic(err)
	}
	h := icon.h.Sum64()

	// Generate foreground color with better contrast
	hue := float64(h%360) / 360.0
	saturation := 0.5 + float64(h%1000)/2000.0
	brightness := 0.5 + float64(h%1000)/2000.0

	r, g, b := hsvToRgb(hue, saturation, brightness)
	fgColor := color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: 255,
	}

	// Background color (light neutral color)
	bgColor := color.RGBA{R: 240, G: 240, B: 240, A: 255}

	// Create image and fill with background color
	img := image.NewRGBA(image.Rect(0, 0, icon.maxX, icon.maxY))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	// Calculate center area for shapes
	margin := int(float64(icon.maxX) * 0.1)
	innerSize := icon.maxX - 2*margin
	cellSize := innerSize / icon.cols

	// Generate a symmetric pattern
	pattern := generateSymmetricPattern(h, icon.rows, icon.cols)

	// Draw the pattern
	for y := 0; y < icon.rows; y++ {
		for x := 0; x < icon.cols; x++ {
			if pattern[y][x] {
				drawShape(img, x, y, cellSize, margin, fgColor, int(h%7))
			}
		}
	}

	return img
}

// Generate a symmetric pattern based on the hash
func generateSymmetricPattern(hash uint64, rows, cols int) [][]bool {
	pattern := make([][]bool, rows)
	for i := range pattern {
		pattern[i] = make([]bool, cols)
	}

	// Generate the left half (or slightly more than half for odd dimensions)
	middleCol := cols / 2
	if cols%2 == 1 {
		middleCol++
	}

	// Fill the left part of the pattern
	bits := hash
	for y := 0; y < rows; y++ {
		for x := 0; x < middleCol; x++ {
			pattern[y][x] = (bits & 1) == 1
			bits >>= 1

			// Mirror horizontally (left to right)
			if x < cols/2 {
				pattern[y][cols-x-1] = pattern[y][x]
			}
		}
	}

	return pattern
}

// Draw a shape at the specified position
func drawShape(img *image.RGBA, x, y, cellSize, margin int, color color.RGBA, shapeType int) {
	startX := margin + x*cellSize
	startY := margin + y*cellSize

	switch shapeType {
	case 0:
		// Fill square
		drawRect(img, startX, startY, cellSize, cellSize, color)
	case 1:
		// Circle
		drawCircle(img, startX+cellSize/2, startY+cellSize/2, cellSize/2, color)
	case 2:
		// Diamond
		drawDiamond(img, startX, startY, cellSize, color)
	case 3:
		// Triangle pointing up
		drawTriangle(img, startX, startY, cellSize, 0, color)
	case 4:
		// Triangle pointing right
		drawTriangle(img, startX, startY, cellSize, 1, color)
	case 5:
		// Triangle pointing down
		drawTriangle(img, startX, startY, cellSize, 2, color)
	case 6:
		// Triangle pointing left
		drawTriangle(img, startX, startY, cellSize, 3, color)
	}
}

// Draw a filled rectangle
func drawRect(img *image.RGBA, x, y, width, height int, color color.RGBA) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			img.Set(x+dx, y+dy, color)
		}
	}
}

// Draw a filled circle
func drawCircle(img *image.RGBA, centerX, centerY, radius int, color color.RGBA) {
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= radius*radius {
				img.Set(centerX+dx, centerY+dy, color)
			}
		}
	}
}

// Draw a filled diamond
func drawDiamond(img *image.RGBA, x, y, size int, color color.RGBA) {
	halfSize := size / 2
	centerX := x + halfSize
	centerY := y + halfSize

	for dy := 0; dy < size; dy++ {
		width := size - abs(dy-halfSize)*2
		startX := centerX - width/2

		for dx := 0; dx < width; dx++ {
			img.Set(startX+dx, centerY+(dy-halfSize), color)
		}
	}
}

// Draw a filled triangle with specified orientation (0=up, 1=right, 2=down, 3=left)
func drawTriangle(img *image.RGBA, x, y, size, orientation int, color color.RGBA) {
	switch orientation {
	case 0: // Up
		for dy := 0; dy < size; dy++ {
			width := size - dy*2
			startX := x + dy
			for dx := 0; dx < width; dx++ {
				img.Set(startX+dx, y+size-dy-1, color)
			}
		}
	case 1: // Right
		for dx := 0; dx < size; dx++ {
			height := size - dx*2
			startY := y + dx
			for dy := 0; dy < height; dy++ {
				img.Set(x+dx, startY+dy, color)
			}
		}
	case 2: // Down
		for dy := 0; dy < size; dy++ {
			width := size - dy*2
			startX := x + dy
			for dx := 0; dx < width; dx++ {
				img.Set(startX+dx, y+dy, color)
			}
		}
	case 3: // Left
		for dx := 0; dx < size; dx++ {
			height := size - dx*2
			startY := y + dx
			for dy := 0; dy < height; dy++ {
				img.Set(x+size-dx-1, startY+dy, color)
			}
		}
	}
}

// HSV to RGB conversion for better color generation
func hsvToRgb(h, s, v float64) (r, g, b float64) {
	if s == 0 {
		return v, v, v
	}

	h *= 6
	i := math.Floor(h)
	f := h - i
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	switch int(i) % 6 {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	default:
		return v, p, q
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
