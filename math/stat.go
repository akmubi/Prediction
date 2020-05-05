package math

import "math"
// import "fmt"

func (vec Vector) GetAverage() (average float64) {
	for _, value := range vec.Array {
		average += value 
	}
	average /= float64(len(vec.Array))
	return
}

func (vec Vector) GetDispersion() (dispersion float64) {
	average := vec.GetAverage()
	for _, value := range vec.Array {
		dispersion += (value - average) * (value - average)
	}
	dispersion /= float64(len(vec.Array))
	// dispersion = math.Sqrt(dispersion)
	return
}

// Средние значения столбцов матрицы
func (mat Matrix) GetAverages() (averages Vector) {
	vectors := mat.ConvertToVec()
	averages.New(len(vectors))
	for i := range averages.Array {
		averages.Array[i] = vectors[i].GetAverage()
	} 
	return
}

// Дисперсии столбцов матрицы
func (mat Matrix) GetDispersions() (dispersions Vector) {
	// Создаём вектор дисперсий
	dispersions.New(mat.Column_count)

	// Вычисляем средние значения столбцов
	averages := mat.GetAverages()

	// Вычисляем дисперсии
	for j := 0; j < mat.Column_count; j++ {
		for i := 0; i < mat.Row_count; i++ {
			dispersions.Array[j] += (mat.Array[i][j] - averages.Array[j]) * (mat.Array[i][j] - averages.Array[j]) 
		}
		dispersions.Array[j] /= float64(mat.Row_count)
		// dispersions.Array[j] = math.Sqrt(dispersions.Array[j])
	}
	return
} 

// Выполнить стандартизацию матрицы
func (mat *Matrix) Standartize() {
	// Вычисляем средние и дисперсии каждого столбца матрицы 
	averages := mat.GetAverages()
	dispersions := mat.GetDispersions()

	for i := range mat.Array {
		for j := range mat.Array[i] {
			mat.Array[i][j] = (mat.Array[i][j] - averages.Array[j]) /  math.Sqrt(dispersions.Array[j])
		}
	}
}

// Вычислить ковариационную матрицу
func (mat Matrix) GetCovariation() (result Matrix) {
	result.New(mat.Column_count, mat.Column_count)
	// Вычисляем средние по столбцам
	averages := mat.GetAverages() // Vector
	for i := 0; i < mat.Column_count; i++ {
		for j := 0; j < mat.Column_count; j++ {
			for k := 0; k < mat.Row_count; k++ {
				result.Array[i][j] += (mat.Array[k][i] - averages.Array[i]) * (mat.Array[k][j] - averages.Array[j])
			}
			result.Array[i][j] /= float64(mat.Row_count)
		}
	}
	return
}

// Вычислить корреляционную матрицу
func (mat Matrix) GetCorrelation() (result Matrix) {
	result.New(mat.Column_count, mat.Column_count)

	for i := 0; i < mat.Column_count; i++ {
		for j := 0; j < mat.Column_count; j++ {
			for k := 0; k < mat.Row_count; k++ {
				result.Array[i][j] += mat.Array[k][i] * mat.Array[k][j]
			}
			result.Array[i][j] /= float64(mat.Row_count)
		}
	}
	return
}

// Расчёт проекций объектов на главные компоненты
func CalculateMainComponents(standartized Matrix, eigenvectors []Vector) (main_components []Vector) {
	// Преобразуем стандартизованную матрицу в векторы-признаки
	samples := standartized.ConvertToVec() // []Vector
	// Длина собственных векторов
	P := eigenvectors[0].Size
	// Количество значений признаков
	N := samples[0].Size
	main_components = make([]Vector, P)
	for i := range main_components {
		main_components[i].New(N)
		for j, v := range samples {
			mul := v.MulScalar(eigenvectors[j].Array[i])
			main_components[i].Add(mul)
		}
	}
	return
}

// Проверка равенства сумм выборочных дисперсий исходных признаков и
// выборочных дисперсий проекций объектов на главные компоненты
func CheckDispersionEquality(standartized []Vector, main_components []Vector) (sum_dipers1, sum_dipers2 float64) {
	for i := range standartized {
		sum_dipers1 += standartized[i].GetDispersion()
	}
	for i := range main_components {
		sum_dipers2 += main_components[i].GetDispersion()
	}
	return
}

// Вычисление относительной доли разброса I(p')
func CalculateIValue(eigenvalues Vector) (int, float64) {
	var p, next_p int
	var full_sum float64
	// Количество собственных значений
	p = eigenvalues.Size
	// Сумма всех собственных значений
	full_sum = eigenvalues.getSum(p)

	// Находим минимальный next_p, при котором I(next_p) > 0.95
	for next_p = p - 1; eigenvalues.getSum(next_p) / full_sum > 0.95; next_p-- {}
	// Возвращаем next_p и I(next_p)
	return next_p + 1, eigenvalues.getSum(next_p + 1) / full_sum
}

// Проверка корреляционной матрицы на значимое отличие от единичной матрицы
func (mat Matrix) ExistDifference(N int) (d float64) {
	mat.checkSquareness()
	for i := range mat.Array {
		for _, v := range mat.Array[i] {
			d += v * v
		}
	}
	d /= float64(N)
	return
}