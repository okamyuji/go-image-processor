package processor

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T) string {
	// システムの一時ディレクトリ内にテスト用ディレクトリを作成
	testDir, err := os.MkdirTemp("", "image-processor-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp test directory: %v", err)
	}
	return testDir
}

// 単一のテスト画像を生成する関数
func generateSingleTestImage(outputPath string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// チェッカーパターンを描画
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

	// ファイルを作成
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// JPEGとして保存
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 90})
}

func TestResizeImage(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	// 入力画像のパス
	testInputPath := filepath.Join(testDir, "test_input.jpg")
	// 出力画像のパス
	testOutputPath := filepath.Join(testDir, "test_output_resize.jpg")

	// テスト画像を生成
	err := generateSingleTestImage(testInputPath, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image: %v", err)
	}

	// リサイズを実行
	err = ResizeImage(testInputPath, testOutputPath, 50, 50)
	if err != nil {
		t.Fatalf("Failed to resize image: %v", err)
	}

	// 結果を検証
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

	// 入力画像と出力画像のパス
	testInputPath := filepath.Join(testDir, "test_input.jpg")
	testOutputPath := filepath.Join(testDir, "test_output_rotate.jpg")

	// テスト画像を生成
	err := generateSingleTestImage(testInputPath, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image: %v", err)
	}

	// 画像を回転
	err = RotateImage(testInputPath, testOutputPath, 90)
	if err != nil {
		t.Fatalf("Failed to rotate image: %v", err)
	}

	// 結果を検証
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

	// 入力画像のパス
	testInputPath := filepath.Join(testDir, "test_input.jpg")
	testOutputPath := filepath.Join(testDir, "test_output_binarize.jpg")

	// テスト画像を生成
	err := generateSingleTestImage(testInputPath, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image: %v", err)
	}

	// 二値化を実行
	err = BinarizeImage(testInputPath, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to binarize image: %v", err)
	}

	// 結果を検証
	_, err = os.Stat(testOutputPath)
	if os.IsNotExist(err) {
		t.Errorf("Binarized image was not created")
	}
}

func TestConcatenateImagesHorizontally(t *testing.T) {
	testDir := setupTestDir(t)
	defer os.RemoveAll(testDir)

	// 入力画像と出力画像のパス
	testInputPath1 := filepath.Join(testDir, "test_input1.jpg")
	testInputPath2 := filepath.Join(testDir, "test_input2.jpg")
	testOutputPath := filepath.Join(testDir, "test_output_concathorz.jpg")

	// テスト画像を生成
	err := generateSingleTestImage(testInputPath1, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 1: %v", err)
	}

	err = generateSingleTestImage(testInputPath2, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 2: %v", err)
	}

	// 画像を水平方向に連結
	err = ConcatenateImagesHorizontally([]string{testInputPath1, testInputPath2}, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to concatenate images horizontally: %v", err)
	}

	// 結果を検証
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

	// 入力画像のパス
	testInputPath1 := filepath.Join(testDir, "test_input1.jpg")
	testInputPath2 := filepath.Join(testDir, "test_input2.jpg")
	testOutputPath := filepath.Join(testDir, "test_output_concatvert.jpg")

	// テスト画像を生成
	err := generateSingleTestImage(testInputPath1, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 1: %v", err)
	}
	err = generateSingleTestImage(testInputPath2, 100, 100)
	if err != nil {
		t.Fatalf("Failed to generate test image 2: %v", err)
	}

	// 画像を連結
	err = ConcatenateImagesVertically([]string{testInputPath1, testInputPath2}, testOutputPath)
	if err != nil {
		t.Fatalf("Failed to concatenate images vertically: %v", err)
	}

	// 結果を検証
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
