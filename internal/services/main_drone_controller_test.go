package services

import (
	"testing"

	"drone/internal/models"
	"drone/internal/utils"
)

func TestInputState(t *testing.T) {
	input := InputState{
		Forward:  true,
		Backward: false,
		Left:     false,
		Right:    false,
		Up:       false,
		Down:     false,
	}

	if !input.Forward {
		t.Error("Expected Forward to be true")
	}
	if input.Backward {
		t.Error("Expected Backward to be false")
	}
}

func TestMainDroneControllerNew(t *testing.T) {
	ctrl := NewMainDroneController()

	if ctrl == nil {
		t.Fatal("Expected controller to be created")
	}
}

func TestMainDroneControllerSetInput(t *testing.T) {
	ctrl := NewMainDroneController()

	input := InputState{
		Forward: true,
		Up:      true,
	}

	ctrl.SetInput(input)

	if !ctrl.input.Forward {
		t.Error("Expected Forward to be true after SetInput")
	}
	if !ctrl.input.Up {
		t.Error("Expected Up to be true after SetInput")
	}
}

func TestMainDroneControllerUpdate(t *testing.T) {
	ctrl := NewMainDroneController()
	drone := models.NewDroneAtPosition(0, utils.NewVector3D(0, 10, 0))

	input := InputState{
		Forward: true,
	}
	ctrl.SetInput(input)

	// Обновляем дрон
	ctrl.Update(drone, 0.016)

	// Проверяем, что позиция изменилась
	newPos := drone.GetPosition()
	if newPos.Z >= 0 {
		t.Errorf("Expected Z to decrease when moving forward, got %f", newPos.Z)
	}
}

func TestMainDroneControllerUpdateHorizontalNormalization(t *testing.T) {
	ctrl := NewMainDroneController()
	drone := models.NewDroneAtPosition(0, utils.NewVector3D(0, 10, 0))

	// Диагональное движение
	input := InputState{
		Forward: true,
		Right:   true,
	}
	ctrl.SetInput(input)

	ctrl.Update(drone, 0.016)

	// Проверяем, что дрон движется
	newPos := drone.GetPosition()
	if newPos.Z >= 0 && newPos.X <= 0 {
		t.Errorf("Expected diagonal movement, got position (%f, %f, %f)", newPos.X, newPos.Y, newPos.Z)
	}
}

func TestMainDroneControllerUpdateVerticalMovement(t *testing.T) {
	ctrl := NewMainDroneController()
	drone := models.NewDroneAtPosition(0, utils.NewVector3D(0, 10, 0))

	input := InputState{
		Up: true,
	}
	ctrl.SetInput(input)

	ctrl.Update(drone, 0.016)

	newPos := drone.GetPosition()
	if newPos.Y <= 10 {
		t.Errorf("Expected Y to increase when moving up, got %f", newPos.Y)
	}
}

func TestMainDroneControllerUpdateDownMovement(t *testing.T) {
	ctrl := NewMainDroneController()
	drone := models.NewDroneAtPosition(0, utils.NewVector3D(0, 10, 0))

	input := InputState{
		Down: true,
	}
	ctrl.SetInput(input)

	ctrl.Update(drone, 0.016)

	newPos := drone.GetPosition()
	if newPos.Y >= 10 {
		t.Errorf("Expected Y to decrease when moving down, got %f", newPos.Y)
	}
}

func TestMainDroneControllerUpdateGroundClamp(t *testing.T) {
	ctrl := NewMainDroneController()
	drone := models.NewDroneAtPosition(0, utils.NewVector3D(0, 1, 0))

	input := InputState{
		Down: true,
	}
	ctrl.SetInput(input)

	// Несколько обновлений для достижения земли
	for i := 0; i < 100; i++ {
		ctrl.Update(drone, 0.016)
	}

	newPos := drone.GetPosition()
	if newPos.Y < utils.GroundLevel-0.01 {
		t.Errorf("Expected Y to be clamped to ground level, got %f", newPos.Y)
	}
}

func TestMainDroneControllerIsMoving(t *testing.T) {
	ctrl := NewMainDroneController()

	if ctrl.IsMoving() {
		t.Error("Expected IsMoving to be false with no input")
	}

	ctrl.SetInput(InputState{Forward: true})
	if !ctrl.IsMoving() {
		t.Error("Expected IsMoving to be true with Forward input")
	}

	ctrl.SetInput(InputState{Down: true})
	if !ctrl.IsMoving() {
		t.Error("Expected IsMoving to be true with Down input")
	}
}

func TestMainDroneControllerUpdateBackward(t *testing.T) {
	ctrl := NewMainDroneController()
	drone := models.NewDroneAtPosition(0, utils.NewVector3D(0, 10, 0))

	input := InputState{
		Backward: true,
	}
	ctrl.SetInput(input)

	ctrl.Update(drone, 0.016)

	newPos := drone.GetPosition()
	if newPos.Z <= 0 {
		t.Errorf("Expected Z to increase when moving backward, got %f", newPos.Z)
	}
}

func TestMainDroneControllerUpdateLeft(t *testing.T) {
	ctrl := NewMainDroneController()
	drone := models.NewDroneAtPosition(0, utils.NewVector3D(0, 10, 0))

	input := InputState{
		Left: true,
	}
	ctrl.SetInput(input)

	ctrl.Update(drone, 0.016)

	newPos := drone.GetPosition()
	if newPos.X >= 0 {
		t.Errorf("Expected X to decrease when moving left, got %f", newPos.X)
	}
}

func TestMainDroneControllerUpdateRight(t *testing.T) {
	ctrl := NewMainDroneController()
	drone := models.NewDroneAtPosition(0, utils.NewVector3D(0, 10, 0))

	input := InputState{
		Right: true,
	}
	ctrl.SetInput(input)

	ctrl.Update(drone, 0.016)

	newPos := drone.GetPosition()
	if newPos.X <= 0 {
		t.Errorf("Expected X to increase when moving right, got %f", newPos.X)
	}
}
