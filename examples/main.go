package main

import (
	"fmt"

	ps "github.com/vcaesar/gops"
)

func main() {
	findid()
	run()
}

func findid() {
	fpid, err := ps.FindIds("chrome")
	if err == nil {
		fmt.Println("pids... ", fpid)
	} else {
		fmt.Println(err)
	}
}

func run() {
	o, err := ps.Run("open /Applications/Xcode.app")
	fmt.Println("run err: ", o, err)
}
