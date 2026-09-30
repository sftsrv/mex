# `mex`

A tool for applying patches using Markdown for readability

## Installation

```sh
go install github.com/sftsrv/mex
```

## Usage

### Interactively

Interactive works with `mex edit` as follows:

```sh
grep -r -n -H my-search | mex edit

# another search command that outputs the same structure as grep may also be used, e.g. ripgrep
rg -n -H my-search | mex edit
```

Saving the file and closing your editor will apply the changes

### In Stages

1. Generate a changeset using `mex gen`

```sh
grep -r -n -H my-search >> out.txt
```

2. Make your changes as needed
3. Apply the changeset using `mex apply`:

```sh
cat out.txt | mex commit
```

Anything outside of the Markdown code blocks will be ignored

## Commands

- `mex edit` - Open the changeset in your $EDITOR to be immediately edited and applied
- `mex gen` - Generate a Markdown changeset that can be edited and applied using `mex commit`
- `mex commit` - Apply changes in the set to disk

## Flags

- `--help` shows this help menu

## Formats

Pipe in some content that `mex` is able to extract from. `mex` expects the input pipe to be the `grep -n -H` output format (or `rg -n -H`) and joins the relevant ranges

The input file should have the following structure:

```txt
my/project/route.ts:1:/**
my/project/route.ts:2: * This is some stuff
my/project/route.ts:3: * @example
my/project/route.ts:4: */
my/project/route.ts:5:console.log("hello");
my/project/template.html:5:<h1>Hello world</h1>
my/project/template.html:6:<style>
my/project/template.html:7:  .some-style, h1 {
my/project/template.html:8:    background-color: red;
my/project/template.html:9:  }
my/project/template.html:10:</style>
```

The resulting changeset file that can be edited looks like so:

```````md
# mex command used to generate this output

Pipe this into `mex commit` to save

`````ts my/project/route.ts:45-65
/**
 * This is some stuff
 * @example
 */
console.log("hello");
`````

`````html my/project/template.html:45-65
<h1>Hello world</h1>
<style>
  .some-style, h1 {
    background-color: red;
  }
</style>
`````
```````
