package downloading

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/edm20627/gopherdojo-studyroom/kadai3-2/edm20627/option"
)

type Download struct {
	options option.Options
}

func New() *Download {
	var options option.Options
	options.Parse()
	return &Download{
		options: options,
	}
}

func (d *Download) Run() int {

	// contentLengthを取得
	contentLength, err := d.getContentLength()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// 対象をGet
	res, err := d.download(contentLength)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer res.Body.Close()

	file, err := os.Create(d.options.Output)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return 0
}

func (d *Download) getContentLength() (int, error) {
	req, err := http.NewRequest("HEAD", d.options.URL, nil)
	if err != nil {
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	if res.Header.Get("Accept-Ranges") != "bytes" {
		return int(res.ContentLength), nil
	} else if int(res.ContentLength) == 0 {
		return 0, errors.New("ContentLength is 0")
	} else {
		return int(res.ContentLength), nil
	}
}

func (d *Download) download(contentLength int) (*http.Response, error) {
	req, err := http.NewRequest("GET", d.options.URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", 0, contentLength-1))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, err
}
