package main

import (
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Database file name must provide")
	}
}
