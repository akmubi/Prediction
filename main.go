package main

import (
	"fmt"
	"math"
	stat "./math"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"strconv"
)

const (
	plot_width = 20 * vg.Centimeter
	plot_height = 10 * vg.Centimeter
)

func main() {
	time_series, steps := InitData("data.txt")
	steps.Write()
	time_series.Write()

	aver, disp := AverageAndDispersion(time_series)
	fmt.Println("\nМат. ожидание:")
	fmt.Println(aver)

	fmt.Println("\nДисперсия:")
	fmt.Println(disp)

	outs, errs := FirstPrediction(time_series)
	fmt.Println("\nS[0] (alpha = 0.1):")
	outs[0].Write()
	fmt.Println("\nS[0] (alpha = 0.3):")
	outs[1].Write()
	fmt.Println("\nОшибки предсказания для S[0] (alpha = 0.1) и S[0] (alpha = 0.3) соответственно:")
	fmt.Println(errs)

	a0, a1, line := Coeffs(time_series, steps)
	fmt.Println("\nПрямая:")
	line.Write()

	fmt.Println("A0 =", a0)
	fmt.Println("A1 =", a1)

	outs2, errs2 := MakePrediction(a0, a1, time_series)
	// outs2, errs2 := MakePrediction(time_series, steps)
	// for i := range outs2 {
	// 	outs[i].Write()
	// }
	fmt.Println("\nErrs2:")
	fmt.Println(errs2)

	//
	// ГРАФИКИ
	//

	chart, _ := plot.New()
	// S[1] alpha = 0.1
	s1a1, _ := plot.New()
	// S[1] alpha = 0.3
	s1a3, _ := plot.New()
	// S[2] alpha = 0.1
	s2a1, _ := plot.New()
	// S[2] alpha = 0.3
	s2a3, _ := plot.New()

	// Названия
	chart.Title.Text = "График временного ряда"
	chart.X.Label.Text = "Недели"
	chart.Y.Label.Text = "Экспериментальные значения"

	// Сетка
	chart.Add(plotter.NewGrid())
	
	// plotutil.AddLinePoints(chart, "S0 (alpha = 0.1)", convertOutToXYs(steps.Array, append(time_series.Array, outs[0].Array...)), "S0 (alpha = 0.3)", convertOutToXYs(steps.Array, append(time_series.Array, outs[1].Array...)))

	// plotutil.AddLinePoints(chart, "Аппроксимирующая линия", convertOutToXYs(steps, line), convertOutToXYs(steps, time_series))
	plotutil.AddLinePoints(chart,	"Значения временного ряда", convertOutToXYs(steps, time_series))//,
									// "S0 (alpha = 0.1)", convertOutToXYs(steps, outs2[0][0]),
									// "S1 (alpha = 0.1)", convertOutToXYs(steps, outs2[0][1]),
									// "S0 (alpha = 0.3)", convertOutToXYs(steps, outs2[1][0]),
									// "S1 (alpha = 0.3)", convertOutToXYs(steps, outs2[1][1]))
	// S[1] alpha = 0.1
	s1a1.Title.Text = "График временного ряда и S[1] (alpha = 0.1)"
	s1a1.X.Label.Text = "Недели"
	s1a1.Y.Label.Text = "Экспериментальные значения"	
	plotutil.AddLinePoints(s1a1,	"Значения временного ряда", convertOutToXYs(steps, time_series),
									"S1 (alpha = 0.1)", convertOutToXYs(steps, outs2[0][0]))
	s1a1.Save(plot_width, plot_height, "S11.jpeg")

	// S[1] alpha = 0.3
	s1a3.Title.Text = "График временного ряда и S[1] (alpha = 0.3)"
	s1a3.X.Label.Text = "Недели"
	s1a3.Y.Label.Text = "Экспериментальные значения"
	plotutil.AddLinePoints(s1a3,	"Значения временного ряда", convertOutToXYs(steps, time_series),
									"S1 (alpha = 0.3)", convertOutToXYs(steps, outs2[0][1]))
	s1a3.Save(plot_width, plot_height, "S13.jpeg")

	// S[2] alpha = 0.1
	s2a1.Title.Text = "График временного ряда и S[2] (alpha = 0.1)"
	s2a1.X.Label.Text = "Недели"
	s2a1.Y.Label.Text = "Экспериментальные значения"
	plotutil.AddLinePoints(s2a1,	"Значения временного ряда", convertOutToXYs(steps, time_series),
									"S2 (alpha = 0.1)", convertOutToXYs(steps, outs2[1][0]))
	s2a1.Save(plot_width, plot_height, "S21.jpeg")

	// S[2] alpha = 0.3
	s2a3.Title.Text = "График временного ряда и S[2] (alpha = 0.3)"
	s2a3.X.Label.Text = "Недели"
	s2a3.Y.Label.Text = "Экспериментальные значения"
	plotutil.AddLinePoints(s2a3,	"Значения временного ряда", convertOutToXYs(steps, time_series),
									"S2 (alpha = 0.3)", convertOutToXYs(steps, outs2[1][1]))
	s2a3.Save(plot_width, plot_height, "S23.jpeg")

	// Сохранение графика в файл
	chart.Save(plot_width, plot_height, "example.jpeg")

	new_chart, _ := plot.New()

	// Названия
	new_chart.Title.Text = "График временного ряда"
	new_chart.X.Label.Text = "Недели"
	new_chart.Y.Label.Text = "Экспериментальные значения"

	// Сетка
	new_chart.Add(plotter.NewGrid())
	new_line := make(plotter.XYs, 2)
	new_line[0].X = 0.0
	new_line[0].Y = a0
	new_line[1].X = 30.0
	new_line[1].Y = a0 + 20.0 * a1
	// plotutil.AddLinePoints(new_chart,	"Значения временного ряда", convertOutToXYs(steps, time_series),
	// 								"Аппроксимирующая линия", new_line)

	fmt.Println("Разности:")
	results := getDiffs(time_series)
	div_results := getDivs(time_series)
	plotutil.AddLinePoints(new_chart, "0", convertOutToXYs(steps, results))
	results.Write()
	div_results.Write()
	fmt.Println("Среднее:", results.GetAverage())
	for i := 1; results.Size > 1; i++ {
		fmt.Println("Разности №", i + 1, ":")
		results = getDiffs(results)
		plotutil.AddLinePoints(new_chart, strconv.Itoa(i), convertOutToXYs(steps, results))
		results.Write()
		fmt.Println("Среднее:", results.GetAverage())
	}	
	new_chart.Save(plot_width * 2, plot_height * 2, "example-line.jpeg")
}

