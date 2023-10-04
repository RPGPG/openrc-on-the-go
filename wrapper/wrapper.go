package wrapper

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CheckServiceStatus(service string) []byte {
	output, _ := exec.Command("rc-service", service, "status").CombinedOutput()
	return output
}

func IsServiceStarted(service string) bool {
	output := CheckServiceStatus(service)
	return strings.Contains(OutTrimmed(output), "started")
}

func IsOpenrcExecutable() bool {
	_, err := exec.Command("which", "openrc").CombinedOutput()
	_, err2 := exec.Command("which", "rc-service").CombinedOutput()
	return (err == err2) && err == nil
}

func ClearTerminal() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func OutTrimmed(inp []byte) string {
	return strings.TrimSpace(string(inp[:]))
}

func CheckOS() {
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		log.Fatal("This program is designed to run on Alpine Linux 🐧")
	}
}
