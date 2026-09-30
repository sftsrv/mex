package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/sftsrv/mex/pkg"
	touchup "github.com/sftsrv/touchup/pkg"
)

const usage = `mex

A tool for applying multi-buffer style patches using Markdown for readability

## Usage

### Interactively

Interactive works with 'mex edit' as follows:

'''sh
grep -r -n -H my-search | mex edit

# another search command that outputs the same structure as grep may also be used, e.g. ripgrep
rg -n -H my-search | mex edit
'''

Saving the file and closing your editor will apply the changes

### In Stages

1. Generate a changeset using 'mex gen'

'''sh
grep -r -n -H my-search >> out.txt
'''

2. Make your changes as needed
3. Apply the changeset using 'mex apply':

'''sh
cat out.txt | mex commit
'''

Anything outside of the Markdown code blocks will be ignored

## Commands

- 'mex edit' - Open the changeset in your $EDITOR to be immediately edited and applied
- 'mex gen' - Generate a Markdown changeset that can be edited and applied using 'mex commit'
- 'mex commit' - Apply changes in the set to disk

## Flags

- '--help' shows this help menu

## Formats

Pipe in some content that 'mex' is able to extract from. 'mex' expects the input pipe to be the 'grep -n -H' output format (or 'rg -n -H') and joins the relevant ranges
`

func printHelp() {
	fmt.Print(usage)

}

func main() {
	helpFlag := flag.Bool("help", false, "show usage info")

	flag.Parse()

	if *helpFlag || len(os.Args) != 2 {
		printHelp()
		return
	}

	arg := os.Args[1]

	switch arg {
	case "edit":
		edit()
		return

	case "gen":
		generate()
		return

	case "commit":
		commit()
		return
	}
}

func generate() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	content := string(input)
	regions, err := pkg.FromGrep(content)
	if err != nil {
		panic(err)
	}

	output := pkg.ToMdFile("mex generate", "Pipe int `mex commit` to apply the file changes", regions)

	os.Stdout.WriteString(output)
}

func edit() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	content := string(input)
	regions, err := pkg.FromGrep(content)
	if err != nil {
		panic(err)
	}

	md := pkg.ToMdFile("mex edit", "Edit and save this file to apply the changes", regions)

	editor, err := touchup.GetDefaultEditor()
	if err != nil {
		panic(err)
	}

	result, err := touchup.EditFile(editor, "mex", "md", md)
	if err != nil {
		panic(err)
	}

	changes, err := pkg.FromMd(result)

	if err != nil {
		panic(err)
	}

	apply(changes)
}

func commit() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	content := string(input)
	regions, err := pkg.FromGrep(content)
	if err != nil {
		panic(err)
	}

	apply(regions)
}

func apply(regions []pkg.Region) {
	byFile := map[string][]pkg.Region{}

	for r := range regions {
		region := regions[r]
		existing, exists := byFile[region.Path]

		if !exists {
			byFile[region.Path] = []pkg.Region{region}
		} else {
			byFile[region.Path] = append(existing, region)
		}
	}

	for path, regions := range byFile {
		file, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}

		result := pkg.ApplyRegions(string(file), regions)

		err = os.WriteFile(path, []byte(result), 0644)
		if err != nil {
			panic(err)
		}

	}
}
