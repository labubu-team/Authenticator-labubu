package app

import (
	"fmt"
	"github.com/kbinani/screenshot"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

// Window defines the region for capturing the screen
type Window struct {
	X, Y int
	W, H int
}

// CaptureScreen captures a portion of the screen as specified by the Window struct
// and saves it to a file, returning the file path.
func CaptureScreen(window Window) (string, error) {

	// Create a rectangle for the capture region
	rect := image.Rect(window.X, window.Y, window.X+window.W, window.Y+window.H)

	// Capture the specified region
	img, err := screenshot.CaptureRect(rect)
	if err != nil {
		return "", fmt.Errorf("error capturing screen: %v", err)
	}

	// Generate the file path in the user's home directory
	fileName := "screenshot.png"
	filePath := filepath.Join(os.Getenv("HOME"), fileName)

	// Create the file to save the screenshot
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("Error closing file: %v\n", closeErr)
		}
	}()

	// Encode the captured image to PNG format and save it
	err = png.Encode(file, img)
	if err != nil {
		return "", fmt.Errorf("error saving image: %v", err)
	}

	// Return the file path
	return filePath, nil
}
