# Tasks

## test

```sh
go test ./...
```

## snap

> Update snapshots for tests

```sh
UPDATE_SNAPSHOTS=true go test ./...
```

## gen

```sh
grep -r -n mex README.md | go run . gen
```

## edit

```sh
grep -r -n mex README.md | go run . edit
```

## rg

> Example run with `rg`

```sh
rg -n mex | go run . edit
```
