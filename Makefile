BINARY_NAME=go-image-processor
GUI_BINARY_NAME=go-image-processor-gui

.PHONY: all build test clean run build-gui run-gui
.PHONY: ensure-examples-dir generate-test-inputs
.PHONY: resize-example denoise-example rotate-example binarize-example
.PHONY: concatvert-example concathorz-example generatetest-example
.PHONY: edges-example autorotate-example benchmark

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
	rm -rf examples

run:
	./${BINARY_NAME}

run-gui:
	./${GUI_BINARY_NAME}

# Ensure examples directory exists
ensure-examples-dir:
	@mkdir -p examples

# Generate test input images
generate-test-inputs: ensure-examples-dir build
	@echo "=== Generating Test Input Images ==="
	./${BINARY_NAME} generatetest examples width 200 height 200

# Example commands
resize-example: ensure-examples-dir generate-test-inputs
	@echo "=== Resizing Image ==="
	@echo "Command: resize -width 800 -height 600 examples/rotation_test.jpg examples/output_resized.jpg"
	./${BINARY_NAME} resize -width 800 -height 600 examples/rotation_test.jpg examples/output_resized.jpg
	@echo "Output saved to examples/output_resized.jpg"
	@ls -lh examples/output_resized.jpg

denoise-example: ensure-examples-dir generate-test-inputs
	@echo "=== Denoising Image ==="
	@echo "Command: denoise examples/noise_test.jpg examples/output_denoised.jpg"
	./${BINARY_NAME} denoise examples/noise_test.jpg examples/output_denoised.jpg
	@echo "Output saved to examples/output_denoised.jpg"
	@ls -lh examples/output_denoised.jpg

rotate-example: ensure-examples-dir generate-test-inputs
	@echo "=== Rotating Image ==="
	@echo "Command: rotate -angle 90 examples/rotation_test.jpg examples/output_rotated.jpg"
	./${BINARY_NAME} rotate -angle 90 examples/rotation_test.jpg examples/output_rotated.jpg
	@echo "Output saved to examples/output_rotated.jpg"
	@ls -lh examples/output_rotated.jpg

binarize-example: ensure-examples-dir generate-test-inputs
	@echo "=== Binarizing Image ==="
	@echo "Command: binarize examples/binary_test.jpg examples/output_binarized.jpg"
	./${BINARY_NAME} binarize examples/binary_test.jpg examples/output_binarized.jpg
	@echo "Output saved to examples/output_binarized.jpg"
	@ls -lh examples/output_binarized.jpg

concatvert-example: ensure-examples-dir generate-test-inputs
	@echo "=== Concatenating Images Vertically ==="
	@echo "Command: concatvert examples/output_concat_vert.jpg examples/concat_test_1.jpg examples/concat_test_2.jpg"
	./${BINARY_NAME} concatvert examples/output_concat_vert.jpg examples/concat_test_1.jpg examples/concat_test_2.jpg
	@echo "Output saved to examples/output_concat_vert.jpg"
	@ls -lh examples/output_concat_vert.jpg

concathorz-example: ensure-examples-dir generate-test-inputs
	@echo "=== Concatenating Images Horizontally ==="
	@echo "Command: concathorz examples/output_concat_horz.jpg examples/concat_test_1.jpg examples/concat_test_2.jpg"
	./${BINARY_NAME} concathorz examples/output_concat_horz.jpg examples/concat_test_1.jpg examples/concat_test_2.jpg
	@echo "Output saved to examples/output_concat_horz.jpg"
	@ls -lh examples/output_concat_horz.jpg

generatetest-example: build ensure-examples-dir
	@echo "=== Generating Test Images ==="
	@echo "Command: generatetest -width 200 -height 200 examples"
	./${BINARY_NAME} generatetest -width 200 -height 200 examples
	@echo "\n=== Generated Test Images ==="
	@find examples -type f -name "*.jpg" | sort | while read file; do \
		echo "$$(basename $$file) - $$(stat -f %z $$file) bytes"; \
	done

edges-example: ensure-examples-dir generate-test-inputs
	@echo "=== Detecting Edges ==="
	@echo "Command: edges examples/gradient_test.jpg examples/output_edges.jpg"
	./${BINARY_NAME} edges examples/gradient_test.jpg examples/output_edges.jpg
	@echo "Output saved to examples/output_edges.jpg"
	@ls -lh examples/output_edges.jpg

autorotate-example: ensure-examples-dir generate-test-inputs
	@echo "=== Auto-rotating Image ==="
	@echo "Command: autorotate examples/skew_test_1.jpg examples/output_autorotate.jpg"
	./${BINARY_NAME} autorotate examples/skew_test_1.jpg examples/output_autorotate.jpg
	@echo "Output saved to examples/output_autorotate.jpg"
	@ls -lh examples/output_autorotate.jpg

benchmark:
	@echo "=== Running Benchmarks ==="
	go test -bench=. ./...