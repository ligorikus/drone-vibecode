package utils

import (
	"math"
	"testing"
)

func BenchmarkVector3DAdd(b *testing.B) {
	v1 := NewVector3D(1, 2, 3)
	v2 := NewVector3D(4, 5, 6)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v1.Add(v2)
	}
}

func BenchmarkVector3DSubtract(b *testing.B) {
	v1 := NewVector3D(1, 2, 3)
	v2 := NewVector3D(4, 5, 6)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v1.Subtract(v2)
	}
}

func BenchmarkVector3DMultiply(b *testing.B) {
	v := NewVector3D(1, 2, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Multiply(2.5)
	}
}

func BenchmarkVector3DLength(b *testing.B) {
	v := NewVector3D(1, 2, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Length()
	}
}

func BenchmarkVector3DNormalize(b *testing.B) {
	v := NewVector3D(1, 2, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Normalize()
	}
}

func BenchmarkVector3DDistance(b *testing.B) {
	v1 := NewVector3D(1, 2, 3)
	v2 := NewVector3D(4, 5, 6)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v1.Distance(v2)
	}
}

func BenchmarkVector3DCopy(b *testing.B) {
	v := NewVector3D(1, 2, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.Copy()
	}
}

func BenchmarkVector3DEqual(b *testing.B) {
	v1 := NewVector3D(1, 2, 3)
	v2 := NewVector3D(1, 2, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v1.Equal(v2)
	}
}

func BenchmarkVector3DDot(b *testing.B) {
	v1 := NewVector3D(1, 2, 3)
	v2 := NewVector3D(4, 5, 6)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v1.Dot(v2)
	}
}

func BenchmarkVector3DCross(b *testing.B) {
	v1 := NewVector3D(1, 2, 3)
	v2 := NewVector3D(4, 5, 6)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v1.Cross(v2)
	}
}

func BenchmarkIsFloatEqual(b *testing.B) {
	a, c := 1.0000001, 1.0000002

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsFloatEqual(a, c)
	}
}

func BenchmarkIsFloatZero(b *testing.B) {
	f := 0.0000001

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsFloatZero(f)
	}
}

func BenchmarkVector3DMath(b *testing.B) {
	v := NewVector3D(3, 4, 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	}
}
