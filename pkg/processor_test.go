package processor

import (
	"image"
	"os"
	"testing"
)

func TestResizeImage(t *testing.T) {
	// Generate a test image
	testInputPath := "test_input.jpg"
	err := GenerateTestImage(testInputPath, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image: %v", err)
	}
	defer os.Remove(testInputPath)

	// Resize the image
	testOutputPath := "test_output_resize.jpg"
	err = ResizeImage(testInputPath, testOutputPath, 50, 50)
	if err != nil {
		t.Fatalf("Failed to resize image: %v", err)
	}
	defer os.Remove(testOutputPath)

	// Check if the resized image has the correct dimensions
	resizedImg, err := os.Open(testOutputPath)
	if err != nil {
		t.Fatalf("Failed to open resized image: %v", err)
	}
	defer resizedImg.Close()

	img, _, err := image.Decode(resizedImg)
	if err != nil {
		t.Fatalf("Failed to decode resized image: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("Resized image dimensions incorrect. Expected 50x50, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestDenoiseImage(t *testing.T) {
	// Generate a test image
	testInputPath := "test_input.jpg"
	err := GenerateTestImage(testInputPath, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image: %v", err)
	}
	defer os.Remove(testInputPath)

	// Denoise the image
	testOutputPath := "test_output_denoise.jpg"
	err = DenoiseImage(testInputPath, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to denoise image: %v", err)
	}
	defer os.Remove(testOutputPath)

	// Check if the denoised image exists (we can't easily check the denoising effect)
	_, err = os.Stat(testOutputPath)
	if os.IsNotExist(err) {
		t.Errorf("Denoised image was not created")
	}
}

func TestRotateImage(t *testing.T) {
	// Generate a test image
	testInputPath := "test_input.jpg"
	err := GenerateTestImage(testInputPath, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image: %v", err)
	}
	defer os.Remove(testInputPath)

	// Rotate the image
	testOutputPath := "test_output_rotate.jpg"
	err = RotateImage(testInputPath, testOutputPath, 90)
	if err != nil {
		t.Fatalf("Failed to rotate image: %v", err)
	}
	defer os.Remove(testOutputPath)

	// Check if the rotated image exists and has the correct dimensions
	rotatedImg, err := os.Open(testOutputPath)
	if err != nil {
		t.Fatalf("Failed to open rotated image: %v", err)
	}
	defer rotatedImg.Close()

	img, _, err := image.Decode(rotatedImg)
	if err != nil {
		t.Fatalf("Failed to decode rotated image: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("Rotated image dimensions incorrect. Expected 100x100, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestBinarizeImage(t *testing.T) {
	// Generate a test image
	testInputPath := "test_input.jpg"
	err := GenerateTestImage(testInputPath, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image: %v", err)
	}
	defer os.Remove(testInputPath)

	// Binarize the image
	testOutputPath := "test_output_binarize.jpg"
	err = BinarizeImage(testInputPath, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to binarize image: %v", err)
	}
	defer os.Remove(testOutputPath)

	// Check if the binarized image exists
	_, err = os.Stat(testOutputPath)
	if os.IsNotExist(err) {
		t.Errorf("Binarized image was not created")
	}
}

func TestConcatenateImagesVertically(t *testing.T) {
	// Generate two test images
	testInputPath1 := "test_input1.jpg"
	testInputPath2 := "test_input2.jpg"
	err := GenerateTestImage(testInputPath1, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 1: %v", err)
	}
	defer os.Remove(testInputPath1)

	err = GenerateTestImage(testInputPath2, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 2: %v", err)
	}
	defer os.Remove(testInputPath2)

	// Concatenate the images vertically
	testOutputPath := "test_output_concatvert.jpg"
	err = ConcatenateImagesVertically([]string{testInputPath1, testInputPath2}, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to concatenate images vertically: %v", err)
	}
	defer os.Remove(testOutputPath)

	// Check if the concatenated image exists and has the correct dimensions
	concatenatedImg, err := os.Open(testOutputPath)
	if err != nil {
		t.Fatalf("Failed to open concatenated image: %v", err)
	}
	defer concatenatedImg.Close()

	img, _, err := image.Decode(concatenatedImg)
	if err != nil {
		t.Fatalf("Failed to decode concatenated image: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 200 {
		t.Errorf("Concatenated image dimensions incorrect. Expected 100x200, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestConcatenateImagesHorizontally(t *testing.T) {
	// Generate two test images
	testInputPath1 := "test_input1.jpg"
	testInputPath2 := "test_input2.jpg"
	err := GenerateTestImage(testInputPath1, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 1: %v", err)
	}
	defer os.Remove(testInputPath1)

	err = GenerateTestImage(testInputPath2, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 2: %v", err)
	}
	defer os.Remove(testInputPath2)

	// Concatenate the images horizontally
	testOutputPath := "test_output_concathorz.jpg"
	err = ConcatenateImagesHorizontally([]string{testInputPath1, testInputPath2}, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to concatenate images horizontally: %v", err)
	}
	defer os.Remove(testOutputPath)

	// Check if the concatenated image exists and has the correct dimensions
	concatenatedImg, err := os.Open(testOutputPath)
	if err != nil {
		t.Fatalf("Failed to open concatenated image: %v", err)
	}
	defer concatenatedImg.Close()

	img, _, err := image.Decode(concatenatedImg)
	if err != nil {
		t.Fatalf("Failed to decode concatenated image: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 200 || bounds.Dy() != 100 {
		t.Errorf("Concatenated image dimensions incorrect. Expected 200x100, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}
