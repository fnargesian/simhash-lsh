# simhash-lsh
Cosine/Simhash locality sensitive hashing (LSH) in Golang with Euclidean distance sort

Implements Random projection LSH
https://en.wikipedia.org/wiki/Locality-sensitive_hashing#Random_projection

[![Build Status](https://travis-ci.org/Cogile/simhash-lsh.svg?branch=master)](https://travis-ci.org/Cogile/simhash-lsh)
[![Go Report Card](https://goreportcard.com/badge/github.com/Cogile/simhash-lsh)](https://goreportcard.com/report/github.com/Cogile/simhash-lsh)
[![GoDoc](https://godoc.org/github.com/Cogile/simhash-lsh?status.svg)](https://godoc.org/github.com/Cogile/simhash-lsh)

## Usage

Use go tool to install the package in your packages tree:
```bash
go get github.com/cogile/simhash-lsh
```

Import package as:
```go
import "github.com/cogile/simhash-lsh"
````

## Basic example
```go
func main() {
	lsh := simhashlsh.NewCosineLsh(8, 1, 6)
	lsh.Insert([]float64{1, 2, 3, 4, 5, 6, 7, 8}, "1")
	lsh.Insert([]float64{2, 3, 4, 5, 6, 7, 8, 9}, "2")
	lsh.Insert([]float64{10, 12, 99, 1, 5, 31, 2, 3}, "3")
	fmt.Println(lsh.Query([]float64{1, 2, 3, 4, 5, 6, 7, 7}))
	// [{[1 2 3 4 5 6 7 8] 1} {[2 3 4 5 6 7 8 9] 2}]
}
```
