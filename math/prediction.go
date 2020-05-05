package math

// Вычисление величины дисперсии ошибки предсказания
// predicted - вектор предсказанных величин
// n - порядок полинома
func (vec Vector) getErrorPredictionDispersion(predicted Vector, n float64) (dispersion float64) {
	N := float64(len(vec.Array))
	for i, value := range vec.Array {
		dispersion += (value - predicted.Array[i]) * (value - predicted.Array[i])
	}
	dispersion /= (N - n - 1.0)
	return
}

// Сделать прогнозирование значений временного ряда
// и вычислить дисперсию ошибки предсказания
func (vec Vector) CalculateFirstPrediction(alpha float64) (S Vector, error_dispersion float64) {
	N := vec.Size
	S.New(N)
	beta := 1.0 - alpha
	// Начальное значение S[0] - среднее арифметическое
	// первых 5-ти значений временного ряда
	// Ограничеваем длину среза в 5 значений
	vec.Array = vec.Array[:5]
	// Вычисляем среднее
	S.Array[0] = vec.GetAverage()
	// Возвращаем длину среза
	vec.Array = vec.Array[:cap(vec.Array)]

	// Вычислим остальные значения
	for i, value := range vec.Array[5:] {
		S.Array[i + 5] = alpha * value + beta * S.Array[i + 4]
	}
	error_dispersion = vec.getErrorPredictionDispersion(S, 0)
	return
}

// Оценка начальных коэффициентов линейной модели
// x, y - экспериментальные значения случайных величин
func GetInitialCoeffsLinearModel(x, y Vector) (a0, a1 float64) {
	N := len(x.Array)
	sum_y, sum_x := y.Sum(), x.Sum()
	var sum_xy, sum_xx float64
	for i := range x.Array {
		sum_xy += x.Array[i] * y.Array[i]
		sum_xx += x.Array[i] * x.Array[i]
	}
	a1 = (float64(N) * sum_xy - sum_x * sum_y) / (float64(N) * sum_xx - sum_x * sum_x)
	a0 = (sum_y - a1 * sum_x) / float64(N)
	return
}

// Линейная модель для прогнозирования значений временного ряда
func (vec Vector) LinearModel(alpha, a0, a1 float64, m int) (S []Vector, error_dispersion float64) {
	N := len(vec.Array)
	S = make([]Vector, 2)
	for i := range S {
		S[i].New(N)
	}
	beta := 1 - alpha
	// Определение начальных оценок
	S[0].Array[0] = a0 - (beta / alpha) * a1
	S[1].Array[0] = a0 - 2.0 * (beta / alpha) * a1

	next_a0 := make([]float64, N - m)
	next_a1 := make([]float64, N - m)
	xt := InitVector()
	xt.New(N)
	vxt := a0 - a1

	next_a0[0], next_a1[0], xt.Array[0] = a0, a1, vxt

	for i := 1; i < N - m; i++ {
		S[0].Array[i] = alpha * vec.Array[i] + beta * S[0].Array[i - 1]
		S[1].Array[i] = alpha * S[0].Array[i] + beta * S[1].Array[i - 1]
		xt.Array[i], vxt = next_a0[i - 1] + next_a1[i - 1], next_a0[i - 1]
		next_a0[i] = vec.Array[i] + beta * beta * (vxt - vec.Array[i])
		next_a1[i] = next_a1[i - 1] + alpha * alpha * (vxt - vec.Array[i])
	}

	// Оценка будущего наблюдения
	for i := N - m; i < N; i++ {
		xt.Array[i] = next_a0[i - m] + float64(m) * next_a1[i - m]
	}
	error_dispersion = vec.getErrorPredictionDispersion(xt, 1.0)
	return
}