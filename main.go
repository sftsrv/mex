package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	touchup "github.com/sftsrv/touchup/pkg"
)

const usage = ``

func printHelp() {
	fmt.Print(usage)
	flag.Usage()

}

func main() {
	helpFlag := flag.Bool("help", false, "show usage info")

	flag.Parse()

	if *helpFlag || len(os.Args) > 2 {
		printHelp()
		return
	}

	if len(os.Args) == 1 {
		read()
	}

	if len(os.Args) == 2 {
		arg := os.Args[1]

		switch arg {
		case "edit":
			edit()
			return

		case "commit":
			commit()
			return
		}
	}

	printHelp()
}

func read() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	fmt.Println(input)
}

func edit() {
	defaultEditor, _ := touchup.GetDefaultEditor()
	fmt.Println(defaultEditor)
}

func commit() {}
