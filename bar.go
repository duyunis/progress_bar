package progress_bar

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type Bar struct {
	lock           sync.RWMutex
	current        int64
	total          int64
	lastCount      int64
	currentPercent int64
	currentRate    string
	graph          string
	backGraph      string
	graphTotal     int
	describe       string
	duration       string
	showBytes      bool
	showPercent    bool
	showDuration   bool
	finish         bool
	color          *BarColor
	quit           chan bool
}

type Options struct {
	Graph        string
	Describe     string
	IsBytes      bool
	ShowPercent  bool
	ShowDuration bool
	Color        *BarColor
}

type BarColor struct {
	DescribeColor  color
	GraphColor     color
	GraphBackColor backcolor
	PercentColor   color
	RatioColor     color
	DurationColor  color
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

func (bar *Bar) Add(curr int64) {
	bar.lock.Lock()
	defer bar.lock.Unlock()
	if bar.finish {
		return
	}
	bar.SetCurrValue(curr)
	printBar := "\r" + bar.getDescribe()
	if bar.showPercent {
		printBar += bar.getPercentPrintString()
	}
	printBar += bar.getProgressBarString()
	printBar += bar.getRatioPrintString()
	if bar.showDuration {
		printBar += bar.getDuration()
	}
	fmt.Print(printBar + "  ")
	fmt.Printf("\033[%dA\033[%dB", 1, 0)
}

func (bar *Bar) getDescribe() string {
	if bar.describe == "" {
		return ""
	}
	return fmt.Sprintf(" %c[%vm%v%c[0m", 0x1B, bar.color.DescribeColor, bar.describe, 0x1B)
}

func (bar *Bar) getProgressBarString() string {
	full := bar.getFullBackGraph()
	printBar := "[" + strings.Replace(full, bar.backGraph, bar.graph, bar.currentBarGraphNumber()) + "]"
	if runtime.GOOS == "windows" || (bar.color.GraphColor == 0 && bar.color.GraphBackColor == 0) {
		return printBar
	} else {
		return fmt.Sprintf("%c[%v;%vm%s%c[0m", 0x1B, bar.color.GraphColor, bar.color.GraphBackColor, printBar, 0x1B)
	}
}

func (bar *Bar) getPercentPrintString() string {
	if runtime.GOOS == "windows" || bar.color.PercentColor == 0 {
		return fmt.Sprintf(" %v%%", bar.currentPercent)
	}
	return fmt.Sprintf(" %c[%vm%v%%%c[0m", 0x1B, bar.color.PercentColor, bar.currentPercent, 0x1B)
}

func (bar *Bar) getDuration() string {
	if runtime.GOOS == "windows" || bar.color.DurationColor == 0 {
		return fmt.Sprintf(" [%s]", bar.duration)
	}
	return fmt.Sprintf(" %c[%vm[%v]%c[0m", 0x1B, bar.color.DurationColor, bar.duration, 0x1B)
}

func (bar *Bar) getRatioPrintString() string {
	bar.startCalcRate()
	var curr any
	var total any
	if bar.showBytes {
		curr = formatBytes(float64(bar.current))
		total = formatBytes(float64(bar.total))
		if runtime.GOOS == "windows" || bar.color.RatioColor == 0 {
			return fmt.Sprintf(" (%v/%v, %s)", curr, total, bar.currentRate)
		} else {
			return fmt.Sprintf(" %c[%vm\t(%v/%v, %s)%c[0m", 0x1B, bar.color.RatioColor, curr, total, bar.currentRate, 0x1B)
		}
	} else {
		curr = bar.current
		total = bar.total
		if runtime.GOOS == "windows" || bar.color.RatioColor == 0 {
			return fmt.Sprintf(" %v/%v", curr, total)
		} else {
			return fmt.Sprintf(" %c[%vm\t%v/%v%c[0m", 0x1B, bar.color.RatioColor, curr, total, 0x1B)
		}
	}
}

func (bar *Bar) getFullBackGraph() string {
	res := ""
	for i := 0; i < bar.graphTotal; i++ {
		res += bar.backGraph
	}
	return res
}

func (bar *Bar) countCurrentPercent() {
	if bar.current == 0 {
		bar.currentPercent = 0
	} else {
		bar.currentPercent = bar.current * 100 / bar.total
	}
}

func (bar *Bar) currentBarGraphNumber() int {
	if bar.currentPercent == 100 {
		return bar.graphTotal
	} else {
		return int(float64(bar.currentPercent) * (float64(bar.graphTotal) / float64(100)))
	}
}

func (bar *Bar) Finish() {
	bar.lock.Lock()
	defer bar.lock.Unlock()
	fmt.Println()
	bar.finish = true
	bar.quit <- true
}

func (bar *Bar) SetCurrValue(curr int64) {
	bar.current = curr
	bar.countCurrentPercent()
}

func (bar *Bar) SetGraph(graph string) {
	bar.graph = graph
	bar.backGraph = ""
	for _, g := range graph {
		_, size := utf8.DecodeRuneInString(string(g))
		if size > 3 {
			bar.backGraph += "  "
		} else {
			bar.backGraph += " "
		}
	}
}

func (bar *Bar) SetGraphTotal(graphTotal int) {
	bar.graphTotal = graphTotal
}

func (bar *Bar) SetGraphColor(c color) {
	bar.color.GraphColor = c
}

func (bar *Bar) SetGraphBackColor(b backcolor) {
	bar.color.GraphBackColor = b
}

func (bar *Bar) SetDescribeColor(c color) {
	bar.color.DescribeColor = c
}

func (bar *Bar) SetRatioColor(c color) {
	bar.color.RatioColor = c
}

func (bar *Bar) SetPercentColor(c color) {
	bar.color.PercentColor = c
}

func (bar *Bar) SetDescribe(describe string) {
	bar.describe = describe
}

func (bar *Bar) ShowBytes(show bool) {
	bar.showBytes = show
}

func (bar *Bar) ShowPercent(show bool) {
	bar.showPercent = show
}

func (bar *Bar) ShowDuration(show bool) {
	bar.showDuration = show
}

func (bar *Bar) setBackGraph(bg string) {
	bar.backGraph = bg
}

func (bar *Bar) genSpace(n int) string {
	res := ""
	for i := 0; i < n; i++ {
		res += " "
	}
	return res
}

func (bar *Bar) startCalcDuration() {
	ticker := time.NewTicker(time.Second)
	go func() {
	LOOP:
		for {
			select {
			case <-bar.quit:
				break LOOP
			case <-ticker.C:
				bar.lock.Lock()
				if bar.current <= 0 {
					bar.lock.Unlock()
					break
				}
				sp := bar.current - bar.lastCount
				if sp > 0 {
					speed := (bar.total - bar.current) / sp
					bar.duration = formatTime(int(speed))
				}
				bar.lastCount = bar.current
				bar.lock.Unlock()
			}
		}
	}()
}

func (bar *Bar) startCalcRate() {
	count := bar.current - bar.lastCount
	bar.currentRate = formatBytes(float64(count*10)) + "/s"
}

func initOptions(options *Options) {
	if options.Graph == "" {
		options.Graph = "█"
	}
	if options.Color == nil {
		options.Color = &BarColor{}
	}
}

func NewBar(total int64) *Bar {
	bar := &Bar{
		total:       total,
		graph:       "█",
		backGraph:   " ",
		graphTotal:  50,
		currentRate: "0 B/s",
		color:       &BarColor{},
		quit:        make(chan bool, 1),
	}
	bar.countCurrentPercent()
	bar.startCalcDuration()
	return bar
}

func NewBarWithOptions(total int64, options *Options) *Bar {
	initOptions(options)
	bar := &Bar{
		total:        total,
		graph:        options.Graph,
		backGraph:    " ",
		graphTotal:   50,
		currentRate:  "0 B/s",
		describe:     options.Describe,
		showPercent:  options.ShowPercent,
		showBytes:    options.IsBytes,
		showDuration: options.ShowDuration,
		color:        options.Color,
		quit:         make(chan bool, 1),
	}
	bar.countCurrentPercent()
	bar.startCalcDuration()
	return bar
}
