package z_test

import (
	"encoding/json"
	"errors"
	"testing"

	"svc/internal/z"
)

func TestNumber(t *testing.T) {
	t.Run("Sort", testNumberSort)
	t.Run("Sum", testNumberSum)
	t.Run("Mean", testNumberMean)
	t.Run("Percentile", testNumberPercentile)
	t.Run("Quartile", testNumberQuartile)
	t.Run("Variance", testNumberVariance)
	t.Run("BoxPlot", testNumberBoxPlot)
}

var (
	empty  []int
	num    = z.NewNumbers(14999, 477, 446, 3809, 1101, 2759, 5647, 8788, 2050, 7789, 7193, 892, 8322, 4102, 212, 4613, 813, 7172, 1802, 3134, 3182, 6136, 1707, 9679, 2810, 4023, 349, 9594, 6955, 1808, 5394, 5101, 6935, 4571, 5347, 9155, 8700, 7967, 1160, 2637, 1641, 1182, 8756, 2061, 9554, 717, 4408, 4664, 7681, 8014, 9567, 533, 2867, 1804, 1405, 6882, 8558, 2564, 950, 4663, 5612, 5561, 3244, 3451, 1693, 332, 3407, 7678, 7994, 6673, 4560, 4258, 3150, 2292, 556, 3837, 2306, 5723, 3787, 8786, 9446, 8497, 928, 6847, 5650, 3276, 2736, 1753, 805, 5383, 3637, 6407, 194, 5566, 4815, 7760, 8569, 4439, 7383, 2562)
	sum    = 459354
	mean   = 4593.54
	sorted = num.Clone().Sort()
)

func testNumberSort(t *testing.T) {
	if z.SliceEqual(num, sorted) {
		t.Error("num & sorted should not be equal")
	}
	if !z.SliceEqual(num.Sort(), sorted) {
		t.Error("num.Sort & sorted should be equal")
	}
}

func testNumberSum(t *testing.T) {
	v, err := num.Sum()
	if err != nil {
		t.Error(err.Error())
	}
	if v != sum {
		t.Errorf("v (%d) should be equal to sum (%d)", v, sum)
	}
	_, err = z.NewNumbers(empty...).Sum()
	if !errors.Is(err, z.ErrEmptyNumbers) {
		t.Error("err should be sdk.ErrEmptyNumbers")
	}
}

func testNumberMean(t *testing.T) {
	v, err := num.Mean()
	if err != nil {
		t.Error(err.Error())
	}
	if v != mean {
		t.Errorf("v (%f) should be equal to mean (%f)", v, mean)
	}
	_, err = z.NewNumbers(empty...).Mean()
	if !errors.Is(err, z.ErrEmptyNumbers) {
		t.Error("err should be sdk.ErrEmptyNumbers")
	}
}

func testNumberPercentile(t *testing.T) {
	msg := "v (%f) should be equal to percentile (%f)"
	num2 := z.NewNumbers(1, 2, 9)
	v, err := num2.Median()
	if err != nil {
		t.Error(err.Error())
	}
	if percentile := 2.0; v != percentile {
		t.Errorf(msg, v, percentile)
	}

	num2 = z.NewNumbers(1, 2, 3, 9)
	v, err = num2.Median()
	if err != nil {
		t.Error(err.Error())
	}
	if percentile := 2.5; v != percentile {
		t.Errorf(msg, v, percentile)
	}

	num2 = z.NewNumbers(1, 2, 3, 9)
	v, err = num2.Percentile(.2)
	if err != nil {
		t.Error(err.Error())
	}
	if percentile := 1.0; v != percentile {
		t.Errorf(msg, v, percentile)
	}

	num2 = z.NewNumbers(1, 2, 3, 9)
	v, err = num2.Percentile(.25)
	if err != nil {
		t.Error(err.Error())
	}
	if percentile := 1.5; v != percentile {
		t.Errorf(msg, v, percentile)
	}

	num2 = z.NewNumbers(1, 2, 3, 9)
	v, err = num2.Percentile(.3)
	if err != nil {
		t.Error(err.Error())
	}
	if percentile := 2.0; v != percentile {
		t.Errorf(msg, v, percentile)
	}

	num2 = z.NewNumbers(9)
	v, err = num2.Percentile(.3)
	if err != nil {
		t.Error(err.Error())
	}
	if percentile := 9.0; v != percentile {
		t.Errorf(msg, v, percentile)
	}

	num2 = z.NewNumbers(9)
	v, err = num2.Percentile(-.3)
	if !errors.Is(err, z.ErrInvalidRange) || v != 0 {
		t.Error("err should be sdk.ErrInvalidRange")
	}
}

func testNumberQuartile(t *testing.T) {
	q1, q2, q3, iqr, err := num.Quartile()
	t.Log(q1, q2, q3, iqr, err)

	num2 := z.NewNumbers(empty...)
	q1, q2, q3, iqr, err = num2.Quartile()
	t.Log(q1, q2, q3, iqr, err)
}

func testNumberVariance(t *testing.T) {
	num2 := z.NewNumbers(1245, 1415, 1312, 1427, 1510, 1590)
	num3 := z.NewNumbers(100, 123, 129, 143, 150, 197)
	v, err := num2.PopulationVariance()
	t.Log(v, err)

	v, err = num2.SampleVariance()
	t.Log(v, err)

	v, err = num2.PopulationCovariance(num3)
	t.Log(v, err)

	v, err = num2.SampleCovariance(num3)
	t.Log(v, err)

	v, err = num2.PopulationStdDev()
	t.Log(v, err)

	v, err = num2.SampleStdDev()
	t.Log(v, err)
}

func testNumberBoxPlot(t *testing.T) {
	num := z.NewNumbers(
		1_261_261.261,
		1_540_090.090,
		2_366_674.775,
		1_447_747.748, // outlier
	).Sort()

	v, err := num.BoxPlot()
	t.Log(err)
	p, _ := json.MarshalIndent(v, "", "\t")
	t.Logf("\n%s\n", p)

	va, err := num.Variance()
	t.Log(err)
	t.Log(va)

	sd, err := num.StdDev()
	t.Log(err)
	t.Log(sd)
}
