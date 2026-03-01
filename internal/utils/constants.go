package utils

import "math"

const (
	// DefaultDeltaTime стандартный шаг времени для физики (16мс ~ 60 FPS)
	DefaultDeltaTime = 0.016

	// Sqrt2 квадратный корень из 2 для нормализации диагоналей
	Sqrt2 = 1.41421356237

	// FloatEpsilon погрешность для сравнения float
	FloatEpsilon = 1e-6

	// DefaultFPS стандартная частота кадров
	DefaultFPS = 60

	// MinFPS минимальная частота кадров
	MinFPS = 30
)

// IsFloatEqual сравнивает два float с погрешностью
func IsFloatEqual(a, b float64) bool {
	return math.Abs(a-b) < FloatEpsilon
}

// IsFloatZero проверяет, равен ли float нулю с погрешностью
func IsFloatZero(f float64) bool {
	return math.Abs(f) < FloatEpsilon
}
