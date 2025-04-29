package model

import (
	"fmt"
	"image/png"
	"os"
	"sync"
)

// createProgressTracker returns a function that tracks the progress of pixel processing
func createProgressTracker(totalPixels int) func(int) (float64, bool) {
	var processedPixels int
	var mu sync.Mutex
	lastReportedPercent := -1 // Track the last reported percentage

	return func(increment int) (float64, bool) {
		mu.Lock()
		defer mu.Unlock()

		processedPixels += increment
		progress := float64(processedPixels) / float64(totalPixels) * 100.0

		// Calculate the integer percentage (e.g., 62 for 62.1%)
		currentPercent := int(progress)

		// Determine if this is a new percentage to report
		reportChange := currentPercent > lastReportedPercent

		if reportChange {
			lastReportedPercent = currentPercent
		}

		return progress, reportChange
	}
}

func ProcessImage(src string) chan float64 {
	// Create buffered channel to avoid blocking
	outchan := make(chan float64, 101) // Buffer for 0-100%

	// Launch a goroutine to do the processing and send progress updates
	go func() {
		// Ensure the channel gets closed when we're done
		defer close(outchan)

		file, err := os.Open(src)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Decode the image
		img, err := png.Decode(file)
		if err != nil {
			panic(err)
		}

		// Get the image bounds
		bounds := img.Bounds()
		width, height := bounds.Max.X, bounds.Max.Y

		// Total number of pixels
		totalPixels := width * height

		// Create progress tracker
		trackProgress := createProgressTracker(totalPixels)

		// Process the image
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// This is where you would actually process the pixel if needed
				// pixel := img.At(x, y)
				// r, g, b, a := pixel.RGBA() // Values are 16-bit (0-65535)

				// Update progress after each pixel (increment by 1)
				progress, shouldReport := trackProgress(1)

				if shouldReport {
					outchan <- progress
					fmt.Printf("Processing: %.1f%% complete\n", progress)
				}
			}
		}

		// Make sure we send 100% at the end
	}()

	return outchan
}

/*
	fmt.Println("Starting image processing...")

	progChan := ProcessImage("input.png")

	// Read from the channel until it's closed
	for progress := range progChan {
		fmt.Printf("Main routine received: %.1f%%\n", progress)
	}

	fmt.Println("Processing complete!")
}
*/
