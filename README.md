# tdigest

[![Build Status](https://travis-ci.org/bsm/tdigest.png?branch=master)](https://travis-ci.org/bsm/tdigest)
[![GoDoc](https://godoc.org/github.com/bsm/tdigest?status.png)](http://godoc.org/github.com/bsm/tdigest)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This is an implementation of Ted Dunning's [t-digest](https://github.com/tdunning/t-digest/) in Go.

The implementation is based on [influxdata/tdigest](https://github.com/influxdata/tdigest).

Supported:

## Example:

```go
package main

import (
	"fmt"

	"github.com/bsm/tdigest"
)

func main() {
	// Create a new instance
	t := tdigest.New()

	// Add values
	for _, x := range []float64{1, 2, 3, 4, 5, 5, 4, 3, 2, 1} {
		t.Add(x, 1.0)
	}

	// Process and output
	fmt.Printf("MIN    : %.1f\n", t.Min())
	fmt.Printf("MAX    : %.1f\n", t.Max())
	fmt.Printf("Q.50   : %.1f\n", t.Quantile(0.5))
	fmt.Printf("Q.95   : %.1f\n", t.Quantile(0.95))
	fmt.Printf("CDF(1) : %.1f\n", t.CDF(1))
	fmt.Printf("CDF(2) : %.1f\n", t.CDF(2))
	fmt.Printf("CDF(3) : %.1f\n", t.CDF(3))

}
```
