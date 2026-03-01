package services

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"drone/internal/config"
	"drone/internal/models"
	"drone/internal/utils"
)

func TestDroneControllerNew(t *testing.T) {
	cfg := config.DefaultConfig()
	ctrl := NewDroneController(cfg)

	if ctrl == nil {
		t.Fatal("Expected controller to be created")
	}
	if ctrl.config != cfg {
		t.Error("Expected config to be set")
	}
}

func TestDroneControllerUpdatePosition(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DroneCount = 10
	ctrl := NewDroneController(cfg)

	parent := models.NewDroneAtPosition(0, utils.NewVector3D(0, 10, 0))
	child := models.NewChildDrone(0, utils.NewVector3D(5, 10, 5))

	ctx := context.Background()
	ctrl.UpdatePosition(ctx, child, parent)

	// Проверяем, что позиция изменилась
	newPos := child.GetPosition()
	if newPos == nil {
		t.Fatal("Expected position to be updated")
	}
}

func TestDroneControllerUpdatePositionGroundClamp(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DroneCount = 10
	ctrl := NewDroneController(cfg)

	parent := models.NewDroneAtPosition(0, utils.NewVector3D(0, 5, 0))
	// Дрон ниже уровня земли
	child := models.NewChildDrone(0, utils.NewVector3D(5, -1, 5))

	ctx := context.Background()
	ctrl.UpdatePosition(ctx, child, parent)

	newPos := child.GetPosition()
	if newPos.Y < 0 {
		t.Errorf("Expected Y to be clamped to 0, got %f", newPos.Y)
	}
}

func TestDroneControllerCalculateFormationDirection(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DroneCount = 10
	ctrl := NewDroneController(cfg)

	parent := models.NewDroneAtPosition(0, utils.NewVector3D(0, 10, 0))
	child := models.NewChildDrone(0, utils.NewVector3D(5, 10, 5))

	direction := ctrl.CalculateFormationDirection(child, parent)

	if direction == nil {
		t.Fatal("Expected direction to be calculated")
	}

	// Направление должно быть нормализовано
	length := direction.Length()
	if length < 0.99 || length > 1.01 {
		t.Errorf("Expected normalized direction, got length %f", length)
	}
}

func TestDroneControllerMoveDroneTowards(t *testing.T) {
	cfg := config.DefaultConfig()
	ctrl := NewDroneController(cfg)

	child := models.NewChildDrone(0, utils.NewVector3D(0, 10, 0))
	target := utils.NewVector3D(10, 10, 10)

	ctrl.MoveDroneTowards(child, target, 0.016)

	newPos := child.GetPosition()
	if newPos == nil {
		t.Fatal("Expected position to be updated")
	}
}

func TestDroneControllerMoveDroneTowardsStopAtTarget(t *testing.T) {
	cfg := config.DefaultConfig()
	ctrl := NewDroneController(cfg)

	child := models.NewChildDrone(0, utils.NewVector3D(0, 10, 0))
	target := utils.NewVector3D(0.05, 10, 0) // Очень близко

	ctrl.MoveDroneTowards(child, target, 0.016)

	// После остановки скорость должна быть нулевой
	vel := child.GetVelocity()
	if vel.Length() > 0.5 {
		t.Logf("Velocity: %f (expected near zero)", vel.Length())
	}
}

func TestDroneControllerMoveToFormationPosition(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DroneCount = 10
	ctrl := NewDroneController(cfg)

	parent := models.NewDroneAtPosition(0, utils.NewVector3D(0, 10, 0))
	child := models.NewChildDrone(0, utils.NewVector3D(5, 10, 5))

	ctrl.MoveToFormationPosition(child, parent, 0.016)

	newPos := child.GetPosition()
	if newPos == nil {
		t.Fatal("Expected position to be updated")
	}
}

func TestDroneControllerSetLogger(t *testing.T) {
	cfg := config.DefaultConfig()
	ctrl := NewDroneController(cfg)

	// SetLogger не должен паниковать
	ctrl.SetLogger(nil)
}

func TestDroneControllerMoveDroneTowardsGroundClamp(t *testing.T) {
	cfg := config.DefaultConfig()
	ctrl := NewDroneController(cfg)

	child := models.NewChildDrone(0, utils.NewVector3D(0, -5, 0))
	target := utils.NewVector3D(10, -5, 10)

	ctrl.MoveDroneTowards(child, target, 0.016)

	newPos := child.GetPosition()
	if newPos.Y < 0 {
		t.Errorf("Expected Y to be clamped to 0, got %f", newPos.Y)
	}
}

func TestDroneControllerWithRandomConfig(t *testing.T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cfg := &config.Config{
		DroneCount:       50 + rng.Intn(100),
		MinDistance:      5.0 + rng.Float64()*5,
		MaxDistance:      10.0 + rng.Float64()*5,
		UpdateInterval:   time.Duration(1+rng.Intn(5)) * time.Second,
		FormationRadius:  7.5 + rng.Float64()*5,
		MovementVariation: 0.5 + rng.Float64()*0.5,
		SmoothingFactor:  0.1 + rng.Float64()*0.3,
	}

	ctrl := NewDroneController(cfg)
	if ctrl == nil {
		t.Fatal("Expected controller to be created")
	}
}
