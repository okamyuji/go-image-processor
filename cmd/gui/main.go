package main

import (
	"fmt"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Image Processor")

	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder("Input file path")

	outputEntry := widget.NewEntry()
	outputEntry.SetPlaceHolder("Output file path")

	widthEntry := widget.NewEntry()
	widthEntry.SetPlaceHolder("Width")

	heightEntry := widget.NewEntry()
	heightEntry.SetPlaceHolder("Height")

	angleEntry := widget.NewEntry()
	angleEntry.SetPlaceHolder("Angle")

	operationSelect := widget.NewSelect([]string{"resize", "rotate", "denoise", "binarize", "edges"}, func(value string) {})

	processButton := widget.NewButton("Process", func() {
		operation := operationSelect.Selected
		input := inputEntry.Text
		output := outputEntry.Text

		if operation == "" || input == "" || output == "" {
			dialog.ShowError(fmt.Errorf("please fill in all required fields"), w)
			return
		}

		var cmd *exec.Cmd

		switch operation {
		case "resize":
			width := widthEntry.Text
			height := heightEntry.Text
			if width == "" || height == "" {
				dialog.ShowError(fmt.Errorf("please specify width and height for resize operation"), w)
				return
			}
			cmd = exec.Command("go-image-processor", "resize", input, output, "-width", width, "-height", height)
		case "rotate":
			angle := angleEntry.Text
			if angle == "" {
				dialog.ShowError(fmt.Errorf("please specify angle for rotate operation"), w)
				return
			}
			cmd = exec.Command("go-image-processor", "rotate", input, output, "-angle", angle)
		case "denoise":
			cmd = exec.Command("go-image-processor", "denoise", input, output)
		case "binarize":
			cmd = exec.Command("go-image-processor", "binarize", input, output)
		case "edges":
			cmd = exec.Command("go-image-processor", "edges", input, output)
		}

		cmdOutput, err := cmd.CombinedOutput()
		if err != nil {
			dialog.ShowError(fmt.Errorf("error processing image: %v\n%s", err, string(cmdOutput)), w)
			return
		}

		dialog.ShowInformation("Success", "Image processed successfully", w)
	})

	content := container.NewVBox(
		widget.NewLabel("Select operation:"),
		operationSelect,
		widget.NewLabel("Input file:"),
		inputEntry,
		widget.NewLabel("Output file:"),
		outputEntry,
		widget.NewLabel("Width (for resize):"),
		widthEntry,
		widget.NewLabel("Height (for resize):"),
		heightEntry,
		widget.NewLabel("Angle (for rotate):"),
		angleEntry,
		processButton,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(300, 400))
	w.ShowAndRun()
}
