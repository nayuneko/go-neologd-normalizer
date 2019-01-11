package main

import (
	"fmt"
	"os"

	normalizer "github.com/nayuneko/go-neologd-normalizer"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s \"string\"\n", os.Args[0])
		os.Exit(1)
	}
	fmt.Println(normalizer.NormalizeNeologd(os.Args[1]))
}
