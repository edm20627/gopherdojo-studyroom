package main

import (
	"fmt"
	"os"

	"github.com/edm20627/gopherdojo-studyroom/kadai3-2/edm20627/downloading"
)

func main() {
	fmt.Println("Download started")
	var client = &downloading.Client{downloading.NewDownload()}

	status := client.Run()
	switch status {
	case downloading.StatusOK:
		fmt.Println("Download complete")
	case downloading.StatusErr:
		fmt.Println("Download abnormal termination")
	}
	os.Exit(status)
}
