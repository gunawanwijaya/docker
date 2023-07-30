package z

import (
	"errors"
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

func NewNumbers[T Number](a ...T) Numbers[T] { return Numbers[T](a) }

type Numbers[T Number] []T // generic
type Number interface {
	constraints.Integer | constraints.Float
}

func (n Numbers[T]) Len() int           { return len(n) }
func (n Numbers[T]) Less(i, j int) bool { return n[i] < n[j] }
func (n Numbers[T]) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n Numbers[T]) Sort() Numbers[T]   { sort.Sort(n); return n }
func (n Numbers[T]) Clone() (o Numbers[T]) {
	o = make(Numbers[T], n.Len())
	for i, e := range n {
		o[i] = e
	}
	return o
}

var ErrEmptyNumbers = errors.New("sdk:Numbers: Empty numbers.")
var ErrInvalidRange = errors.New("sdk:Numbers: Invalid range.")
var ErrDifferentDimension = errors.New("sdk:Numbers: Different dimension.")

type ( // type alias
	f = float64
	e = error
)

func (n Numbers[T]) Sum() (v T, err e) {
	if l := n.Len(); l == 0 {
		return v, ErrEmptyNumbers
	}
	for _, e := range n {
		v = v + e
	}
	return v, nil
}

func (n Numbers[T]) Mean() (v f, err e) {
	if l := n.Len(); l == 1 {
		return f(n[0]), nil
	} else if s, err := n.Sum(); err != nil {
		return v, err
	} else {
		return f(s) / f(l), nil
	}
}

func (n Numbers[T]) Median() (v f, err e) {
	return n.Percentile(.50)
}

func (n Numbers[T]) Percentile(p f) (v f, err e) {
	if p <= 0 || p >= 1 {
		return v, ErrInvalidRange
	} else if l := n.Len(); l == 0 {
		return v, ErrEmptyNumbers
	} else if l == 1 {
		return n.Mean()
	} else {
		xf := p * f(l)
		xi := int(xf)
		// check whether p * len is integer
		if xf == f(xi) && xi == int(xf) {
			return n[xi-1 : xi+1].Mean()
		}
		// this slice contains exactly one item
		return n[xi : xi+1].Mean()
	}
}

func (n Numbers[T]) Quartile() (q1, q2, q3, iqr f, err e) {
	if q2, err = n.Median(); err != nil {
		return
	}
	q3, _ = n.Percentile(.75)
	q1, _ = n.Percentile(.25)
	return q1, q2, q3, q3 - q1, nil
}

func (n Numbers[T]) Variance() (v f, err e)           { return n.SampleVariance() }
func (n Numbers[T]) PopulationVariance() (v f, err e) { return n.xVar(false) }
func (n Numbers[T]) SampleVariance() (v f, err e)     { return n.xVar(true) }
func (n Numbers[T]) xVar(isSample bool) (v f, err e) {
	if l := n.Len(); l == 0 {
		return v, ErrEmptyNumbers
	} else if l == 1 {
		return 0, nil
	} else if m, err := n.Mean(); err != nil {
		return v, err
	} else {
		for _, n := range n {
			f := f(n)
			d := f - m
			v += d * d
		}
		if isSample {
			l -= 1
		}
		return v / f(l), nil

	}
}

func (n Numbers[T]) Covariance(x Numbers[T]) (v f, err e)           { return n.PopulationCovariance(x) }
func (n Numbers[T]) PopulationCovariance(x Numbers[T]) (v f, err e) { return n.xCov(false, x) }
func (n Numbers[T]) SampleCovariance(x Numbers[T]) (v f, err e)     { return n.xCov(true, x) }
func (n Numbers[T]) xCov(isSample bool, x Numbers[T]) (v f, err e) {
	if ln := n.Len(); ln == 0 {
		return v, ErrEmptyNumbers
	} else if lx := x.Len(); lx == 0 {
		return v, ErrEmptyNumbers
	} else if lx != ln {
		return v, ErrDifferentDimension
	} else if mn, err := n.Mean(); err != nil {
		return v, err
	} else if mx, err := x.Mean(); err != nil {
		return v, err
	} else {
		for i := 0; i < ln; i++ {
			fn, fx := f(n[i]), f(x[i])
			dn, dx := (fn - mn), (fx - mx)
			v += (dn * dx)
		}
		if isSample {
			ln -= 1
		}
		return v / f(ln), nil
	}
}

