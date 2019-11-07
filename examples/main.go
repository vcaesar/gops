package main

import (
	"fmt"

	ps "github.com/vcaesar/gops"
)

func main() {
	fpid, err := ps.FindIds("chrome")
	if err == nil {
		fmt.Println("pids... ", fpid)
	} else {
		fmt.Println(err)
	}
}
