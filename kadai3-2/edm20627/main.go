package main

import (
	"os"

	"github.com/edm20627/gopherdojo-studyroom/kadai3-2/edm20627/downloading"
)

func main() {
	os.Exit(downloading.New().Run())
}
