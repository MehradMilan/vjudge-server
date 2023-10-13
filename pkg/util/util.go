package util

import (
	"log/slog"
	"math"
	"os"
	"os/exec"
)

func SlogError(err error) slog.Attr {
	return slog.String("error", err.Error())
}

// SlogFatal will log a message with an error and kill the program
func SlogFatal(msg string, err error, attrs ...slog.Attr) {
	slog.With(attrs).With(SlogError(err)).Error(msg)
	os.Exit(1)
}

// FileExists will check either if a file exists or not
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// ExecToSTD will execute a command and print its result in stdout and stderr
func ExecToSTD(command string, arguments ...string) error {
	cmd := exec.Command(command, arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
