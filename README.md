# `mex`

A mechanism for extracting and persisting multiple simultaneous file changes

## File Format

Pipe in some content that `mex` is able to extract from, expects input pipe to be the `grep -n` output format (or `rg -n`) and joins the relevant ranges

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
