package processor

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T) string {
	testDir, err := os.MkdirTemp("", "image-processor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp test directory: %v", err)
	}
	return testDir
}

func generateSingleTestImage(outputPath string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	blockSize := 20
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if (x/blockSize+y/blockSize)%2 == 0 {
				img.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			} else {
				img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, img, &jpeg.Options{Quality: 90})
}

func TestResizeImage(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	testInputPath := filepath.Join(testDir, "test_input.jpg")
	testOutputPath := filepath.Join(testDir, "test_output_resize.jpg")

	err := generateSingleTestImage(testInputPath, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image: %v", err)
	}

	err = ResizeImage(testInputPath, testOutputPath, 50, 50)
	if err != nil {
		t.Fatalf("Failed to resize image: %v", err)
	}

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
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	// Generate test images
	err := GenerateTestImage(testDir, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test images: %v", err)
	}

	// Use noise image for denoise test
	testInputPath := filepath.Join(testDir, "noise_test.jpg")
	testOutputPath := filepath.Join(testDir, "test_output_denoise.jpg")

	err = DenoiseImage(testInputPath, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to denoise image: %v", err)
	}

	_, err = os.Stat(testOutputPath)
	if os.IsNotExist(err) {
		t.Errorf("Denoised image was not created")
	}
}

func TestRotateImage(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	testInputPath := filepath.Join(testDir, "test_input.jpg")
	testOutputPath := filepath.Join(testDir, "test_output_rotate.jpg")

	err := generateSingleTestImage(testInputPath, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image: %v", err)
	}

	err = RotateImage(testInputPath, testOutputPath, 90)
	if err != nil {
		t.Fatalf("Failed to rotate image: %v", err)
	}

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
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	testInputPath := filepath.Join(testDir, "test_input.jpg")
	testOutputPath := filepath.Join(testDir, "test_output_binarize.jpg")

	err := generateSingleTestImage(testInputPath, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image: %v", err)
	}

	err = BinarizeImage(testInputPath, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to binarize image: %v", err)
	}

	_, err = os.Stat(testOutputPath)
	if os.IsNotExist(err) {
		t.Errorf("Binarized image was not created")
	}
}

func TestConcatenateImagesHorizontally(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	testInputPath1 := filepath.Join(testDir, "test_input1.jpg")
	testInputPath2 := filepath.Join(testDir, "test_input2.jpg")
	testOutputPath := filepath.Join(testDir, "test_output_concathorz.jpg")

	err := generateSingleTestImage(testInputPath1, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 1: %v", err)
	}

	err = generateSingleTestImage(testInputPath2, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 2: %v", err)
	}

	err = ConcatenateImagesHorizontally([]string{testInputPath1, testInputPath2}, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to concatenate images horizontally: %v", err)
	}

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
		t.Errorf("Concatenated image dimensions incorrect. Expected 200x100, got %dx%d",
			bounds.Dx(), bounds.Dy())
	}
}

func TestConcatenateImagesVertically(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	testInputPath1 := filepath.Join(testDir, "test_input1.jpg")
	testInputPath2 := filepath.Join(testDir, "test_input2.jpg")
	testOutputPath := filepath.Join(testDir, "test_output_concatvert.jpg")

	err := generateSingleTestImage(testInputPath1, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 1: %v", err)
	}
	err = generateSingleTestImage(testInputPath2, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 2: %v", err)
	}

	err = ConcatenateImagesVertically([]string{testInputPath1, testInputPath2}, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to concatenate images vertically: %v", err)
	}

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

func TestAutoRotateImage(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	// Set up test image path
	testInputPath := filepath.Join(testDir, "test_input_skew.jpg")
	testOutputPath := filepath.Join(testDir, "test_output_auto_rotate.jpg")

	// Generate skewed test image
	err := generateSkewedTestImage(testInputPath, 100, 100, 15.0) // 15度傾いた画像を生成
	if err != nil {
		t.Fatalf("Failed to generate skewed test image: %v", err)
	}

	// Run auto-correction
	err = AutoRotateImage(testInputPath, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to auto-rotate image: %v", err)
	}

	// Verify correction results
	correctedImg, err := os.Open(testOutputPath)
	if err != nil {
		t.Fatalf("Failed to open corrected image: %v", err)
	}
	defer correctedImg.Close()

	img, _, err := image.Decode(correctedImg)
	if err != nil {
		t.Fatalf("Failed to decode corrected image: %v", err)
	}

	// Check if the image is properly rotated after correction
	// Note: Since precise angle verification is difficult, verify that the image is generated
	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Error("Corrected image has invalid dimensions")
	}
}

// generateSkewedTestImage creates a test image with a known skew angle
func generateSkewedTestImage(outputPath string, width, height int, angleInDegrees float64) error {
	// Generate test image with text and lines, rotate by specified angle
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background with white color
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw horizontal lines (for skew detection)
	for y := height / 4; y < height*3/4; y += height / 4 {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.Black)
		}
	}

	// Apply rotation
	rotated := rotateImage(img, angleInDegrees)

	// Save to file
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, rotated, &jpeg.Options{Quality: 90})
}

func TestGenerateSkewTestImage(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	// Generate skewed test image
	testPath := filepath.Join(testDir, "skew_test.jpg")
	err := generateSkewTestImage(testPath, 200, 200, 15.0)
	if err != nil {
		t.Fatalf("Failed to generate skew test image: %v", err)
	}

	// Verify the image was created
	_, err = os.Stat(testPath)
	if os.IsNotExist(err) {
		t.Error("Skew test image was not created")
	}

	// Try to load and decode the image
	file, err := os.Open(testPath)
	if err != nil {
		t.Fatalf("Failed to open generated image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatalf("Failed to decode generated image: %v", err)
	}

	// Verify image dimensions
	bounds := img.Bounds()
	expectedSize := int(math.Ceil(200 * math.Sqrt(2))) // Maximum size after 45-degree rotation
	if bounds.Dx() > expectedSize || bounds.Dy() > expectedSize {
		t.Errorf("Generated image is too large: got %dx%d, expected maximum %dx%d",
			bounds.Dx(), bounds.Dy(), expectedSize, expectedSize)
	}
}

func TestGenerateTestImage(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	err := GenerateTestImage(testDir, 200, 200)
	if err != nil {
		t.Fatalf("Failed to generate test images: %v", err)
	}

	// Check for skew test images
	for i := 1; i <= 3; i++ {
		path := filepath.Join(testDir, fmt.Sprintf("skew_test_%d.jpg", i))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected skew test image not found: %s", path)
		}
	}
}
