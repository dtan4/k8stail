package main

import (
	"fmt"
)

var (
	Version  string
	Revision string
)

func printVersion() {
	fmt.Printf("k8stail version %s, %s\n", Version, Revision)
}
