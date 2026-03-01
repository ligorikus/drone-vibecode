package models

import (
	"math"
	"testing"

	"drone/internal/utils"
)

func TestCalculateSphericalPosition(t *testing.T) {
	center := utils.NewVector3D(0, 0, 0)
	
	pos := CalculateSphericalPosition(center, 5, 10)
	
// Проверяем, что позиция не nil
	if pos == nil {
		t.Fatal("Position is nil")
	}
	
// Проверяем, что расстояние в пределах диапазона
	distance := center.Distance(pos)
	if distance < 5 || distance > 10 {
		t.Errorf("Distance out of range [5, 10]: %f", distance)
	}
}

func TestCalculateSphericalPositionNonZeroCenter(t *testing.T) {
	center := utils.NewVector3D(10, 20, 30)
	
	pos := CalculateSphericalPosition(center, 5, 10)
	
// Расстояние от центра должно быть в диапазоне
	distance := center.Distance(pos)
	if distance < 5 || distance > 10 {
		t.Errorf("Distance out of range [5, 10]: %f", distance)
	}
}

func TestSmoothMove(t *testing.T) {
	current := utils.NewVector3D(0, 0, 0)
	target := utils.NewVector3D(10, 10, 10)
	
// При factor=0.5 должна быть половина пути
	result := SmoothMove(current, target, 0.5)
	
	expected := 5.0
	if math.Abs(result.X-expected) > 1e-9 {
		t.Errorf("Expected X=%f, got %f", expected, result.X)
	}
	if math.Abs(result.Y-expected) > 1e-9 {
		t.Errorf("Expected Y=%f, got %f", expected, result.Y)
	}
	if math.Abs(result.Z-expected) > 1e-9 {
		t.Errorf("Expected Z=%f, got %f", expected, result.Z)
	}
}

func TestSmoothMoveFactorZero(t *testing.T) {
	current := utils.NewVector3D(0, 0, 0)
	target := utils.NewVector3D(10, 10, 10)
	
	result := SmoothMove(current, target, 0)
	
// При factor=0 должна остаться текущая позиция
	if !result.Equal(current) {
		t.Errorf("Expected %v, got %v", current, result)
	}
}

func TestSmoothMoveFactorOne(t *testing.T) {
	current := utils.NewVector3D(0, 0, 0)
	target := utils.NewVector3D(10, 10, 10)
	
	result := SmoothMove(current, target, 1)
	
// При factor=1 должна стать целевой позицией
	if !result.Equal(target) {
		t.Errorf("Expected %v, got %v", target, result)
	}
}

func TestCalculateTargetPosition(t *testing.T) {
	parentPos := utils.NewVector3D(0, 0, 0)
	currentPos := utils.NewVector3D(5, 0, 0)
	
	config := FormationConfig{
		MinDistance:       5,
		MaxDistance:       10,
		MovementVariation: 0.5,
		SmoothingFactor:   0.2,
	}
	
	result := CalculateTargetPosition(parentPos, currentPos, config)
	
	if result == nil {
		t.Fatal("Result position is nil")
	}
}

func TestGetFormationStats(t *testing.T) {
	center := utils.NewVector3D(0, 0, 0)
	
	children := []*ChildDrone{
		NewChildDrone(1, utils.NewVector3D(5, 0, 0)),
		NewChildDrone(2, utils.NewVector3D(10, 0, 0)),
		NewChildDrone(3, utils.NewVector3D(7.5, 0, 0)),
	}
	
	stats := GetFormationStats(center, children)
	
	if stats.Count != 3 {
		t.Errorf("Expected count=3, got %d", stats.Count)
	}
	
// Средняя дистанция: (5 + 10 + 7.5) / 3 = 7.5
	expectedAvg := 7.5
	if math.Abs(stats.AvgDistance-expectedAvg) > 1e-9 {
		t.Errorf("Expected avg=%f, got %f", expectedAvg, stats.AvgDistance)
	}
	
// Минимальная дистанция: 5
	if math.Abs(stats.MinDistance-5) > 1e-9 {
		t.Errorf("Expected min=5, got %f", stats.MinDistance)
	}
	
// Максимальная дистанция: 10
	if math.Abs(stats.MaxDistance-10) > 1e-9 {
		t.Errorf("Expected max=10, got %f", stats.MaxDistance)
	}
}

func TestGetFormationStatsEmpty(t *testing.T) {
	center := utils.NewVector3D(0, 0, 0)
	
	stats := GetFormationStats(center, []*ChildDrone{})
	
	if stats.Count != 0 {
		t.Errorf("Expected count=0, got %d", stats.Count)
	}
	if stats.AvgDistance != 0 {
		t.Errorf("Expected avg=0, got %f", stats.AvgDistance)
	}
}

func TestGetFormationStatsSingle(t *testing.T) {
	center := utils.NewVector3D(0, 0, 0)
	children := []*ChildDrone{
		NewChildDrone(1, utils.NewVector3D(8, 0, 0)),
	}
	
	stats := GetFormationStats(center, children)
	
	if stats.Count != 1 {
		t.Errorf("Expected count=1, got %d", stats.Count)
	}
	if math.Abs(stats.AvgDistance-8) > 1e-9 {
		t.Errorf("Expected avg=8, got %f", stats.AvgDistance)
	}
	if math.Abs(stats.MinDistance-8) > 1e-9 {
		t.Errorf("Expected min=8, got %f", stats.MinDistance)
	}
	if math.Abs(stats.MaxDistance-8) > 1e-9 {
		t.Errorf("Expected max=8, got %f", stats.MaxDistance)
	}
}
