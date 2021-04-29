package downloading

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/edm20627/gopherdojo-studyroom/kadai3-2/edm20627/option"
	"golang.org/x/sync/errgroup"
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
	bc := context.Background()
	ctx, cancel := context.WithTimeout(bc, d.options.Timeout)
	defer cancel()

	// contentLengthを取得
	contentLength, err := d.getContentLength(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// ダウンロード先を作成
	dir, err := ioutil.TempDir("", "download")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer os.RemoveAll(dir)

	// ダウンロード
	err = d.download(ctx, contentLength, dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	// マージ
	err = d.merge(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func (d *Download) getContentLength(ctx context.Context) (int, error) {
	req, err := http.NewRequest("HEAD", d.options.URL, nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)

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

func (d *Download) download(ctx context.Context, contentLength int, dir string) error {
	var preMin, min, max int
	errCh := make(chan error)

	g, ctx := errgroup.WithContext(ctx)

	for n := d.options.ParallelNum; 0 < n; n-- {
		n := n
		min = preMin
		max = contentLength/n - 1
		preMin = contentLength / n
		go d.parallelDownload(ctx, n, min, max, dir, errCh)
	}

	for n := d.options.ParallelNum; 0 < n; n-- {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-errCh:
				return err
			}
		})
	}

	// for n := d.options.ParallelNum; 0 < n; n-- {
	// 	n := n
	// 	g.Go(func() error {
	// 		min = preMin
	// 		max = contentLength/n - 1
	// 		preMin = contentLength / n
	// 		return d.parallelDownload(ctx, n, min, max, dir)

	// 	})
	// }

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (d *Download) parallelDownload(ctx context.Context, n, min, max int, dir string, errCh chan error) {
	fmt.Println(min, max)

	req, err := http.NewRequest("GET", d.options.URL, nil)
	if err != nil {
		errCh <- err
		return
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", min, max))

	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		errCh <- err
		return
	}

	defer res.Body.Close()

	file, err := os.Create(fmt.Sprintf("%v/%v-%v", dir, n, d.options.Output))
	fmt.Println(file.Name())

	if err != nil {
		errCh <- err
		return
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		errCh <- err
		return
	}

	errCh <- nil
}

func (d *Download) merge(dir string) error {
	file, err := os.Create(d.options.Output)
	if err != nil {
		return err
	}
	defer file.Close()

	for n := d.options.ParallelNum; 0 < n; n-- {
		src, err := os.Open(fmt.Sprintf("%v/%v-%v", dir, n, d.options.Output))
		if err != nil {
			return err
		}
		_, err = io.Copy(file, src)
		if err != nil {
			return err
		}
	}
	return nil
}
