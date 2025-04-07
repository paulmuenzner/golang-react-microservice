package main

import (
	"fmt"

	"github.com/paulmuenzner/shared/date"
	logger "github.com/paulmuenzner/shared/logging"
)

func main() {
	currentTime := date.GetCurrentUTCTime()
	logger.Info("Starte user_service")
	fmt.Println("Current timet:", date.FormatDate(currentTime))
}
