package utils

import (
	"math"
	"math/rand"
)

// Vector3D представляет 3D вектор
type Vector3D struct {
	X, Y, Z float64
}

// NewVector3D создаёт новый 3D вектор
func NewVector3D(x, y, z float64) *Vector3D {
	return &Vector3D{X: x, Y: y, Z: z}
}

// NewRandomVector3D создаёт случайный 3D вектор с X/Z в диапазоне [-maxRange, maxRange] и Y=0
func NewRandomVector3D(maxRange float64) *Vector3D {
	return &Vector3D{
		X: (rand.Float64()*2 - 1) * maxRange,
		Y: 0,
		Z: (rand.Float64()*2 - 1) * maxRange,
	}
}

// Zero создаёт нулевой вектор
func Zero() *Vector3D {
	return &Vector3D{X: 0, Y: 0, Z: 0}
}

// Add складывает два вектора
func (v *Vector3D) Add(other *Vector3D) *Vector3D {
	return &Vector3D{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

// Subtract вычитает другой вектор из этого вектора
func (v *Vector3D) Subtract(other *Vector3D) *Vector3D {
	return &Vector3D{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

// Multiply масштабирует вектор на скаляр
func (v *Vector3D) Multiply(scalar float64) *Vector3D {
	return &Vector3D{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
	}
}

// Length вычисляет длину вектора
func (v *Vector3D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize нормализует вектор к единичной длине
func (v *Vector3D) Normalize() *Vector3D {
	length := v.Length()
	if length == 0 {
		return &Vector3D{}
	}
	return &Vector3D{
		X: v.X / length,
		Y: v.Y / length,
		Z: v.Z / length,
	}
}

// Distance вычисляет расстояние между двумя точками
func (v *Vector3D) Distance(other *Vector3D) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

// Copy создаёт копию вектора
func (v *Vector3D) Copy() *Vector3D {
	return &Vector3D{X: v.X, Y: v.Y, Z: v.Z}
}

// Equal проверяет равенство векторов
func (v *Vector3D) Equal(other *Vector3D) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

// Dot вычисляет скалярное произведение
func (v *Vector3D) Dot(other *Vector3D) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross вычисляет векторное произведение
func (v *Vector3D) Cross(other *Vector3D) *Vector3D {
	return &Vector3D{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// IsFloatEqual проверяет равенство двух float с учётом погрешности
func IsFloatEqual(a, b float64) bool {
	return math.Abs(a-b) < FloatEpsilon
}

// IsFloatZero проверяет, равен ли float нулю с учётом погрешности
func IsFloatZero(a float64) bool {
	return math.Abs(a) < FloatEpsilon
}
