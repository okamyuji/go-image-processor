package processor

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log/slog"
	"math"
	"os"
	"sort"
	"testing"

	"github.com/nfnt/resize"
	"golang.org/x/exp/rand"

	"github.com/okamyuji/go-image-processor/config"
)

var cfg *config.Config

func init() {
	cfg = config.GetConfig()
	// JSONフォーマットのロガーを設定
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)
}

// ResizeImage resizes the input image to the specified width and height.
// It takes the paths of the input and output files, and the desired width and height.
// Returns an error if the operation fails.
func ResizeImage(inputPath string, outputPath string, width, height uint) error {
	slog.Info("resizing image",
		"input", inputPath,
		"width", width,
		"height", height)

	// Open the input file
	file, err := os.Open(inputPath)
	if err != nil {
		return &ErrInvalidInput{Path: inputPath}
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return &ErrProcessing{Op: "decode", Err: err}
	}

	// Calculate new dimensions while maintaining aspect ratio
	bounds := img.Bounds()
	origWidth := float64(bounds.Dx())
	origHeight := float64(bounds.Dy())
	ratio := origWidth / origHeight

	var newWidth, newHeight uint
	if float64(width)/float64(height) > ratio {
		// Height is the limiting factor
		newHeight = height
		newWidth = uint(float64(height) * ratio)
	} else {
		// Width is the limiting factor
		newWidth = width
		newHeight = uint(float64(width) / ratio)
	}

	// Resize the image
	resizedImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return &ErrInvalidOutput{Path: outputPath}
	}
	defer out.Close()

	// Encode and save the resized image
	if err := jpeg.Encode(out, resizedImg, &jpeg.Options{Quality: cfg.JpegQuality}); err != nil {
		return &ErrProcessing{Op: "encode", Err: err}
	}

	return nil
}

// DenoiseImage applies a simple denoising filter to the input image.
// It takes the paths of the input and output files.
// Returns an error if the operation fails.
func DenoiseImage(inputPath string, outputPath string) error {
	slog.Info("denoising image", "input", inputPath)

	// Open the input file
	file, err := os.Open(inputPath)
	if err != nil {
		return &ErrInvalidInput{Path: inputPath}
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return &ErrProcessing{Op: "decode", Err: err}
	}

	// Apply median filter for denoising
	bounds := img.Bounds()
	denoised := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			denoised.Set(x, y, medianFilter(img, x, y))
		}
	}

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return &ErrInvalidOutput{Path: outputPath}
	}
	defer out.Close()

	// Encode and save the denoised image
	return jpeg.Encode(out, denoised, nil)
}

func medianFilter(img image.Image, x, y int) color.Color {
	var r, g, b []int
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			c := img.At(x+dx, y+dy)
			r1, g1, b1, _ := c.RGBA()
			r = append(r, int(r1>>8))
			g = append(g, int(g1>>8))
			b = append(b, int(b1>>8))
		}
	}
	sort.Ints(r)
	sort.Ints(g)
	sort.Ints(b)
	return color.RGBA{uint8(r[4]), uint8(g[4]), uint8(b[4]), 255}
}

// RotateImage rotates the input image by the specified angle in degrees.
// It takes the paths of the input and output files, and the rotation angle.
// Returns an error if the operation fails.
func RotateImage(inputPath string, outputPath string, angle float64) error {
	slog.Info("rotating image",
		"input", inputPath,
		"angle", angle)

	// Open the input file
	file, err := os.Open(inputPath)
	if err != nil {
		return &ErrInvalidInput{Path: inputPath}
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return &ErrProcessing{Op: "decode", Err: err}
	}

	// Convert angle to radians
	radians := angle * (math.Pi / 180)

	// Calculate new image size
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	newW, newH := rotatedSize(w, h, radians)

	// Create a new image with the rotated size
	rotated := image.NewRGBA(image.Rect(0, 0, newW, newH))

	// Rotate the image
	centerX, centerY := float64(w)/2, float64(h)/2
	newCenterX, newCenterY := float64(newW)/2, float64(newH)/2

	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			// Translate to origin
			xr := float64(x) - newCenterX
			yr := float64(y) - newCenterY

			// Rotate
			xr, yr = rotatePoint(xr, yr, -radians)

			// Translate back
			xr += centerX
			yr += centerY

			// If the point is within the original image, copy the color
			if xr >= 0 && xr < float64(w) && yr >= 0 && yr < float64(h) {
				rotated.Set(x, y, img.At(int(xr), int(yr)))
			}
		}
	}

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return &ErrInvalidOutput{Path: outputPath}
	}
	defer out.Close()

	// Encode and save the rotated image
	return jpeg.Encode(out, rotated, nil)
}

