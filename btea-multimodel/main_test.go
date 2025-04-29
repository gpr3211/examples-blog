package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
	"time"
)

// createTestImage creates a test PNG image with the specified dimensions
// and saves it to a temporary file
func createTestImage(width, height int) (string, error) {
	// Create a new image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill it with some test pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a simple pattern
			img.Set(x, y, color.RGBA{
				R: uint8((x * 255) / width),
				G: uint8((y * 255) / height),
				B: 100,
				A: 255,
			})
		}
	}

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "testimage_*.png")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	// Encode the image to PNG
	if err := png.Encode(tmpfile, img); err != nil {
		os.Remove(tmpfile.Name())
		return "", err
	}

	return tmpfile.Name(), nil
}

func TestProcessImage(t *testing.T) {
	// Create a test image (100x100 pixels)
	imgPath, err := createTestImage(100, 100)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	defer os.Remove(imgPath) // Clean up the test image when done

	// Redirect standard output temporarily to avoid cluttering test output
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Start processing the image
	progressChan := ProcessImage(imgPath)

	// Collect progress updates
	var progressValues []float64
	var lastProgress float64

	// Set up a timeout in case the channel never closes
	timeout := time.After(5 * time.Second)

	// Loop until we get all progress updates or timeout
loop:
	for {
		select {
		case progress, ok := <-progressChan:
			if !ok {
				// Channel was closed, we're done
				break loop
			}
			progressValues = append(progressValues, progress)
			lastProgress = progress
		case <-timeout:
			// Test timed out
			t.Fatal("Test timed out waiting for progress updates")
			break loop
		}
	}

	// Restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Consume the captured output if needed
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	// Validate the progress reports
	if len(progressValues) == 0 {
		t.Fatal("No progress updates were received")
	}

	// Check that progress increases monotonically
	for i := 1; i < len(progressValues); i++ {
		if progressValues[i] < progressValues[i-1] {
			t.Errorf("Progress decreased from %.1f%% to %.1f%%",
				progressValues[i-1], progressValues[i])
		}
	}

	// Check that we eventually reach 100%
	if lastProgress != 100.0 {
		t.Errorf("Final progress was %.1f%%, expected 100.0%%", lastProgress)
	}

	t.Logf("Received %d progress updates", len(progressValues))
	t.Logf("First few updates: %v", progressValues[:min(5, len(progressValues))])
	t.Logf("Last update: %.1f%%", lastProgress)
}

// TestProcessImageCancellation tests that the processing can be canceled
// by simply abandoning the progress channel
func TestProcessImageCancellation(t *testing.T) {
	// Create a test image (200x200 pixels - bigger to ensure it takes longer)
	imgPath, err := createTestImage(200, 200)
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	defer os.Remove(imgPath) // Clean up the test image when done

	// Redirect standard output temporarily
	originalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	// Start processing the image
	progressChan := ProcessImage(imgPath)

	// Read just a few progress updates and then abandon the channel
	for i := 0; i < 3; i++ {
		select {
		case progress, ok := <-progressChan:
			if !ok {
				// Channel closed too early
				t.Fatal("Progress channel closed unexpectedly")
			}
			t.Logf("Received progress update: %.1f%%", progress)
		case <-time.After(1 * time.Second):
			t.Fatal("Timed out waiting for initial progress updates")
		}
	}

	// We're explicitly NOT draining the channel here
	// This tests that the goroutine in ProcessImage doesn't leak

	// Restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Give some time for the goroutine to complete
	// In a real application, you'd want to use proper context cancellation
	time.Sleep(100 * time.Millisecond)

	// There's no real assertion here except that the test completes without deadlock
	t.Log("Test completed without deadlock")
}

// min returns the smaller of x or y
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
