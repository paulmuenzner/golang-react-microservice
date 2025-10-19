package logger_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/app/shared/go/utils/logger"
)

func TestLoggingLevels(t *testing.T) {
	// Setup: Capture output
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	logger.SetLoggerWriters(wOut, wErr)

	// Log some messages
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Flush and close writers
	wOut.Close()
	wErr.Close()

	// Read captured output
	var outBuf, errBuf bytes.Buffer
	io.Copy(&outBuf, rOut)
	io.Copy(&errBuf, rErr)

	output := outBuf.String() + errBuf.String()

	// Check expected output
	expected := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	for _, level := range expected {
		if !strings.Contains(output, level) {
			t.Errorf("log output missing level %s", level)
		}
	}
}
