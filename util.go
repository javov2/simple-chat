package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func runCmd(name string, arg ...string) {
	fmt.Println("Testing")
	fmt.Println("Hey")
	fmt.Println("looking gor autocompletition")
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("Yuppp")
}

func ClearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}
