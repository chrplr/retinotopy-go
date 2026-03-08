package generators

import (
	"image"
	"image/color"
	"math"
)

// CreateWedgeMask creates a wedge-shaped mask.
// Translated from CreateWedgeMask.m
// rmin, rmax: radius in pixels
// startAngle: start angle in degrees
// widthAngle: width of the wedge in degrees
func CreateWedgeMask(size int, rmin, rmax, startAngle, widthAngle float64) ([]*image.Gray, error) {
	mask := make([]*image.Gray, 2)
	mask[0] = image.NewGray(image.Rect(0, 0, size, size))
	mask[1] = image.NewGray(image.Rect(0, 0, size, size))

	centerX, centerY := float64(size)/2.0, float64(size)/2.0

	// Convert to radians and normalize
	startRad := math.Mod(startAngle*math.Pi/180.0, 2*math.Pi)
	if startRad < 0 {
		startRad += 2 * math.Pi
	}
	widthRad := widthAngle * math.Pi / 180.0
	endRad := math.Mod(startRad+widthRad, 2*math.Pi)

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			r := math.Sqrt(dx*dx + dy*dy)
			theta := math.Atan2(dy, dx)
			if theta < 0 {
				theta += 2 * math.Pi
			}

			inRadius := r >= rmin && r <= rmax
			inAngle := false

			if widthAngle >= 360.0 {
				inAngle = true
			} else if startRad <= endRad {
				inAngle = theta >= startRad && theta <= endRad
			} else {
				inAngle = theta >= startRad || theta <= endRad
			}

			if inRadius && inAngle {
				mask[0].SetGray(x, y, color.Gray{Y: 255})
				mask[1].SetGray(x, y, color.Gray{Y: 0})
			} else {
				mask[0].SetGray(x, y, color.Gray{Y: 0})
				mask[1].SetGray(x, y, color.Gray{Y: 255})
			}
		}
	}

	return mask, nil
}
