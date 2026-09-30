package pkg_test

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/google/go-cmp/cmp"

	"github.com/sftsrv/mex/pkg"
)

const md = `# mex command used to generate this output

Pipe this into ` + "`mex commit`" + ` to save

` + "`````" + `ts my/project/route.ts:1-5
/**
 * This is some stuff
 * @example
 */
console.log("hello");
` + "`````" + `

Here are some general comments that should be excluded from the changeset

` + "`````" + `html my/project/template.html:5-10
<h1>Hello world</h1>
<style>
  .some-style, h1 {
    background-color: red;
  }
</style>
` + "`````" + `

And here are some more comments
`

const grep = `
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
`

const file = `1 some file
2
3
4
5
6
7
8
9
10
11
12
15
14
15
`

func TestMdFormat(t *testing.T) {
	cupaloy.SnapshotT(t, md)
}

func TestGrepFormat(t *testing.T) {
	cupaloy.SnapshotT(t, grep)
}

func TestParseFromMd(t *testing.T) {
	parsed, err := pkg.FromMd(md)

	if err != nil {
		t.Error(err)
	}

	cupaloy.SnapshotT(t, parsed)
}

func TestParserFromGrep(t *testing.T) {
	parsed, err := pkg.FromGrep(grep)

	if err != nil {
		t.Error(err)
	}

	cupaloy.SnapshotT(t, parsed)
}

func TestParseFromMdAndGrep(t *testing.T) {
	fromMd, err := pkg.FromMd(md)

	if err != nil {
		t.Error(err)
	}

	fromGrep, err := pkg.FromGrep(grep)

	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(fromMd, fromGrep); diff != "" {
		t.Errorf("Parsing from md is not the same as from grep:\n%s", diff)
	}
}

func TestApplyRegions(t *testing.T) {
	regions := []pkg.Region{
		{
			Path:    "my/project/route.ts",
			Start:   5,
			End:     7,
			Content: "[start]This replaces the more\n lines\n than were in in\nthe original region[end]\n",
		},
		{
			Path:    "my/project/route.ts",
			Start:   10,
			End:     12,
			Content: "[start]This replaces the same number\nof lines as in\nthe original region[end]\n",
		},
	}

	result := pkg.ApplyRegions(file, regions)

	cupaloy.SnapshotT(t, result)
}
