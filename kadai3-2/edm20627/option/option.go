package option

import "flag"

type Options struct {
	URL         string
	Output      string
	ParallelNum int
}

func (options *Options) Parse() {
	flag.StringVar(&options.URL, "u", "https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js", "Specify url")
	flag.StringVar(&options.Output, "o", "save.js", "Specify url")
	flag.IntVar(&options.ParallelNum, "p", 2, "paralle number")
	flag.Parse()
}
