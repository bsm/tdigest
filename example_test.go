package tdigest_test

import (
	"fmt"

	"github.com/bsm/tdigest"
)

func Example() {
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

	// Output:
	// MIN    : 1.0
	// MAX    : 5.0
	// Q.50   : 3.0
	// Q.95   : 5.0
	// CDF(1) : 0.0
	// CDF(2) : 0.3
	// CDF(3) : 0.6
}
