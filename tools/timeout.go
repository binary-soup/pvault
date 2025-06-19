package tools

import "time"

func Timeout(seconds float32) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
