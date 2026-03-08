package generators

import (
	"image"
	"image/color"
	"math"
)

// Bar_GenerateCheckerBar1D generates a bar-shaped checkerboard pattern.
// Translated from bar_GenerateCheckerBar1D.m
func Bar_GenerateCheckerBar1D(size int, barWidth float64, angle float64, offset float64, ndivsL int, ndivsS int, phase float64) ([]*image.Gray, error) {
	// Returns two images: one and its inverse
	checkers := make([]*image.Gray, 2)
	checkers[0] = image.NewGray(image.Rect(0, 0, size, size))
	checkers[1] = image.NewGray(image.Rect(0, 0, size, size))

	centerX, centerY := float64(size)/2.0, float64(size)/2.0
	maxRadius := float64(size) / 2.0

	// Angle in radians (Matlab version uses counter-clockwise)
	rad := angle * math.Pi / 180.0
	cosA := math.Cos(rad)
	sinA := math.Sin(rad)
	phaseRad := phase * math.Pi / 180.0

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY

			// Distance from center (for circular aperture)
			r := math.Sqrt(dx*dx + dy*dy)
			if r > maxRadius {
				checkers[0].SetGray(x, y, color.Gray{Y: 0})
				checkers[1].SetGray(x, y, color.Gray{Y: 0})
				continue
			}

			// Coordinate transform: rotate and shift
			// xRot is position across bar's width (short axis)
			// yRot is position along bar's length (long axis)
			xRot := dx*cosA + dy*sinA - offset
			yRot := -dx*sinA + dy*cosA

			// Check if we are inside the bar
			if xRot >= -barWidth/2.0 && xRot <= barWidth/2.0 {
				// Normalize positions to [0, nDivs]
				// Short axis includes phase shift
				normX := (xRot + barWidth/2.0) / barWidth
				if phaseRad != 0 {
					// Apply phase shift relative to short axis
					normX = math.Mod(normX + phaseRad/(2*math.Pi), 1.0)
					if normX < 0 { normX += 1.0 }
				}
				shortIdx := int(normX * float64(ndivsS))

				// Long axis
				normY := (yRot + maxRadius) / (2 * maxRadius)
				longIdx := int(normY * float64(ndivsL))

				if (shortIdx+longIdx)%2 == 0 {
					checkers[0].SetGray(x, y, color.Gray{Y: 255})
					checkers[1].SetGray(x, y, color.Gray{Y: 0})
				} else {
					checkers[0].SetGray(x, y, color.Gray{Y: 0})
					checkers[1].SetGray(x, y, color.Gray{Y: 255})
				}
			} else {
				checkers[0].SetGray(x, y, color.Gray{Y: 0})
				checkers[1].SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	return checkers, nil
}