func rotatedSize(w, h int, angle float64) (int, int) {
	sin, cos := math.Abs(math.Sin(angle)), math.Abs(math.Cos(angle))
	newW := int(float64(w)*cos + float64(h)*sin)
	newH := int(float64(w)*sin + float64(h)*cos)
	return newW, newH
}

func rotatePoint(x, y float64, angle float64) (float64, float64) {
	return x*math.Cos(angle) - y*math.Sin(angle),
		x*math.Sin(angle) + y*math.Cos(angle)
}

// BinarizeImage applies Otsu's method to binarize the input image.
// It takes the paths of the input and output files.
// Returns an error if the operation fails.
func BinarizeImage(inputPath string, outputPath string) error {
	slog.Info("binarizing image", "input", inputPath)

	// Open the input file
	file, err := os.Open(inputPath)
	if err != nil {
		return &ErrInvalidInput{Path: inputPath}
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return &ErrProcessing{Op: "decode", Err: err}
	}

	// Convert to grayscale and calculate histogram
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	histogram := make([]int, 256)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldColor := img.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor).(color.Gray)
			grayImg.Set(x, y, grayColor)
			histogram[grayColor.Y]++
		}
	}

	// Calculate Otsu's threshold
	threshold := otsuThreshold(histogram, bounds.Dx()*bounds.Dy())

	// Apply threshold
	binarized := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if grayImg.GrayAt(x, y).Y > threshold {
				binarized.Set(x, y, color.White)
			} else {
				binarized.Set(x, y, color.Black)
			}
		}
	}

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return &ErrInvalidOutput{Path: outputPath}
	}
	defer out.Close()

	// Encode and save the binarized image
	return jpeg.Encode(out, binarized, nil)
}

func otsuThreshold(histogram []int, total int) uint8 {
	sum := 0
	for i := 0; i < 256; i++ {
		sum += i * histogram[i]
	}

	sumB, wB, wF := 0, 0, 0
	varMax, threshold := 0.0, 0

	for t := 0; t < 256; t++ {
		wB += histogram[t]
		if wB == 0 {
			continue
		}
		wF = total - wB
		if wF == 0 {
			break
		}
		sumB += t * histogram[t]
		mB := float64(sumB) / float64(wB)
		mF := float64(sum-sumB) / float64(wF)
		varBetween := float64(wB) * float64(wF) * (mB - mF) * (mB - mF)
		if varBetween > varMax {
			varMax = varBetween
			threshold = t
		}
	}

	return uint8(threshold)
}

// ConcatenateImagesVertically combines multiple images vertically into a single image.
// It takes a slice of input file paths and the output file path.
// Returns an error if the operation fails.
func ConcatenateImagesVertically(inputPaths []string, outputPath string) error {
	slog.Info("concatenating images vertically",
		"count", len(inputPaths),
		"output", outputPath)

	var images []image.Image
	var totalHeight int
	maxWidth := 0

	// Open and decode all input images
	for _, path := range inputPaths {
		file, err := os.Open(path)
		if err != nil {
			return &ErrInvalidInput{Path: path}
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return &ErrProcessing{Op: "decode", Err: err}
		}

		images = append(images, img)
		bounds := img.Bounds()

		// Find the maximum width
		if bounds.Dx() > maxWidth {
			maxWidth = bounds.Dx()
		}
	}

	// Resize images to match the maximum width while maintaining aspect ratios
	var resizedImages []image.Image
	totalHeight = 0
	for _, img := range images {
		bounds := img.Bounds()
		ratio := float64(bounds.Dx()) / float64(bounds.Dy())
		newHeight := int(float64(maxWidth) / ratio)

		resized := resize.Resize(uint(maxWidth), uint(newHeight), img, resize.Lanczos3)
		resizedImages = append(resizedImages, resized)
		totalHeight += newHeight
	}

	// Create a new image with the maximum width and total height
	concatenated := image.NewRGBA(image.Rect(0, 0, maxWidth, totalHeight))

	// Draw each resized image onto the concatenated image
	y := 0
	for _, img := range resizedImages {
		bounds := img.Bounds()
		draw.Draw(concatenated, image.Rect(0, y, maxWidth, y+bounds.Dy()), img, bounds.Min, draw.Src)
		y += bounds.Dy()
	}

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return &ErrInvalidOutput{Path: outputPath}
	}
	defer out.Close()

	// Encode and save the concatenated image
	if err := jpeg.Encode(out, concatenated, &jpeg.Options{Quality: cfg.JpegQuality}); err != nil {
		return &ErrProcessing{Op: "encode", Err: err}
	}

	return nil
}

