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

` + "`````" + `html my/project/template.html:45-50
<h1>Hello world</h1>
<style>
  .some-style, h1 {
    background-color: red;
  }
</style>
` + "`````"

const grep = `
my/project/route.ts:1:/**
my/project/route.ts:2: * This is some stuff
my/project/route.ts:3: * @example
my/project/route.ts:4: */
my/project/route.ts:5:console.log("hello");
my/project/template.html:45:<h1>Hello world</h1>
my/project/template.html:46:<style>
my/project/template.html:47:  .some-style, h1 {
my/project/template.html:48:    background-color: red;
my/project/template.html:49:  }
my/project/template.html:50:</style>
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