func convertOutToXYs(steps, time_series stat.Vector) plotter.XYs {
	points := make(plotter.XYs, time_series.Size)
	for i := range points {
		points[i].X = steps.Array[i]
		points[i].Y = time_series.Array[i]
	}
	return points
}

func getDiffs(time_series stat.Vector) (result stat.Vector) {
	result.New(time_series.Size - 1)
	for i := 0; i < time_series.Size - 1; i++ {
		result.Array[i] = time_series.Array[i + 1] - time_series.Array[i]
	}
	return
}

func getDivs(time_series stat.Vector) (result stat.Vector) {
	result.New(time_series.Size - 1)
	for i := 0; i < time_series.Size - 1; i++ {
		result.Array[i] = time_series.Array[i + 1] / time_series.Array[i]
	}
	return
}

func InitData(filepath string) (time_series, steps stat.Vector) {
	// Экспериментальные данные
	time_series.Read(filepath)
	// Промежутки времени (недели)
	steps.New(time_series.Size)
	for i := range steps.Array {
		steps.Array[i] = float64(i + 1)
	}
	return
}

func FirstPrediction(time_series stat.Vector) (outs []stat.Vector, errs []float64) {
	out1, err1 := time_series.CalculateFirstPrediction(0.1)
	out2, err2 := time_series.CalculateFirstPrediction(0.3)
	outs = append(outs, out1, out2)
	errs = append(errs, err1, err2)
	return
}

func AverageAndDispersion(time_series stat.Vector) (average, dispersion float64) {
	average = time_series.GetAverage()
	dispersion = math.Sqrt(time_series.GetDispersion())
	return
}

func Coeffs(time_series, steps stat.Vector) (a0, a1 float64, line stat.Vector) {
	a0, a1 = stat.GetInitialCoeffsLinearModel(steps, time_series)
	line.New(time_series.Size)
	for i, x := range time_series.Array {
		line.Array[i] = a0 + a1 * x
	}
	return
}

func MakePrediction(a0, a1 float64, time_series stat.Vector) (outs [][]stat.Vector, errs []float64){
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