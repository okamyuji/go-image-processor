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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
