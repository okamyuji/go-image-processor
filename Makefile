BINARY_NAME=go-image-processor
GUI_BINARY_NAME=go-image-processor-gui

.PHONY: all build test clean run build-gui run-gui

all: build build-gui

build:
	go build -o ${BINARY_NAME} ./cmd

build-gui:
	go build -o ${GUI_BINARY_NAME} ./cmd/gui

test:
	go test -v ./...

clean:
	go clean
	rm -f ${BINARY_NAME}
	rm -f ${GUI_BINARY_NAME}

run:
	./${BINARY_NAME}

run-gui:
	./${GUI_BINARY_NAME}

# Example commands
resize-example:
	./${BINARY_NAME} resize examples/input.jpg examples/output_resized.jpg -width 800 -height 600

denoise-example:
	./${BINARY_NAME} denoise examples/input.jpg examples/output_denoised.jpg

rotate-example:
	./${BINARY_NAME} rotate examples/input.jpg examples/output_rotated.jpg -angle 90

binarize-example:
	./${BINARY_NAME} binarize examples/input.jpg examples/output_binarized.jpg

concatvert-example:
	./${BINARY_NAME} concatvert examples/output_concat_vert.jpg examples/input1.jpg examples/input2.jpg

concathorz-example:
	./${BINARY_NAME} concathorz examples/output_concat_horz.jpg examples/input1.jpg examples/input2.jpg

generatetest-example:
	./${BINARY_NAME} generatetest examples/test_image.jpg -width 200 -height 200

edges-example:
	./${BINARY_NAME} edges examples/input.jpg examples/output_edges.jpg

benchmark:
	go test -bench=. ./...