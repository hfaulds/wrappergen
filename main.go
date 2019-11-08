package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	args := flag.Args()
	interfaces, err := parseDir(args[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(generate(interfaces))
}