// ConcatenateImagesHorizontally combines multiple images horizontally into a single image.
// It takes a slice of input file paths and the output file path.
// Returns an error if the operation fails.
func ConcatenateImagesHorizontally(inputPaths []string, outputPath string) error {
	slog.Info("concatenate image horizontally", "input", inputPaths)

	var images []image.Image
	var totalWidth int
	maxHeight := 0

	// Open and decode all input images
	for _, path := range inputPaths {
		file, err := os.Open(path)
		if err != nil {
			return &ErrInvalidInput{Path: path}
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return &ErrProcessing{Op: "decode", Err: err}
		}

		images = append(images, img)
		bounds := img.Bounds()

		// Find the maximum height
		if bounds.Dy() > maxHeight {
			maxHeight = bounds.Dy()
		}
	}

	// Resize images to match the maximum height while maintaining aspect ratios
	var resizedImages []image.Image
	totalWidth = 0
	for _, img := range images {
		bounds := img.Bounds()
		ratio := float64(bounds.Dx()) / float64(bounds.Dy())
		newWidth := int(float64(maxHeight) * ratio)

		resized := resize.Resize(uint(newWidth), uint(maxHeight), img, resize.Lanczos3)
		resizedImages = append(resizedImages, resized)
		totalWidth += newWidth
	}

	// Create a new image with the total width and maximum height
	concatenated := image.NewRGBA(image.Rect(0, 0, totalWidth, maxHeight))

	// Draw each resized image onto the concatenated image
	x := 0
	for _, img := range resizedImages {
		bounds := img.Bounds()
		draw.Draw(concatenated, image.Rect(x, 0, x+bounds.Dx(), maxHeight), img, bounds.Min, draw.Src)
		x += bounds.Dx()
	}

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return &ErrInvalidOutput{Path: outputPath}
	}
	defer out.Close()

	// Encode and save the concatenated image
	if err := jpeg.Encode(out, concatenated, &jpeg.Options{Quality: cfg.JpegQuality}); err != nil {
		return &ErrProcessing{Op: "encode", Err: err}
	}

	return nil
}

// GenerateTestImage creates a test image with random colored pixels.
// It takes the output file path and the desired width and height of the image.
// Returns an error if the operation fails.
func GenerateTestImage(outputPath string, width, height int) error {
	slog.Info("generating test image",
		"output", outputPath,
		"width", width,
		"height", height)

	// Create a new image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill the image with random colors
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(rand.Intn(256)),
				G: uint8(rand.Intn(256)),
				B: uint8(rand.Intn(256)),
				A: 255,
			})
		}
	}

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return &ErrInvalidOutput{Path: outputPath}
	}
	defer out.Close()

	// Encode and save the image as JPEG
	return jpeg.Encode(out, img, nil)
}

