package main

import (
	"bufio"
	"fmt"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"os"
	"time"
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

func paintSegment(x, y, len int, color termbox.Attribute) {
	for i := 0; i < len; i++ {
		termbox.SetFg(x+i, y, color)
	}
}

func printAndPaint(width, height int, str string, strpos int, typo bool) {
	leftMargin := width / 4
	typed, typedMargin, toBeTyped := makeSubstrs(str, strpos, width, leftMargin)
	printStr(typedMargin, height/2, typed)
	paintSegment(0, height/2, leftMargin, typedColor)
	printStr(leftMargin, height/2, toBeTyped)
	if typo {
		paintSegment(leftMargin, height/2, width-leftMargin, errorColor)
	} else {
		paintSegment(leftMargin, height/2, width-leftMargin, toBeTypedColor)
	}
	termbox.SetCursor(leftMargin, height/2)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: \"./typingapp file\", where \"file\" contains the text you want to type.")
		return
	}

	err := termbox.Init()
	if err != nil {
		fmt.Println("Couldn't initialize \"termbox\" library.")
		return
	}

	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Couldn't open file \"", os.Args[1], "\".")
		return
	}

	str := string(bytes)
	strpos := 0

	if len(str) == 0 {
		fmt.Println("The file is empty!")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	typo := false
	typosNum := 0
	start := time.Now()

	for strpos != len(str) {
		width, height := termbox.Size()
		if width < 24 || height < 5 {
			fmt.Println("The window is too small!")
			return
		}

		err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
		if err != nil {
			fmt.Println("Couldn't clear the console.")
			return
		}

		printAndPaint(width, height, str, strpos, typo)

		err = termbox.Flush()
		if err != nil {
			fmt.Println("Couldn't flush the buffer.")
			return
		}

		ch, _, _ := reader.ReadRune()
		if ch != rune(str[strpos]) {
			typo = true
			typosNum++
			continue
		}

		typo = false
		strpos++
	}

	end := time.Now()

	fmt.Printf("Your average speed is: %.2f characters per second.\n", float64(len(str))/end.Sub(start).Seconds())
	fmt.Println("You mamde ", typosNum, " typos.")
}
