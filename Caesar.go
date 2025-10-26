package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Page struct {
	Input  string
	Result string
	Shift  int
	Lang   string
	Brute  []string
}

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

func viewHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Lang: "ru"}

	if r.Method == http.MethodPost {
		text := r.FormValue("text")
		shiftStr := r.FormValue("shift")
		mode := r.FormValue("mode")
		lang := r.FormValue("lang")

		p.Input = text
		p.Lang = lang

		shift, _ := strconv.Atoi(shiftStr)
		p.Shift = shift

		switch mode {
		case "encrypt":
			p.Result = caesar(text, shift, lang, false)
		case "decrypt":
			p.Result = caesar(text, shift, lang, true)
		case "brute":
			p.Brute = bruteForce(text, lang)
		}
	}

	t, _ := template.ParseFiles("main.html")

	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/", viewHandler)
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
