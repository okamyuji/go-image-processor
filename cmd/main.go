package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/okamyuji/go-image-processor/config"
	processor "github.com/okamyuji/go-image-processor/pkg"
)

func init() {
	log.SetPrefix("go-image-processor: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	_ = config.GetConfig()
	// JSONフォーマットのロガーを設定
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

func printUsage() {
	fmt.Println("Usage: go-image-processor <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  resize -width <width> -height <height> <input> <output>")
	fmt.Println("  denoise <input> <output>")
	fmt.Println("  rotate -angle <angle> <input> <output>")
	fmt.Println("  autorotate <input> <output>")
	fmt.Println("  binarize <input> <output>")
	fmt.Println("  concatvert <output> <input1> <input2> [input3...]")
	fmt.Println("  concathorz <output> <input1> <input2> [input3...]")
	fmt.Println("  generatetest -width <width> -height <height> <output>")
	fmt.Println("\nUse 'go-image-processor <command> -h' for more information about a command.")
}

func handleError(err error) {
	switch e := err.(type) {
	case *processor.ErrInvalidInput:
		slog.Error("invalid input file",
			"path", e.Path)
	case *processor.ErrInvalidOutput:
		slog.Error("invalid output file",
			"path", e.Path)
	case *processor.ErrProcessing:
		slog.Error("processing error",
			"operation", e.Op,
			"error", e.Err)
	case *processor.ErrUnsupportedFormat:
		slog.Error("unsupported format",
			"format", e.Format)
	default:
		slog.Error("unexpected error",
			"error", err)
	}
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "resize":
		resizeCmd := flag.NewFlagSet("resize", flag.ExitOnError)
		width := resizeCmd.Int("width", 0, "Width to resize the image to")
		height := resizeCmd.Int("height", 0, "Height to resize the image to")
		if err := resizeCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Usage: go-image-processor resize <input> <output> -width <width> -height <height>")
			os.Exit(1)
		}
		if resizeCmd.NArg() < 2 || *width == 0 || *height == 0 {
			fmt.Println("Usage: go-image-processor resize <input> <output> -width <width> -height <height>")
			os.Exit(1)
		}
		if err := resizeCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Usage: go-image-processor resize <input> <output> -width <width> -height <height>")
			os.Exit(1)
		}
		if resizeCmd.NArg() < 2 || *width == 0 || *height == 0 {
			fmt.Println("Usage: go-image-processor resize <input> <output> -width <width> -height <height>")
			os.Exit(1)
		}

		err := processor.ResizeImage(resizeCmd.Arg(0), resizeCmd.Arg(1), uint(*width), uint(*height))
		if err != nil {
			handleError(err)
		}
		fmt.Println("Image resized successfully")

	case "denoise":
		denoiseCmd := flag.NewFlagSet("denoise", flag.ExitOnError)
		if err := denoiseCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Usage: go-image-processor denoise <input> <output>")
			os.Exit(1)
		}
		if denoiseCmd.NArg() < 2 {
			fmt.Println("Usage: go-image-processor denoise <input> <output>")
			os.Exit(1)
		}
		err := processor.DenoiseImage(denoiseCmd.Arg(0), denoiseCmd.Arg(1))
		if err != nil {
			handleError(err)
		}
		fmt.Println("Image denoised successfully")

	case "rotate":
		rotateCmd := flag.NewFlagSet("rotate", flag.ExitOnError)
		angle := rotateCmd.Float64("angle", 0, "Angle to rotate the image by")
		if err := rotateCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Usage: go-image-processor rotate <input> <output> -angle <angle>")
			os.Exit(1)
		}
		if rotateCmd.NArg() < 2 || *angle == 0 {
			fmt.Println("Usage: go-image-processor rotate <input> <output> -angle <angle>")
			os.Exit(1)
		}
		if rotateCmd.NArg() < 2 || *angle == 0 {
			fmt.Println("Usage: go-image-processor rotate <input> <output> -angle <angle>")
			os.Exit(1)
		}

		err := processor.RotateImage(rotateCmd.Arg(0), rotateCmd.Arg(1), *angle)
		if err != nil {
			handleError(err)
		}
		fmt.Println("Image rotated successfully")
	case "autorotate":
		autoRotateCmd := flag.NewFlagSet("autorotate", flag.ExitOnError)
		if err := autoRotateCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Usage: go-image-processor autorotate <input> <output>")
			os.Exit(1)
		}
		if autoRotateCmd.NArg() < 2 {
			fmt.Println("Usage: go-image-processor autorotate <input> <output>")
			os.Exit(1)
		}

		err := processor.AutoRotateImage(autoRotateCmd.Arg(0), autoRotateCmd.Arg(1))
		if err != nil {
			handleError(err)
		}
		fmt.Println("Image auto-rotated successfully")
	case "binarize":
		binarizeCmd := flag.NewFlagSet("binarize", flag.ExitOnError)
		if err := binarizeCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Usage: go-image-processor binarize <input> <output>")
			os.Exit(1)
		}

		if binarizeCmd.NArg() < 2 {
			fmt.Println("Usage: go-image-processor binarize <input> <output>")
			os.Exit(1)
		}

		err := processor.BinarizeImage(binarizeCmd.Arg(0), binarizeCmd.Arg(1))
		if err != nil {
			handleError(err)
		}
		fmt.Println("Image binarized successfully")

	case "concatvert":
		concatVertCmd := flag.NewFlagSet("concatvert", flag.ExitOnError)
		if err := concatVertCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Usage: go-image-processor concatvert <output> <input1> <input2> [input3...]")
			os.Exit(1)
		}

		if concatVertCmd.NArg() < 3 {
			fmt.Println("Usage: go-image-processor concatvert <output> <input1> <input2> [input3...]")
			os.Exit(1)
		}

		outputPath := concatVertCmd.Arg(0)
		inputPaths := concatVertCmd.Args()[1:]
		err := processor.ConcatenateImagesVertically(inputPaths, outputPath)
		if err != nil {
			handleError(err)
		}
		fmt.Println("Images concatenated vertically successfully")

	case "concathorz":
		concatHorzCmd := flag.NewFlagSet("concathorz", flag.ExitOnError)
		if err := concatHorzCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Usage: go-image-processor concathorz <output> <input1> <input2> [input3...]")
			os.Exit(1)
		}

		if concatHorzCmd.NArg() < 3 {
			fmt.Println("Usage: go-image-processor concathorz <output> <input1> <input2> [input3...]")
			os.Exit(1)
		}

		outputPath := concatHorzCmd.Arg(0)
		inputPaths := concatHorzCmd.Args()[1:]
		err := processor.ConcatenateImagesHorizontally(inputPaths, outputPath)
		if err != nil {
			handleError(err)
		}
		fmt.Println("Images concatenated horizontally successfully")

	case "generatetest":
		generateTestCmd := flag.NewFlagSet("generatetest", flag.ExitOnError)
		width := generateTestCmd.Int("width", 100, "Width of the test image")
		height := generateTestCmd.Int("height", 100, "Height of the test image")
		if err := generateTestCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Usage: go-image-processor generatetest <output> -width <width> -height <height>")
			os.Exit(1)
		}
		if generateTestCmd.NArg() < 1 || *width == 0 || *height == 0 {
			fmt.Println("Usage: go-image-processor generatetest <output> -width <width> -height <height>")
			os.Exit(1)
		}

		if generateTestCmd.NArg() < 1 {
			fmt.Println("Usage: go-image-processor generatetest <output> -width <width> -height <height>")
			os.Exit(1)
		}

		outputPath := generateTestCmd.Arg(0)
		err := processor.GenerateTestImage(outputPath, *width, *height)
		if err != nil {
			handleError(err)
		}
		fmt.Println("Test image generated successfully")
	case "edges":
		edgesCmd := flag.NewFlagSet("edges", flag.ExitOnError)
		if err := edgesCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println("Usage: go-image-processor edges <input> <output>")
			os.Exit(1)
		}

		if edgesCmd.NArg() < 2 {
			fmt.Println("Usage: go-image-processor edges <input> <output>")
			os.Exit(1)
		}

		err := processor.DetectEdges(edgesCmd.Arg(0), edgesCmd.Arg(1))
		if err != nil {
			handleError(err)
		}
		fmt.Println("Edge detection completed successfully")
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}
