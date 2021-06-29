package signal_detection_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/edm20627/gopherdojo-studyroom/kadai3-2/edm20627/signal_detection"
)

func TestListen(t *testing.T) {
	cases := []struct {
		name    string
		signal  bool
		timeout time.Duration
	}{
		{
			name:    "execute Ctrl+C",
			signal:  true,
			timeout: 3 * time.Second,
		},
		{
			name:    "timeout",
			signal:  false,
			timeout: 1 * time.Second,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx, _ := signal_detection.Listen(c.timeout)

			if c.signal {
				doneCh := make(chan int)
				signal_detection.OsExit = func(code int) { doneCh <- code }

				process, err := os.FindProcess(os.Getpid())
				if err != nil {
					t.Fatal(err)
				}

				err = process.Signal(syscall.SIGTERM)
				if err != nil {
					t.Fatal(err)
				}
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
				t.Error("Cancellation failure")
			}
		})
	}
}
