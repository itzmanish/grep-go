package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/itzmanish/grep-go/cmd"
)

func main() {
	outFile := flag.String("o", "", "Output file name")
	flag.Parse()

	searchStr := flag.Arg(0)
	inpFile := flag.Arg(1)
	if flag.Arg(2) == "-o" && len(flag.Arg(3)) != 0 {
		*outFile = flag.Arg(3)
	}

	if len(searchStr) == 0 {
		fmt.Println("Usages: ./grep foo filename.txt")
		flag.PrintDefaults()
		return
	}

	result, err := cmd.Run(searchStr, inpFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	for out := range result.Result {
		_, err = cmd.Write(*outFile, out)
		if err != nil {
			log.Fatal(err)
		}
	}
}
