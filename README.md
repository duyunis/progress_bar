# progress-bar

> terminal progress bar


## Usage

### install
```bash
go get -u github.com/duyunis/progress_bar
```

### example

```go
package main

import (
	"github.com/duyunis/progress_bar"
	"time"
)

func main() {
	options := &progress_bar.Options{
		Graph:        ">",
		Describe:     "ProgressBar",
		IsBytes:      true,
		ShowPercent:  true,
		ShowDuration: true,
		Color: &progress_bar.BarColor{
			DescribeColor:  progress_bar.Yellow,
			GraphColor:     progress_bar.Blue,
			GraphBackColor: progress_bar.BackGreen,
			PercentColor:   progress_bar.Purple,
			RatioColor:     progress_bar.Green,
			DurationColor:  progress_bar.Red,
		},
	}
	bar := progress_bar.NewBarWithOptions(1000000, options)
	add := int64(100)
	for i := 0; i < 10000; i++ {
		bar.Add(add)
		time.Sleep(time.Millisecond * 1)
		add += 100
	}
	bar.Finish()
}
```

### effect

![](https://cdn.duyunis.cn/image/progressBar.gif)