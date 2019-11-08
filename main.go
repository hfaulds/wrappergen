package main

import (
	"flag"
	"fmt"
)

func main() {
	flag.Parse()
	args := flag.Args()
	pkg, err := ParseDir(args[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(Generate(pkg))
}
