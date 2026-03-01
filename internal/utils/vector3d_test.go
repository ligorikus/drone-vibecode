package utils

import (
	"math"
	"testing"
)

func TestNewVector3D(t *testing.T) {
	v := NewVector3D(1, 2, 3)
	
	if v.X != 1 {
		t.Errorf("Expected X=1, got %f", v.X)
	}
	if v.Y != 2 {
		t.Errorf("Expected Y=2, got %f", v.Y)
	}
	if v.Z != 3 {
		t.Errorf("Expected Z=3, got %f", v.Z)
	}
}

func TestVector3DAdd(t *testing.T) {
	v1 := NewVector3D(1, 2, 3)
	v2 := NewVector3D(4, 5, 6)
	
	result := v1.Add(v2)
	
	if result.X != 5 {
		t.Errorf("Expected X=5, got %f", result.X)
	}
	if result.Y != 7 {
		t.Errorf("Expected Y=7, got %f", result.Y)
	}
	if result.Z != 9 {
		t.Errorf("Expected Z=9, got %f", result.Z)
	}
	
	// Проверяем, что оригиналы не изменились
	if v1.X != 1 || v2.X != 4 {
		t.Error("Original vectors were modified")
	}
}

func TestVector3DSubtract(t *testing.T) {
	v1 := NewVector3D(5, 7, 9)
	v2 := NewVector3D(1, 2, 3)
	
	result := v1.Subtract(v2)
	
	if result.X != 4 {
		t.Errorf("Expected X=4, got %f", result.X)
	}
	if result.Y != 5 {
		t.Errorf("Expected Y=5, got %f", result.Y)
	}
	if result.Z != 6 {
		t.Errorf("Expected Z=6, got %f", result.Z)
	}
}

func TestVector3DMultiply(t *testing.T) {
	v := NewVector3D(1, 2, 3)
	
	result := v.Multiply(2)
	
	if result.X != 2 {
		t.Errorf("Expected X=2, got %f", result.X)
	}
	if result.Y != 4 {
		t.Errorf("Expected Y=4, got %f", result.Y)
	}
	if result.Z != 6 {
		t.Errorf("Expected Z=6, got %f", result.Z)
	}
	
	// Проверяем, что оригинал не изменился
	if v.X != 1 {
		t.Error("Original vector was modified")
	}
}

func TestVector3DLength(t *testing.T) {
	v := NewVector3D(3, 4, 0)
	
	length := v.Length()
	
	expected := 5.0
	if math.Abs(length-expected) > 1e-9 {
		t.Errorf("Expected length=%f, got %f", expected, length)
	}
}

func TestVector3DLengthZero(t *testing.T) {
	v := NewVector3D(0, 0, 0)
	
	length := v.Length()
	
	if length != 0 {
		t.Errorf("Expected length=0, got %f", length)
	}
}

func TestVector3DNormalize(t *testing.T) {
	v := NewVector3D(3, 0, 0)
	
	normalized := v.Normalize()
	
	expected := 1.0
	if math.Abs(normalized.Length()-expected) > 1e-9 {
		t.Errorf("Expected normalized length=%f, got %f", expected, normalized.Length())
	}
	
	if normalized.X != 1 {
		t.Errorf("Expected X=1, got %f", normalized.X)
	}
	
	// Проверяем, что оригинал не изменился
	if v.X != 3 {
		t.Error("Original vector was modified")
	}
}

func TestVector3DNormalizeZero(t *testing.T) {
	v := NewVector3D(0, 0, 0)
	
	normalized := v.Normalize()
	
	if normalized.X != 0 || normalized.Y != 0 || normalized.Z != 0 {
		t.Error("Expected zero vector for zero-length input")
	}
}

func TestVector3DDistance(t *testing.T) {
	v1 := NewVector3D(0, 0, 0)
	v2 := NewVector3D(3, 4, 0)
	
	distance := v1.Distance(v2)
	
	expected := 5.0
	if math.Abs(distance-expected) > 1e-9 {
		t.Errorf("Expected distance=%f, got %f", expected, distance)
	}
}

func TestVector3DDistanceSame(t *testing.T) {
	v := NewVector3D(1, 2, 3)
	
	distance := v.Distance(v)
	
	if distance != 0 {
		t.Errorf("Expected distance=0, got %f", distance)
	}
}

func TestVector3DCopy(t *testing.T) {
	v := NewVector3D(1, 2, 3)
	
	copy := v.Copy()
	
	if copy.X != v.X || copy.Y != v.Y || copy.Z != v.Z {
		t.Error("Copy does not match original")
	}
	
	// Изменяем копию
	copy.X = 100
	
	// Проверяем, что оригинал не изменился
	if v.X != 1 {
		t.Error("Original vector was modified when copy was changed")
	}
}

func TestVector3DEqual(t *testing.T) {
	v1 := NewVector3D(1, 2, 3)
	v2 := NewVector3D(1, 2, 3)
	v3 := NewVector3D(1, 2, 4)
	
	if !v1.Equal(v2) {
		t.Error("Expected vectors to be equal")
	}
	
	if v1.Equal(v3) {
		t.Error("Expected vectors to be not equal")
	}
}

func TestVector3DDot(t *testing.T) {
	v1 := NewVector3D(1, 2, 3)
	v2 := NewVector3D(4, 5, 6)
	
	dot := v1.Dot(v2)
	
	expected := 32.0 // 1*4 + 2*5 + 3*6
	if math.Abs(dot-expected) > 1e-9 {
		t.Errorf("Expected dot=%f, got %f", expected, dot)
	}
}

func TestVector3DCross(t *testing.T) {
	v1 := NewVector3D(1, 0, 0)
	v2 := NewVector3D(0, 1, 0)
	
	cross := v1.Cross(v2)
	
	if cross.X != 0 || cross.Y != 0 || cross.Z != 1 {
		t.Errorf("Expected cross product (0, 0, 1), got (%f, %f, %f)", cross.X, cross.Y, cross.Z)
	}
}
