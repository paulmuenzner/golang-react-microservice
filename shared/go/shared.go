// /project/shared/shared.go

package shared

import (
	"fmt"
	"log"
	"os"
)

func NewLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", prefix), log.LstdFlags)
}

func Greet(name string) string {
	return fmt.Sprintf("Hello from %s! 👋", name)
}
