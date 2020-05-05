package main

import (
	"fmt"
	"./math"
)

func main() {
	// Временной ряд представим в виде ассоциативного массива
	// в качестве ключа будет приниматься номер недели, а
	// в качестве значения значение продаж в данную неделю
	time_series, steps := InitData("data.txt")
	steps.Write()
	time_series.Write()

	aver, disp := AverageAndDispersion(time_series)
	fmt.Println("\nСреднее:")
	fmt.Println(aver)
	fmt.Println("\nДисперсия:")
	fmt.Println(disp)

	outs, errs := FirstPrediction(time_series)
	fmt.Println("\nOuts:")
	for i := range outs {
		outs[i].Write()
	}
	fmt.Println(errs)

	a0, a1, line := Coeffs(time_series, steps)

	fmt.Println("\nПрямая:")
	fmt.Println("A0 =", a0)
	fmt.Println("A1 =", a1)
	line.Write()

	_, errs2 := MakePrediction(a0, a1, time_series)
	// outs2, errs2 := MakePrediction(time_series, steps)
	// for i := range outs2 {
	// 	outs[i].Write()
	// }
	fmt.Println("\nErrs2:")
	fmt.Println(errs2)
}

func InitData(filepath string) (time_series, steps math.Vector) {
	// Экспериментальные данные
	time_series.Read("data.txt")
	// Промежутки времени (недели)
	steps.New(time_series.Size)
	for i := range steps.Array {
		steps.Array[i] = float64(i + 1)
	}
	return
}

func FirstPrediction(time_series math.Vector) (outs []math.Vector, errs []float64) {
	out1, err1 := time_series.CalculateFirstPrediction(0.1)
	out2, err2 := time_series.CalculateFirstPrediction(0.3)
	outs = append(outs, out1, out2)
	errs = append(errs, err1, err2)
	return
}

func AverageAndDispersion(time_series math.Vector) (average, dispersion float64) {
	average = time_series.GetAverage()
	dispersion = time_series.GetDispersion()
	return
}

func Coeffs(time_series, steps math.Vector) (a0, a1 float64, line math.Vector) {
	a0, a1 = math.GetInitialCoeffsLinearModel(steps, time_series)
	line.New(time_series.Size)
	for i, x := range time_series.Array {
		line.Array[i] = a0 + a1 * x
	}
	return
}

func MakePrediction(a0, a1 float64, time_series math.Vector) (outs [][]math.Vector, errs []float64){
	out3, err3 := time_series.LinearModel(0.1, a0, a1, 0)
	out4, err4 := time_series.LinearModel(0.3, a0, a1, 0)

	out5, err5 := time_series.LinearModel(0.1, a0, a1, 1)
	out6, err6 := time_series.LinearModel(0.3, a0, a1, 1)

	out7, err7 := time_series.LinearModel(0.1, a0, a1, 5)
	out8, err8 := time_series.LinearModel(0.3, a0, a1, 5)
	outs = append(outs, out3, out4, out5, out6, out7, out8)
	errs = append(errs, err3, err4, err5, err6, err7, err8)
	return
}