func BenchmarkResizeImage(b *testing.B) {
	inputPath := "../examples/input.jpg"
	outputPath := "../examples/output_resized.jpg"
	for i := 0; i < b.N; i++ {
		if err := ResizeImage(inputPath, outputPath, 800, 600); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDenoiseImage(b *testing.B) {
	inputPath := "../examples/input.jpg"
	outputPath := "../examples/output_denoised.jpg"
	for i := 0; i < b.N; i++ {
		if err := DenoiseImage(inputPath, outputPath); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRotateImage(b *testing.B) {
	inputPath := "../examples/input.jpg"
	outputPath := "../examples/output_rotated.jpg"
	for i := 0; i < b.N; i++ {
		if err := RotateImage(inputPath, outputPath, 90); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBinarizeImage(b *testing.B) {
	inputPath := "../examples/input.jpg"
	outputPath := "../examples/output_binarized.jpg"
	for i := 0; i < b.N; i++ {
		if err := BinarizeImage(inputPath, outputPath); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConcatenateImagesVertically(b *testing.B) {
	inputPaths := []string{"../examples/input1.jpg", "../examples/input2.jpg"}
	outputPath := "../examples/output_concat_vert.jpg"
	for i := 0; i < b.N; i++ {
		if err := ConcatenateImagesVertically(inputPaths, outputPath); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConcatenateImagesHorizontally(b *testing.B) {
	inputPaths := []string{"../examples/input1.jpg", "../examples/input2.jpg"}
	outputPath := "../examples/output_concat_horz.jpg"
	for i := 0; i < b.N; i++ {
		if err := ConcatenateImagesHorizontally(inputPaths, outputPath); err != nil {
			b.Fatal(err)
		}
	}
}

// DetectEdges applies Sobel edge detection to the input image.
// It takes the paths of the input and output files.
// Returns an error if the operation fails.
func DetectEdges(inputPath string, outputPath string) error {
	slog.Info("concatenating images horizontally",
		"count", len(inputPath),
		"output", outputPath)

	// Open the input file
	file, err := os.Open(inputPath)
	if err != nil {
		return &ErrInvalidInput{Path: inputPath}
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return &ErrProcessing{Op: "decode", Err: err}
	}

	// Convert to grayscale
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayImg.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}

	// Apply Sobel operator
	edgeImg := image.NewGray(bounds)
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			// Sobel kernels
			gx := -1*int(grayImg.GrayAt(x-1, y-1).Y) + 1*int(grayImg.GrayAt(x+1, y-1).Y) +
				-2*int(grayImg.GrayAt(x-1, y).Y) + 2*int(grayImg.GrayAt(x+1, y).Y) +
				-1*int(grayImg.GrayAt(x-1, y+1).Y) + 1*int(grayImg.GrayAt(x+1, y+1).Y)

			gy := -1*int(grayImg.GrayAt(x-1, y-1).Y) + 1*int(grayImg.GrayAt(x-1, y+1).Y) +
				-2*int(grayImg.GrayAt(x, y-1).Y) + 2*int(grayImg.GrayAt(x, y+1).Y) +
				-1*int(grayImg.GrayAt(x+1, y-1).Y) + 1*int(grayImg.GrayAt(x+1, y+1).Y)

			magnitude := uint8(math.Sqrt(float64(gx*gx + gy*gy)))
			edgeImg.Set(x, y, color.Gray{magnitude})
		}
	}

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return &ErrInvalidOutput{Path: outputPath}
	}
	defer out.Close()

	// Encode and save the edge-detected image
	if err := jpeg.Encode(out, edgeImg, nil); err != nil {
		return &ErrProcessing{Op: "encode", Err: err}
	}

	return nil
}

// ErrInvalidInput represents an error when the input file is invalid or cannot be opened.
type ErrInvalidInput struct {
	Path string
}

func (e *ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input file: %s", e.Path)
}

// ErrInvalidOutput represents an error when the output file cannot be created or written to.
type ErrInvalidOutput struct {
	Path string
}

func (e *ErrInvalidOutput) Error() string {
	return fmt.Sprintf("invalid output file: %s", e.Path)
}

// ErrProcessing represents a general error during image processing.
type ErrProcessing struct {
	Op  string
	Err error
}

func (e *ErrProcessing) Error() string {
	return fmt.Sprintf("error during %s: %v", e.Op, e.Err)
}

// ErrUnsupportedFormat represents an error when the image format is not supported.
type ErrUnsupportedFormat struct {
	Format string
}

func (e *ErrUnsupportedFormat) Error() string {
	return fmt.Sprintf("unsupported image format: %s", e.Format)
}
