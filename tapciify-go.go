package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"math"
	"os"

	"github.com/nfnt/resize"
)

func main() {
	// You can register another format here
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	file, err := os.Open("./image.png")

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	defer file.Close()

	pixels, err := toAscii(file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(pixels)
}

// Get the bi-dimensional pixel array
func toAscii(file io.Reader) (string, error) {
	orgImg, _, err := image.Decode(file)
	img := resize.Resize(128, 0, orgImg, resize.Lanczos3)

	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var output = ""

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			char, err := toAsciiCharacter(getLightness(rgbaToPixel(img.At(x, y).RGBA())))

			if err != nil {
				return "", err
			}

			output += string(char)

		}

		output += "\n"
	}

	return output, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) RgbaPixel {
	return RgbaPixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// RgbaPixel struct example
type RgbaPixel struct {
	R int
	G int
	B int
	A int
}

func getLightness(p RgbaPixel) float32 {
	max := max(max(p.R, p.G), p.B)
	min := min(min(p.R, p.G), p.B)

	return float32((max+min)*p.A) / 130050
}

func toAsciiCharacter(lightness float32) (rune, error) {
	asciiString := " .,:;+*?%S#@"

	charCount := len(asciiString)
	index := int(math.Floor(float64(charCount-1) * float64(lightness)))

	if index < 0 || index >= charCount {
		return 0, fmt.Errorf("LIGHTNESS %f IS OUT OF ARRAY", lightness)
	}

	return rune(asciiString[index]), nil
}
