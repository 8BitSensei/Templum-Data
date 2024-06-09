package main

import (
	"fmt"
	"os"
)

func main() {
	switch os.Args[1] {
	case "meta":
		GenerateMetadata()
	case "bib":
		GenerateBib()
	default:
		fmt.Println("Unrecognised argument, please use bib or meta")
	}
}
