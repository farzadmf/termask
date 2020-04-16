package main

import (
	"fmt"
	"os"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(1)
	}
}
