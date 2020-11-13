package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/bsm/tdigest"
)

var flags struct {
	iterations int
	maxEntries int
	randSeed   int64
}

func init() {
	flag.IntVar(&flags.iterations, "n", 100_000, "Number of iterations")
	flag.IntVar(&flags.maxEntries, "max-entries", 100_000, "Maximum number of entries per digest")
	flag.Int64Var(&flags.randSeed, "seed", 0, "Random seed")
}

func main() {
	flag.Parse()

	if flags.randSeed == 0 {
		flags.randSeed = time.Now().Unix()
	}
	rnd := rand.New(rand.NewSource(flags.randSeed))

	log.Printf("fuzzing in progress [iterations:%d seed:%d]", flags.iterations, flags.randSeed)
	for i := 0; i < flags.iterations; i++ {
		compression := rnd.Float64() * 2_000
		numEntries := rnd.Intn(flags.maxEntries)

		td := tdigest.NewWithCompression(compression)
		for j := 0; j < numEntries; j++ {
			v := rnd.NormFloat64() * 1e4
			w := rnd.NormFloat64()*1 + 1
			td.Add(v, w)
		}

		td.Quantile(0.0001)
		td.Quantile(0.001)
		td.Quantile(0.01)
		td.Quantile(0.1)
		td.Quantile(0.5)
		td.Quantile(0.9)
		td.Quantile(0.99)
		td.Quantile(0.999)
		td.Quantile(0.9999)

		td.CDF(td.Min() - 1)
		td.CDF(td.Quantile(0.5))
		td.CDF(td.Max() + 1)

		if n := (i + 1); n < flags.iterations && n%(flags.iterations/20) == 0 {
			log.Printf("completed %d iterations", n)
		}
	}
	log.Printf("completed %d iterations", flags.iterations)
}
