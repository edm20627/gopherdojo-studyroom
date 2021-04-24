package downloading

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

	// ダウンロード先を作成
	dir, err := ioutil.TempDir("", "download")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer os.RemoveAll(dir)

	// ダウンロード
	err = d.download(contentLength, dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// マージ
	err = d.merge(dir)
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

func (d *Download) download(contentLength int, dir string) error {
	req, err := http.NewRequest("GET", d.options.URL, nil)
	if err != nil {
		return err
	}

	preMin := 0
	min, max := 0, 0
	for n := d.options.ParallelNum; 0 < n; n-- {
		min = preMin
		max = contentLength/n - 1
		preMin = contentLength / n
		fmt.Println(min, max)
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", min, max))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		file, err := os.Create(fmt.Sprintf("%v/%v-%v", dir, n, d.options.Output))
		fmt.Println(file.Name())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		defer file.Close()

		_, err = io.Copy(file, res.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	}
	return nil
}

func (d *Download) merge(dir string) error {
	file, err := os.Create(d.options.Output)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer file.Close()

	for n := d.options.ParallelNum; 0 < n; n-- {
		src, err := os.Open(fmt.Sprintf("%v/%v-%v", dir, n, d.options.Output))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		_, err = io.Copy(file, src)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	return nil
}