func (n Numbers[T]) StdDev() (v f, err e)           { return n.SampleStdDev() }
func (n Numbers[T]) PopulationStdDev() (v f, err e) { return n.xS(false) }
func (n Numbers[T]) SampleStdDev() (v f, err e)     { return n.xS(true) }
func (n Numbers[T]) xS(isSample bool) (v f, err e) {
	if v, err = n.xVar(isSample); err != nil {
		return
	}
	return math.Sqrt(v), nil
}

func (n Numbers[T]) MedianAbsDev() (v f, err e) {
	x := make([]f, n.Len())
	m, _ := n.Median()
	for i := 0; i < n.Len(); i++ {
		x[i] = math.Abs(f(n[i]) - m)
	}
	return NewNumbers(x...).Median()
}

type BoxPlotResult[T Number] struct {
	// Lower and Upper bound of the box plot
	//  1.5*IQR
	Lower, Upper float64

	// Outliers contains elements outside of the `Lower` & `Upper` bound
	Outliers Numbers[T]

	// P02, P09, P91 and P98 are the division used in the alternative box plot
	// along with P25 (First Quartile), P50 (Second Quartile) and P75 (Third Quartile)
	P02, P09, P91, P25, P50, P75, P98 float64

	// Quarter1, Quarter2, Quarter3, Quarter4 contains element divided by the quartiles,
	// while Quarter2 and Quarter3 are used by both of standard & alternative box plot,
	// Quarter1 and Quarter4 are used exclusively by the standard box plot EXCEPT those element
	// listed in the outliers
	Quarter1, Quarter2, Quarter3, Quarter4 Numbers[T]

	// Below02, Below09, Below25, Above75, Above91 and Above98 contains element
	// exlusive separated by P02, P09, P91, P98 plus P25 (First Quartile) and P75 (Third Quartile)
	// usually used in the normal distribution that dont have any outlier to better divide the
	// distribution of https://en.wikipedia.org/wiki/Seven-number_summary, this seven-number summary
	// divided the distribution into 8 groups, Below02, Below09, Below25, Above75, Above91, Above98
	// plus Quarter2 and Quarter3
	Below02, Below09, Below25,
	Above75, Above91, Above98 Numbers[T]
}

func (n Numbers[T]) BoxPlot() (v BoxPlotResult[T], err e) {
	q1, q2, q3, iqr, err := n.Quartile()
	if err != nil {
		return v, err
	}

	v.Lower = q1 - (1.5 * iqr)
	v.Upper = q3 + (1.5 * iqr)
	v.P02, _ = n.Percentile(.02) // (better:  2.15%)
	v.P09, _ = n.Percentile(.09) // (better:  8.87%)
	v.P25, v.P50, v.P75 = q1, q2, q3
	v.P91, _ = n.Percentile(.91) // (better: 91.13%)
	v.P98, _ = n.Percentile(.98) // (better: 97.85%)

	for _, e := range n {
		f := f(e)
		switch {
		case f < v.Lower || f > v.Upper:
			v.Outliers = append(v.Outliers, e)
		case f < q1:
			v.Quarter1 = append(v.Quarter1, e)
		case f >= q1 && f < q2:
			v.Quarter2 = append(v.Quarter2, e)
		case f >= q2 && f < q3:
			v.Quarter3 = append(v.Quarter3, e)
		case f >= q3:
			v.Quarter4 = append(v.Quarter4, e)
		}
		switch {
		case f < v.P02:
			v.Below02 = append(v.Below02, e)
		case f >= v.P02 && f < v.P09:
			v.Below09 = append(v.Below09, e)
		case f >= v.P09 && f < q1:
			v.Below25 = append(v.Below25, e)
		case f >= q3 && f < v.P91:
			v.Above75 = append(v.Above75, e)
		case f >= v.P91 && f < v.P98:
			v.Above91 = append(v.Above91, e)
		case f >= v.P98:
			v.Above98 = append(v.Above98, e)
		}
	}

	return
}
