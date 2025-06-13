package logger

import (
	"fmt"
)

const CollectorLogPrefix = "collector"
const WebLogPrefix = "web"
const MainLogPrefix = "main"

func Logf(prefix string, frmt string, args ...any) {
	fmt.Printf("[%s] %s\n", prefix, fmt.Sprintf(frmt, args...))
}
