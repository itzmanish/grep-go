package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/itzmanish/grep-go/cmd"
)

func main() {
	out_file := flag.String("o", "", "Output file name")
	flag.Parse()

	search_str := flag.Arg(0)
	inp_file := flag.Arg(1)
	if flag.Arg(2) == "-o" && len(flag.Arg(3)) != 0 {
		*out_file = flag.Arg(3)
	}

	if len(search_str) == 0 {
		fmt.Println("Usages: ./grep foo filename.txt")
		flag.PrintDefaults()
		return
	}

	result, err := cmd.Run(search_str, inp_file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	_, err = cmd.Write(*out_file, result)
	if err != nil {
		log.Fatal(err)
	}
}
