package utils

import "math"

// Константы симуляции
const (
	// DefaultDeltaTime стандартный шаг времени для обновлений (16мс ~ 60 FPS)
	DefaultDeltaTime = 0.016

	// FloatEpsilon погрешность для сравнения float
	FloatEpsilon = 1e-6

	// Sqrt2 квадратный корень из 2 для нормализации диагонального движения
	Sqrt2 = 1.4142135623730951
)

// GroundLevel уровень земли (Y координата)
const GroundLevel = 0.0

// IsFloatZero проверяет, равен ли float нулю с учётом погрешности
func IsFloatZero(x float64) bool {
	return math.Abs(x) <= FloatEpsilon
}

// IsFloatEqual проверяет равенство двух float значений с учётом погрешности
func IsFloatEqual(a, b float64) bool {
	return math.Abs(a-b) <= FloatEpsilon
}
