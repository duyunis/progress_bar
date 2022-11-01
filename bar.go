package progress_bar

import (
	"fmt"
	"runtime"
	"strings"
)

type Bar struct {
	percent     int64
	curr        int64
	total       int64
	currentRate int64
	graph       string
	backGraph   string
	graphTotal  int
	describe    string
	color       barColor
}

type barColor struct {
	describeColor  color
	graphColor     color
	graphBackColor backcolor
	percentColor   color
	ratioColor     color
}

type color int
type backcolor int

var (
	Black      color
	White      color
	Red        color
	Blue       color
	Green      color
	Yellow     color
	Purple     color
	BackWhite  backcolor
	BackRed    backcolor
	BackBlue   backcolor
	BackGreen  backcolor
	BackYellow backcolor
	BackBlack  backcolor
	BackPurple backcolor
)

func init() {
	if runtime.GOOS == "windows" {
		//
	} else {
		Black = 30
		BackBlack = 40
		White = 37
		BackWhite = 47
		Red = 31
		BackRed = 41
		Blue = 34
		BackBlue = 44
		Green = 32
		BackGreen = 42
		Yellow = 33
		BackYellow = 43
		Purple = 35
		BackPurple = 45
	}
}

func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.curr) / float32(bar.total) * 100)
}

func (bar *Bar) PrintBar(curr int64) {
	bar.SetCurrValue(curr)
	printBar := "\r" + bar.getDescribe()
	printBar += bar.getProgressBarString()
	printBar += bar.getPercentPrintString()
	printBar += "     "
	printBar += bar.RatioPrintString()
	fmt.Print(printBar)
}

func (bar *Bar) getDescribe() string {
	if bar.describe == "" {
		return ""
	}
	return fmt.Sprintf(" %c[%vm%v%c[0m", 0x1B, bar.color.describeColor, bar.describe, 0x1B)
}

func (bar *Bar) getProgressBarString() string {
	full := bar.getFullBackGraph()
	printBar := "[" + strings.Replace(full, bar.backGraph, bar.graph, bar.currentBarGraphNumber()) + "]"
	if runtime.GOOS == "windows" {
		return printBar
	} else {
		return fmt.Sprintf("%c[%v;%vm%s%c[0m", 0x1B, bar.color.graphColor, bar.color.graphBackColor, printBar, 0x1B)
	}
}

func (bar *Bar) getPercentPrintString() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf(" %v%%", bar.currentRate)
	}
	return fmt.Sprintf(" %c[%vm%v%%%c[0m", 0x1B, bar.color.percentColor, bar.currentRate, 0x1B)
}

func (bar *Bar) RatioPrintString() string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf(" %v/%v", bar.curr, bar.total)
	}
	return fmt.Sprintf(" %c[%vm\t%v/%v%c[0m", 0x1B, bar.color.ratioColor, bar.curr, bar.total, 0x1B)
}

func (bar *Bar) getFullBackGraph() string {
	res := ""
	for i := 0; i < bar.graphTotal; i++ {
		res += bar.backGraph
	}
	return res
}

func (bar *Bar) countCurrentRate() {
	if bar.curr == 0 {
		bar.currentRate = 0
	} else {
		bar.currentRate = bar.curr * 100 / bar.total
	}
}

func (bar *Bar) currentBarGraphNumber() int {
	if bar.currentRate == 100 {
		return bar.graphTotal
	} else {
		return int(float64(bar.currentRate) * (float64(bar.graphTotal) / float64(100)))
	}
}

func (bar *Bar) Finish() {
	fmt.Println()
}

func (bar *Bar) SetCurrValue(curr int64) {
	bar.curr = curr
	bar.countCurrentRate()
}

func (bar *Bar) SetGraph(graph string) {
	bar.graph = graph
	gl := len(graph)
	if gl > 1 {
		for i := 1; i < gl; i++ {
			bar.backGraph += " "
		}
	}
}

func (bar *Bar) SetGraphTotal(graphTotal int) {
	bar.graphTotal = graphTotal
}

func (bar *Bar) SetGraphColor(c color) {
	bar.color.graphColor = c
}

func (bar *Bar) SetGraphBackColor(b backcolor) {
	bar.color.graphBackColor = b
}

func (bar *Bar) SetDescribeColor(c color) {
	bar.color.describeColor = c
}

func (bar *Bar) SetRatioColor(c color) {
	bar.color.ratioColor = c
}

func (bar *Bar) SetPercentColor(c color) {
	bar.color.percentColor = c
}

func (bar *Bar) SetDescribe(describe string) {
	bar.describe = describe
}

func (bar *Bar) setBackGraph(bg string) {
	bar.backGraph = bg
}

func NewBar(total int64) *Bar {
	bar := &Bar{
		total:      total,
		graph:      "█",
		backGraph:  " ",
		graphTotal: 50,
	}
	bar.countCurrentRate()
	return bar
}
