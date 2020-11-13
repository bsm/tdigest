package tdigest_test

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/bsm/tdigest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("TDigest", func() {
	var blank, subject *tdigest.TDigest

	BeforeEach(func() {
		blank = seedTD()
		subject = seedTD(39, 15, 43, 7, 43, 36, 47, 6, 40, 49, 41)
	})

	DescribeTable("Quantile",
		func(q float64, x float64) {
			Expect(subject.Quantile(q)).To(BeNumerically("~", x, 0.01))
		},

		Entry("0%", 0.0, 6.0),
		Entry("25%", 0.25, 20.25),
		Entry("50%", 0.5, 40.0),
		Entry("75%", 0.75, 43.0),
		Entry("95%", 0.95, 48.9),
		Entry("99%", 0.99, 49.0),
		Entry("100%", 1.0, 49.0),
	)

	DescribeTable("CDF",
		func(v float64, x float64) {
			Expect(subject.CDF(v)).To(BeNumerically("~", x, 0.01))
		},

		Entry("0", 0.0, 0.0),
		Entry("10", 10.0, 0.17),
		Entry("20", 20.0, 0.25),
		Entry("30", 30.0, 0.29),
		Entry("40", 40.0, 0.5),
		Entry("50", 50.0, 1.0),
	)

	// inspired by https://github.com/aaw/histosketch/commit/d8284aa#diff-11101c92fbb1d58ccf30ca49764bf202R180
	// released into the public domain
	It("should accurately predict quantile", func() {
		N := 20_000
		Q1 := []float64{0.001, 0.01, 0.1, 0.25, 0.35, 0.65, 0.75, 0.9, 0.99, 0.999}
		Q2 := []float64{0.0001, 0.9999}

		for seed := int64(0); seed < 16; seed++ {
			r := rand.New(rand.NewSource(seed))
			t := tdigest.New()      // tdigest
			x := make([]float64, N) // exact

			for i := 0; i < N; i++ {
				num := r.NormFloat64() * 1
				t.Add(num, 1.0)
				x[i] = num
			}
			sort.Float64s(x)

			for _, q := range Q1 {
				tQ := t.Quantile(q)
				xQ := x[int(float64(len(x))*q)]
				re := math.Abs((tQ - xQ) / xQ)

				Expect(re).To(BeNumerically("<", 0.01), // allow ±1%
					"s.Quantile(%v) (got %.3f, want %.3f with seed = %v)", q, tQ, xQ, seed,
				)
			}

			for _, q := range Q2 {
				tQ := t.Quantile(q)
				xQ := x[int(float64(len(x))*q)]
				re := math.Abs((tQ - xQ) / xQ)

				Expect(re).To(BeNumerically("<", 0.1), // allow ±10%
					"s.Quantile(%v) (got %.3f, want %.3f with seed = %v)", q, tQ, xQ, seed,
				)
			}
		}
	})

	It("should reject bad quantile inputs", func() {
		Expect(math.IsNaN(blank.Quantile(0.5))).To(BeTrue())
		Expect(math.IsNaN(subject.Quantile(-0.1))).To(BeTrue())
		Expect(math.IsNaN(subject.Quantile(1.1))).To(BeTrue())
	})

	It("should calc count", func() {
		Expect(blank.Count()).To(Equal(0.0))
		Expect(subject.Count()).To(Equal(11.0))
	})

	It("should calc min", func() {
		Expect(blank.Min()).To(Equal(math.MaxFloat64))
		Expect(subject.Min()).To(Equal(float64(6)))
	})

	It("should calc max", func() {
		Expect(blank.Max()).To(Equal(-math.MaxFloat64))
		Expect(subject.Max()).To(Equal(float64(49)))
	})

	It("should merge", func() {
		t2 := seedTD(11, 2, 3, 14, 7, 4)
		Expect(t2.Count()).To(Equal(6.0))
		Expect(t2.Quantile(0.5)).To(BeNumerically("~", 5.5, 0.1))

		t2.Merge(subject)
		Expect(t2.Count()).To(Equal(17.0))
		Expect(t2.Quantile(0.5)).To(BeNumerically("~", 15.0, 0.1))
	})

	It("should add with weight", func() {
		Expect(subject.Quantile(0.5)).To(BeNumerically("~", 40.0, 0.1))
		subject.Add(6.5, 2.4)
		subject.Add(15, 3.2)
		Expect(subject.Quantile(0.5)).To(BeNumerically("~", 19.2, 0.1))
	})
})

// --------------------------------------------------------------------

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "tdigest")
}

func seedTD(vv ...float64) *tdigest.TDigest {
	h := tdigest.New()
	for _, v := range vv {
		h.Add(v, 1)
	}
	return h
}
