package math

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"errors"
	// "math"
)

// Структура матрицы
type Matrix struct {
	Array [][]float64
	Row_count, Column_count int
}

// Структура вектора
type Vector struct {
	Array []float64
	Size int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Инициализация матрицы
func InitMatrix() Matrix {
	return Matrix{ Array: nil, Row_count: 0, Column_count: 0 }
}

// Инициализация вектора
func InitVector() Vector {
	return Vector{ Array: nil, Size: 0 }
}

func (vec *Vector) New(size int) {
	vec.Array = make([]float64, size)
	vec.Size = size
}

func (mat *Matrix) New(row_count, column_count int) {
	// Выделяем память под матрицу
	mat.Array = make([][]float64, row_count)
	for i := range mat.Array {
		mat.Array[i] = make([]float64, column_count)
	}
	mat.Row_count = row_count
	mat.Column_count = column_count
}

// Единичная матрица
func Identity(mat_size int) Matrix {
	mat := InitMatrix()
	mat.New(mat_size, mat_size)
	for i := range mat.Array {
		for j := range mat.Array[i] {
			if i == j {
				mat.Array[i][j] = 1.0
			}
		}
	}
	return mat
}

// Проверить матрицу на "квадратность"
// Если матрица не квадратная, то выдать ошибку
func (mat Matrix) checkSquareness() {
	if mat.Column_count != mat.Row_count {
		check(errors.New("Матрица должна быть квадратной!"))
	} 
}

// Проверка симметричности матрицы
func (mat Matrix) IsSimmetric() bool {
	mat.checkSquareness()
	for i := 0; i < mat.Row_count - 1; i++ {
		for j := i + 1; j < mat.Column_count; j++ {
			if mat.Array[i][j] != mat.Array[j][i] {
				return false
			} 
		}
	}
	return true
}

// Проверить симметричность матрицы
// Если матрица не симметрична, то выдать ошибку
func (mat Matrix) checkSimmetry() {
	if !mat.IsSimmetric() {
		check(errors.New("Матрица не симметрична"))
	}
}

// Чтение матрицы из файла
func (mat *Matrix) Read(filepath string) {
	// Открываем файл
	b, err := ioutil.ReadFile(filepath)
	check(err)

	// Построчно делим текст
	lines := strings.Split(string(b), "\r\n")
	
	var Rows, Columns int
	Rows = len(lines)

	// Создаём срез срезов и задаём ему размер (количество строк)
	array := make([][]float64, Rows)

	// Проходимся по каждой строке
	for i, line := range lines {
		// Делим каждую строку на отдельные числовые строки
		string_nums := strings.Split(string(line), " ")
		Columns = len(string_nums)

		// Задаём размер срезу (количество столбцов)
		array[i] = make([]float64, Columns)

		// Проходимся по каждой строке с числом
		for j, number := range string_nums {

			// Конвертируем строку в число
			array[i][j], err = strconv.ParseFloat(number, 64)
			check(err)
		}
	}
	mat.Array = array
	mat.Row_count = Rows
	mat.Column_count = Columns
}

// Чтение вектора из файла
func (vec *Vector) Read(filepath string) {
	b, err := ioutil.ReadFile(filepath)
	check(err)

	// Считываем каждое число построчно
	values := strings.Split(string(b), "\r\n")

	size := len(values)
	// Создаём срез для чисел
	array := make([]float64, size)

	// Преобразуем каждую строку в число
	for i, value := range values {
		array[i], err = strconv.ParseFloat(value, 64)
		check(err)
	}
	vec.Array = array
	vec.Size = size
}

// Вывод содержимого матрицы
func (mat Matrix) Write() {
	fmt.Println("[")
	for i := range mat.Array {
		fmt.Print("\t[ ")
		for _, value := range mat.Array[i] {
			fmt.Print(value, " ")
		}
		fmt.Println("]")
	}
	fmt.Println("]")
}

// Вывод содержимого вектора
func (vec Vector) Write() {
	fmt.Println("[")
	for _, v := range vec.Array {
		fmt.Printf("\t%.4f\n", v)
	}
	fmt.Println("]")
}

// Транспонирование матрицы
func (mat *Matrix) Transpose() {
	new_row_count := mat.Column_count
	new_column_count := mat.Row_count

	// Создание временной матрицы для хранения промежуточных данных
	temp_array := make([][]float64, new_row_count)
	for i := range temp_array {
		temp_array[i] = make([]float64, new_column_count)  
	}

	for i := 0; i < new_row_count; i++ {
		for j := 0; j < new_column_count; j++ {
			temp_array[i][j] = mat.Array[j][i]
		}
	}
	mat.Array = temp_array
	mat.Row_count = new_row_count
	mat.Column_count = new_column_count
}

// Умножение матриц
func (first *Matrix) Mul(second Matrix) {
	if first.Row_count != second.Column_count {
		check(errors.New("Количество строк первой матрицы и количество столбцов второй матрицы не совпадают!"))
	}
	// Временная матрица
	result := InitMatrix()
	result.Array = make([][]float64, first.Row_count)
	for i := range result.Array {
		result.Array[i] = make([]float64, second.Column_count)
	}
	for i := 0; i < first.Row_count; i++ {
		for j := 0; j < second.Column_count; j++ {
			var accum float64
			for k := 0; k < second.Row_count; k++ {
				accum += first.Array[i][k] * second.Array[k][j]
			}
			result.Array[i][j] = accum
		}
	}
	first = &result
}

func (vec Vector) MulScalar(scalar float64) (result Vector) {
	result.New(vec.Size)
	for i := range result.Array {
		result.Array[i] = vec.Array[i] * scalar
	}
	return
}

func (first *Vector) Add(second Vector) {
	if first.Size != second.Size {
		check(errors.New("Длины векторов не совпадают!"))
	}
	for i := range first.Array {
		first.Array[i] += second.Array[i]
	}
}

// func (vec *Vector) Normalize() {
// 	// Вычисляем длину вектора
// 	vec_length := 0.0
// 	for _, v := range vec.Array {
// 		vec_length += v * v
// 	}
// 	vec_length = math.Sqrt(vec_length)
// 	// Делим каждый элемент на длину вектора
// 	for i := range vec.Array {
// 		vec.Array[i] /= vec_length
// 	}
// }

func (vec Vector) getSum(first_elements int) (sum float64) {
	for i := 0; i < first_elements; i++ {
		sum += vec.Array[i]
	}
	return 
}

func (vec Vector) Sum() (sum float64) {
	for _, v := range vec.Array {
		sum += v
	}
	return
} 


// Преобразование столбцов матрицы в векторы
func (mat Matrix) ConvertToVec() (vectors []Vector) {
	vectors = make([]Vector, mat.Column_count)
	for j := 0; j < mat.Column_count; j++ {
		vectors[j].New(mat.Row_count)
		for i := 0; i < mat.Row_count; i++ {
			vectors[j].Array[i] = mat.Array[i][j]
		} 
	}
	return
}

// Преобразование среза векторов в матрицу
func ConvertToMat(vectors []Vector) (mat Matrix) {
	// Будем считать, что векторы представляют собой столбцы матрицы
	columns := len(vectors)
	rows := vectors[0].Size 
	mat.New(rows, columns)
	for i := range mat.Array {
		for j := range mat.Array[i] {
			mat.Array[i][j] = vectors[j].Array[i]
		}
	}
	return
}

// Главная диагональ матрицы
func (mat Matrix) GetDiagonal() (diagonal Vector) {
	mat.checkSquareness()
	diagonal.New(mat.Column_count)
	for i := range mat.Array {
		diagonal.Array[i] = mat.Array[i][i]
	}
	return
}