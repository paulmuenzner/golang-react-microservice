package main

import (
	"fmt"

	"github.com/paulmuenzner/shared/date"
	"github.com/paulmuenzner/shared/logger"
)

func main() {
	currentTime := date.GetCurrentUTCTime()
	logger.Info("Start api_gateway")
	fmt.Println("Current timet:", date.FormatDate(currentTime))
}
