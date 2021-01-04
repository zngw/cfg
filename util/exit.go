package util

import (
	"fmt"
	"os"
)

func WaitExit(code int) {
	fmt.Printf("请按任意键继续...")
	b := make([]byte, 1)
	_, _ = os.Stdin.Read(b)
	os.Exit(code)
}
