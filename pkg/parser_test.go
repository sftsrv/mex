package pkg_test

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
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

` + "`````" + `html my/project/template.html:45-51
<h1>Hello world</h1>
<style>
  .some-style, h1 {
    background-color: red;
  }
</style>
` + "`````" + `
`

const grep = `
my/project/route.ts:1:/**
my/project/route.ts:2: * This is some stuff
my/project/route.ts:3: * @example
my/project/route.ts:4: */
my/project/route.ts:5:console.log("hello");
my/project/template.html:1:<h1>Hello world</h1>
my/project/template.html:2:<style>
my/project/template.html:3:  .some-style, h1 {
my/project/template.html:4:    background-color: red;
my/project/template.html:5:  }
my/project/template.html:6:</style>
`

func TestParseFromMd(t *testing.T) {
	parsed := pkg.FromMd(md)

	cupaloy.SnapshotT(t, parsed)
}

func TestParserFromGrep(t *testing.T) {
	parsed := pkg.FromGrep(grep)

	cupaloy.SnapshotT(t, parsed)
}
