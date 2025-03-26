package main

import (
	"fmt"

	"github.com/paulmuenzner/shared/date"
	"github.com/paulmuenzner/shared/logger"
)

func main() {
	currentTime := date.GetCurrentUTCTime()
	logger.Info("Starte User-Service")
	fmt.Println("Aktuelle Zeit:", date.FormatDate(currentTime))
}
