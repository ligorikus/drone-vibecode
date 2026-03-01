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

func TestCalculateTruncatedSpherePosition_HighAboveGround(t *testing.T) {
	center := utils.NewVector3D(0, 50, 0) // Высоко над землей
	radius := 10.0
	groundLevel := 0.0

	pos := CalculateTruncatedSpherePosition(center, radius, groundLevel, 0, 10)

	if pos == nil {
		t.Fatal("Position is nil")
	}

	// Проверяем, что позиция на сфере (расстояние от центра = радиус)
	distance := center.Distance(pos)
	if math.Abs(distance-radius) > 0.01 {
		t.Errorf("Distance from center should be %f, got %f", radius, distance)
	}
}

func TestCalculateTruncatedSpherePosition_NearGround(t *testing.T) {
	center := utils.NewVector3D(0, 5, 0) // Близко к земле
	radius := 10.0
	groundLevel := 0.0

	// Создаем множество точек для проверки
	allAboveGround := true
	for i := 0; i < 100; i++ {
		pos := CalculateTruncatedSpherePosition(center, radius, groundLevel, i, 100)
		if pos.Y < groundLevel-0.01 { // Небольшой допуск для погрешности
			allAboveGround = false
			t.Errorf("Drone %d is below ground: Y=%f", i, pos.Y)
		}
	}

	if !allAboveGround {
		t.Error("Some drones are below ground level")
	}
}

func TestCalculateTruncatedSpherePosition_AtGroundLevel(t *testing.T) {
	center := utils.NewVector3D(0, 0, 0) // На уровне земли
	radius := 10.0
	groundLevel := 0.0

	// Все точки должны быть в верхней полусфере
	for i := 0; i < 100; i++ {
		pos := CalculateTruncatedSpherePosition(center, radius, groundLevel, i, 100)
		if pos.Y < groundLevel-0.01 {
			t.Errorf("Drone %d is below ground: Y=%f", i, pos.Y)
		}
		// Расстояние от центра должно быть близко к радиусу
		distance := center.Distance(pos)
		if math.Abs(distance-radius) > 0.1 {
			t.Errorf("Distance from center should be ~%f, got %f", radius, distance)
		}
	}
}

func TestCalculateTruncatedSpherePosition_BelowGround(t *testing.T) {
	center := utils.NewVector3D(0, -5, 0) // Ниже уровня земли
	radius := 10.0
	groundLevel := 0.0

	// Все точки все равно должны быть над землей (на верхней части сферы)
	for i := 0; i < 50; i++ {
		pos := CalculateTruncatedSpherePosition(center, radius, groundLevel, i, 50)
		if pos.Y < groundLevel-0.01 {
			t.Errorf("Drone %d is below ground: Y=%f", i, pos.Y)
		}
	}
}

func TestCalculateTruncatedSpherePosition_UniformDistribution(t *testing.T) {
	center := utils.NewVector3D(0, 10, 0)
	radius := 10.0
	groundLevel := 0.0
	total := 100

	// Проверяем, что точки распределены по разным углам
	positions := make([]*utils.Vector3D, total)
	for i := 0; i < total; i++ {
		positions[i] = CalculateTruncatedSpherePosition(center, radius, groundLevel, i, total)
	}

	// Проверяем, что нет дубликатов позиций
	for i := 0; i < total; i++ {
		for j := i + 1; j < total; j++ {
			if positions[i].Distance(positions[j]) < 0.1 {
				t.Errorf("Positions %d and %d are too close: %v and %v", i, j, positions[i], positions[j])
			}
		}
	}
}

func TestCalculateAdaptiveFormationTarget_HighAltitude(t *testing.T) {
	center := utils.NewVector3D(0, 50, 0) // Высоко
	radius := 10.0
	groundLevel := 0.0

	pos := CalculateAdaptiveFormationTarget(center, 0, 100, radius, groundLevel)

	if pos == nil {
		t.Fatal("Position is nil")
	}

	// На большой высоте должна быть обычная сфера
	distance := center.Distance(pos)
	if math.Abs(distance-radius) > 0.1 {
		t.Errorf("At high altitude, should be full sphere: distance=%f, expected=%f", distance, radius)
	}
}

func TestCalculateAdaptiveFormationTarget_LowAltitude(t *testing.T) {
	center := utils.NewVector3D(0, 5, 0) // Низко
	radius := 10.0
	groundLevel := 0.0

	// Все точки должны быть над землей
	for i := 0; i < 50; i++ {
		pos := CalculateAdaptiveFormationTarget(center, i, 50, radius, groundLevel)
		if pos.Y < groundLevel-0.01 {
			t.Errorf("Drone %d is below ground: Y=%f", i, pos.Y)
		}
	}
}

func TestCalculateAdaptiveFormationTarget_Transition(t *testing.T) {
	radius := 10.0
	groundLevel := 0.0

	// Проверяем плавный переход на разных высотах
	heights := []float64{5, 10, 15, 20, 25, 30}
	
	for _, height := range heights {
		center := utils.NewVector3D(0, height, 0)
		allAboveGround := true
		
		for i := 0; i < 20; i++ {
			pos := CalculateAdaptiveFormationTarget(center, i, 20, radius, groundLevel)
			if pos.Y < groundLevel-0.01 {
				allAboveGround = false
			}
		}
		
		if !allAboveGround && height < radius*2 {
			t.Errorf("At height %f, some drones are below ground", height)
		}
	}
}
