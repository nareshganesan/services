package shared

import (
	// "fmt"
	"image"
	"image/color"
	"image/draw"
)

// GetTrackerImage returns a simple pixel image used for email tracking
func GetTrackerImage() *image.Image {
	tmpImg := image.NewRGBA(image.Rect(0, 0, 1, 1))
	color := color.RGBA{0, 0, 0, 0}
	draw.Draw(tmpImg, tmpImg.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)
	var img image.Image = tmpImg
	return &img
}
