// TODO: таймеры, аргументы командной строки, рефактор мейна, чтение из файла

package main

import (
	"bufio"
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
)

const typedColor = termbox.ColorGreen
const toBeTypedColor = termbox.ColorWhite
const errorColor = termbox.ColorRed | termbox.AttrBold

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func makeSubstrs(str string, strpos, width, leftMargin int) (typed string, typedMargin int, toBeTyped string) {
	if strpos-leftMargin < 0 {
		typed = str[0:strpos]
		typedMargin = leftMargin - strpos
	} else {
		typed = str[strpos-leftMargin : strpos]
	}
	toBeTyped = str[strpos:min(strpos+width-leftMargin, len(str))]
	return
}

func printStr(x, y int, str string) {
	for i := 0; i < len(str); i++ {
		termbox.SetChar(x+i, y, rune(str[i]))
	}
}

func colorSegment(x, y, len int, color termbox.Attribute) {
	for i := 0; i < len; i++ {
		termbox.SetFg(x+i, y, color)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	err := termbox.Init()
	if err != nil {
		fmt.Println("Couldn't initialize \"termbox\" library.")
		return
	}

	str := "The quick brown fox jumps over the lazy dog"
	strpos := 0

	typo := false

	for strpos != len(str) {
		width, height := termbox.Size()
		if width < 24 || height < 5 {
			fmt.Println("The window is too small!")
			return
		}

		err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		if err != nil {
			return
		}

		leftMargin := width / 4
		typed, typedMargin, toBeTyped := makeSubstrs(str, strpos, width, leftMargin)
		printStr(typedMargin, height/2, typed)
		colorSegment(0, height/2, leftMargin, typedColor)
		printStr(leftMargin, height/2, toBeTyped)
		if typo {
			colorSegment(leftMargin, height/2, width-leftMargin, errorColor)
		} else {
			colorSegment(leftMargin, height/2, width-leftMargin, toBeTypedColor)
		}
		termbox.SetCursor(leftMargin, height/2)

		termbox.Flush()

		ch, _, _ := reader.ReadRune()
		if ch != rune(str[strpos]) {
			typo = true
			continue
		}
		typo = false
		strpos++
	}
}
