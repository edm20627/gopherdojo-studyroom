package option

import (
	"flag"
	"time"
)

type Options struct {
	URL         string
	Output      string
	ParallelNum int
	Timeout     time.Duration
}

func (options *Options) Parse() {
	flag.StringVar(&options.URL, "u", "https://golang.org/doc/gopher/appenginegophercolor.jpg", "Download destination url")
	flag.StringVar(&options.Output, "o", "download.jpg", "Download file name")
	flag.IntVar(&options.ParallelNum, "p", 4, "Paralle number")
	flag.DurationVar(&options.Timeout, "t", 10*time.Second, "Timeout")
	flag.Parse()
}
