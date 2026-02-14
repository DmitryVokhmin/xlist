package xlist

import (
	"math/rand"
	"slices"
	"testing"
)

// Тестовая структура для бенчмарков
type benchStruct struct {
	Num int
	Str string
}

const (
	benchSize = 100000
)

// Генерация случайной строки
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// Генерация тестовых данных для структур
func generateBenchStructs(size int) []benchStruct {
	data := make([]benchStruct, size)
	for i := 0; i < size; i++ {
		data[i] = benchStruct{
			Num: rand.Intn(1000000),
			Str: randomString(10),
		}
	}
	return data
}

// Генерация тестовых данных для целых чисел
func generateBenchInts(size int) []int {
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Intn(1000000)
	}
	return data
}

// Benchmark сортировки структур с XList
func BenchmarkXListSort_Struct(b *testing.B) {
	// Генерируем данные один раз
	sourceData := generateBenchStructs(benchSize)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Создаем копию данных для каждой итерации
		xlist := New[benchStruct]()
		xlist.Append(sourceData...)
		b.StartTimer()

		// Сортировка только
		xlist.Sort(func(a, b benchStruct) bool {
			return a.Num < b.Num
		})
	}
}

// Benchmark сортировки структур со стандартным slice
func BenchmarkSliceSort_Struct(b *testing.B) {
	// Генерируем данные один раз
	sourceData := generateBenchStructs(benchSize)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Создаем копию данных для каждой итерации
		data := make([]benchStruct, len(sourceData))
		copy(data, sourceData)
		b.StartTimer()

		// Сортировка только
		slices.SortFunc(data, func(a, b benchStruct) int {
			if a.Num < b.Num {
				return -1
			}
			if a.Num > b.Num {
				return 1
			}
			return 0
		})
	}
}

// Benchmark сортировки целых чисел с XList
func BenchmarkXListSort_Int(b *testing.B) {
	// Генерируем данные один раз
	sourceData := generateBenchInts(benchSize)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Создаем копию данных для каждой итерации
		xlist := New[int]()
		xlist.Append(sourceData...)
		b.StartTimer()

		// Сортировка только
		xlist.Sort(func(a, b int) bool {
			return a < b
		})
	}
}

// Benchmark сортировки целых чисел со стандартным slice
func BenchmarkSliceSort_Int(b *testing.B) {
	// Генерируем данные один раз
	sourceData := generateBenchInts(benchSize)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Создаем копию данных для каждой итерации
		data := make([]int, len(sourceData))
		copy(data, sourceData)
		b.StartTimer()

		// Сортировка только
		slices.Sort(data)
	}
}
