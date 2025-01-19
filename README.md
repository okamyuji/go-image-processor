# Go Image Processor

Go Image Processor is a command-line tool for performing various image processing operations on JPEG images.

## Features

- Resize images
- Denoise images
- Rotate images
- Binarize images using Otsu's method
- Concatenate images vertically or horizontally
- Generate test images
- Detect edges using Sobel operator
- Configuration file for default settings
- Graphical User Interface for easier use

## Installation

1. Ensure you have Go installed on your system (version 1.16 or later).
2. Clone this repository

    ```shell
    git clone https://github.com/okamyuji/go-image-processor.git
    ```

3. Navigate to the project directory

    ```shell
    cd go-image-processor
    ```

4. Build the project

    ```shell
    make all
    ```

## Usage

### Command Line Interface

The general syntax for using the CLI tool is:

```shell
./go-image-processor <command> [arguments]
```

### Graphical User Interface

A simple graphical user interface (GUI) is available for easier use of the image processing tool. To build and run the GUI:

```shell
make build-gui
make run-gui
```

The GUI provides a user-friendly interface for selecting operations, inputting file paths, and setting parameters for image processing tasks.

### Available commands

1. Resize an image

    ```shell
    ./go-image-processor resize <input> <output> -width <width> -height <height>
    ```

2. Denoise an image

    ```shell
    ./go-image-processor denoise <input> <output>
    ```

3. Rotate an image

    ```shell
    ./go-image-processor rotate <input> <output> -angle <angle>
    ```

4. Binarize an image

    ```shell
    ./go-image-processor binarize <input> <output>
    ```

5. Concatenate images vertically

    ```shell
    ./go-image-processor concatvert <output> <input1> <input2> [input3...]
    ```

6. Concatenate images horizontally

    ```shell
    ./go-image-processor concathorz <output> <input1> <input2> [input3...]
    ```

7. Generate a test image

    ```shell
    ./go-image-processor generatetest <output> -width <width> -height <height>
    ```

8. Detect edges in an image:

    ```shell
    ./go-image-processor edges <input> <output>
    ```

For more information about a specific command, use

```shell
./go-image-processor <command> -h
```

## Examples

1. Resize an image to 800x600:

    ```shell
    ./go-image-processor resize input.jpg output.jpg -width 800 -height 600
    ```

2. Rotate an image by 90 degrees:

    ```shell
    ./go-image-processor rotate input.jpg output.jpg -angle 90
    ```

3. Concatenate three images vertically:

    ```shell
    ./go-image-processor concatvert output.jpg input1.jpg input2.jpg input3.jpg
    ```

4. Detect edges in an image:

    ```shell
    ./go-image-processor edges input.jpg output_edges.jpg
    ```

## Configuration

The application uses a `config.yaml` file for default settings. You can modify this file to change the default values for various operations. The configuration file should be placed in the root directory of the project.

Example `config.yaml`:

```yaml
default_width: 800
default_height: 600
default_angle: 90
jpeg_quality: 75
```

If the configuration file is not found, the application will use built-in default values.

## Quick Start with Makefile

This project includes a Makefile for easy building, testing, and running example commands.

1. Build the project

    ```shell
    make build
    ```

2. Run tests:

    ```shell
    make test
    ```

3. Clean build artifacts

    ```shell
    make clean
    ```

4. Run example commands:

    ```shell
    make resize-example
    make denoise-example
    make rotate-example
    make binarize-example
    make concatvert-example
    make concathorz-example
    make generatetest-example
    ```

5. Run benchmarks:

```shell
make benchmark
```

These commands will process the example images in the `examples` directory.

## Continuous Integration

This project uses GitHub Actions for continuous integration. On every push and pull request to the main branch, the project is built and all tests are run automatically.

## Benchmarks

To run benchmarks for the image processing functions, use

```shell
make benchmark
```

This will run performance tests on all the main functions, giving you an idea of their execution time and efficiency.

## Documentation

The code is documented using godoc. To view the documentation, run

```shell
godoc -http=:6060
```

Then open your browser and navigate to `http://localhost:6060/pkg/github.com/okamyuji/go-image-processor/pkg/processor/`

## Logging

The application uses Go's built-in logging package to log information about the operations being performed. Logs are printed to stderr by default.

To redirect logs to a file, you can run the application like this:

```shell
./go-image-processor <command> [arguments] 2> logfile.txt
```

This will send all log output to `logfile.txt`.

## How Each Image Processing Feature Works

### Resize Image

Have you ever needed to make a photo smaller or bigger? That's what our resize feature does!

1. It takes your original image
2. Keeps the same shape (like a rectangle stays a rectangle)
3. Makes it bigger or smaller while keeping everything looking natural
4. Saves the new sized image

Example: Making a large 1000x1000 photo smaller to fit on your screen at 500x500.

### Denoise Image (Remove Noise)

Think of noise as tiny unwanted dots in your photo, like static on an old TV.

1. The program looks at each part of the image
2. For each spot, it checks the colors around it
3. If it finds a dot that looks out of place, it smooths it out
4. The result is a cleaner, clearer image

Example: Making a grainy dark photo look smoother and clearer.

### Rotate Image

Just like turning a photo in your hands, this feature rotates your image.

1. You tell it how many degrees to turn (like 90° for a quarter turn)
2. It carefully moves each part of the image to its new position
3. Makes sure nothing gets cut off
4. Saves the turned image

Example: Turning a sideways photo to make it upright.

### Binarize Image (Black and White Conversion)

This turns your image into just black and white - no gray areas!

1. Looks at how bright each part of the image is
2. Decides if each spot should be black or white
3. Uses a smart method (called Otsu) to make the best choice
4. Creates a clear black and white version

Example: Making a color photo look like an old newspaper picture.

### Concatenate Images

This feature can join images together like puzzle pieces!

Vertical Concatenation:

1. Takes two or more images
2. Stacks them on top of each other
3. Makes sure they line up perfectly
4. Creates one tall image

Horizontal Concatenation:

1. Takes two or more images
2. Places them side by side
3. Lines them up evenly
4. Creates one wide image

Example: Joining two holiday photos to make a panorama.

### Edge Detection

This feature finds and highlights the outlines in your image!

1. Looks for places where colors change suddenly
2. Marks these changes as edges
3. Makes the edges stand out
4. Creates an image showing just the outlines

Example: Making a sketch-like version of a photo, showing just the main shapes.

### Auto-rotate Image (Skew Correction)

This feature automatically detects and corrects tilted images!

1. Analyzes the image to find strong lines or text
2. Calculates how much the image is tilted
3. Rotates the image to make it straight
4. Saves the corrected image

Example: Fixing a scanned document that was placed slightly crooked.

### Test Image Generation

Need sample images to practice with? This feature creates them!

1. Makes different types of test images
2. Creates patterns that are perfect for testing
3. Lets you choose the size
4. Saves them as regular image files

Example: Creating a checkerboard pattern to test image processing.

---
Remember: All these features keep your original image safe and create new files with the changes. It's like having a photo copy machine that can do magic tricks with your pictures! 🪄📸

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
