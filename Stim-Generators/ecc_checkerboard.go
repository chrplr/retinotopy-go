package generators

import (
	"image"
	"image/color"
	"math"
)

// ECC_GenerateCheckerBoard1D generates a circular checkerboard pattern.
// Translated and simplified from ecc_GenerateCheckerBoard1D.m
func ECC_GenerateCheckerBoard1D(size int, rmin, rmax float64, width float64, startangle float64, nwedges int, nrings int, phase float64) ([]*image.Gray, error) {
	// Returns two images: one and its inverse
	checkers := make([]*image.Gray, 2)
	checkers[0] = image.NewGray(image.Rect(0, 0, size, size))
	checkers[1] = image.NewGray(image.Rect(0, 0, size, size))

	centerX, centerY := float64(size)/2.0, float64(size)/2.0

	// Convert to radians and normalize
	startRad := math.Mod(startangle*math.Pi/180.0, 2*math.Pi)
	if startRad < 0 {
		startRad += 2 * math.Pi
	}
	widthRad := width * math.Pi / 180.0
	endRad := math.Mod(startRad+widthRad, 2*math.Pi)
	phaseRad := phase * math.Pi / 180.0

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

			if width >= 360.0 {
				inAngle = true
			} else if startRad <= endRad {
				inAngle = theta >= startRad && theta <= endRad
			} else {
				inAngle = theta >= startRad || theta <= endRad
			}

			if inRadius && inAngle {
				// Normalized theta relative to startangle and phase
				normTheta := math.Mod(theta-startRad+phaseRad, 2*math.Pi)
				if normTheta < 0 {
					normTheta += 2 * math.Pi
				}

				// Checker logic
				wedgeIdx := int(normTheta / widthRad * float64(nwedges))
				ringIdx := int((r - rmin) / (rmax - rmin) * float64(nrings))

				if (wedgeIdx+ringIdx)%2 == 0 {
					checkers[0].SetGray(x, y, color.Gray{Y: 255})
					checkers[1].SetGray(x, y, color.Gray{Y: 0})
				} else {
					checkers[0].SetGray(x, y, color.Gray{Y: 0})
					checkers[1].SetGray(x, y, color.Gray{Y: 255})
				}
			} else {
				// Background (Mid-gray or transparent/black depending on use)
				// For a checkerboard generator, usually black (0) outside
				checkers[0].SetGray(x, y, color.Gray{Y: 0})
				checkers[1].SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	return checkers, nil
}
