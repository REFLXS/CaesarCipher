package main

import (
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var ruLower = []rune("абвгдеёжзийклмнопрстуфхцчшщъыьэюя")
var ruUpper = []rune("АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ")
var enLower = []rune("abcdefghijklmnopqrstuvwxyz")
var enUpper = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func shiftRune(r rune, shift int, lower, upper []rune) rune {
	for i, c := range lower {
		if r == c {
			n := (i + shift) % len(lower)
			if n < 0 {
				n += len(lower)
			}
			return lower[n]
		}
	}
	for i, c := range upper {
		if r == c {
			n := (i + shift) % len(upper)
			if n < 0 {
				n += len(upper)
			}
			return upper[n]
		}
	}
	return r
}

func caesar(text string, shift int, lang string, decrypt bool) string {
	if decrypt {
		shift = -shift
	}
	var sb strings.Builder
	for _, r := range text {
		if lang == "ru" {
			sb.WriteRune(shiftRune(r, shift, ruLower, ruUpper))
		} else {
			sb.WriteRune(shiftRune(r, shift, enLower, enUpper))
		}
	}
	return sb.String()
}

func bruteForce(text, lang string) []string {
	var res []string
	max := len(enLower)
	if lang == "ru" {
		max = len(ruLower)
	}
	for i := 1; i < max; i++ {
		res = append(res, fmt.Sprintf("%2d. %s", i, caesar(text, i, lang, true)))
	}
	return res
}

func main() {
	a := app.New()
	w := a.NewWindow("Шифр Цезаря")
	w.Resize(fyne.NewSize(1000, 700))

	input := widget.NewMultiLineEntry()
	input.Wrapping = fyne.TextWrapWord

	output := widget.NewMultiLineEntry()
	output.Wrapping = fyne.TextWrapWord

	shiftEntry := widget.NewEntry()

	langSelect := widget.NewSelect([]string{"ru", "en"}, func(s string) {})
	langSelect.SetSelected("ru")

	modeSelect := widget.NewSelect([]string{"encrypt", "decrypt", "brute"}, func(s string) {})
	modeSelect.SetSelected("encrypt")

	runBtn := widget.NewButton("Выполнить", func() {
		text := input.Text
		lang := langSelect.Selected
		mode := modeSelect.Selected

		var result string

		shift := 0
		if shiftEntry.Text != "" {
			shiftVal, _ := strconv.Atoi(shiftEntry.Text)
			shift = shiftVal
		}

		switch mode {
		case "encrypt":
			result = caesar(text, shift, lang, false)

		case "decrypt":
			result = caesar(text, shift, lang, true)

		case "brute":
			lines := bruteForce(text, lang)
			result = strings.Join(lines, "\n")
		}

		output.SetText(result)
	})

	inputScroll := container.NewScroll(input)
	inputScroll.SetMinSize(fyne.NewSize(400, 400))

	outputScroll := container.NewScroll(output)
	outputScroll.SetMinSize(fyne.NewSize(400, 400))

	left := container.NewVBox(
		widget.NewLabel("Ввод"),
		inputScroll,
		container.NewHBox(
			widget.NewLabel("Язык"),
			langSelect,
			widget.NewLabel("Режим"),
			modeSelect,
		),
		container.NewHBox(
			widget.NewLabel("Сдвиг"),
			shiftEntry,
		),
		runBtn,
	)

	right := container.NewVBox(
		widget.NewLabel("Результат"),
		outputScroll,
	)

	content := container.NewAdaptiveGrid(2, left, right)

	w.SetContent(content)
	w.ShowAndRun()
}
