package services

import (
	"drone/internal/models"
	"drone/internal/utils"
)

const (
	// MovementSpeed скорость движения дрона
	MovementSpeed = 15.0
	// VerticalSpeed скорость вертикального движения
	VerticalSpeed = 10.0
	// GroundLevel уровень земли
	GroundLevel = 0.0
)

// InputState представляет состояние ввода с клавиатуры
type InputState struct {
	Forward  bool // W
	Backward bool // S
	Left     bool // A
	Right    bool // D
	Up       bool // Space
	Down     bool // Shift
}

// MainDroneController контроллер для управления главным дроном
type MainDroneController struct {
	input InputState
}

// NewMainDroneController создаёт новый контроллер главного дрона
func NewMainDroneController() *MainDroneController {
	return &MainDroneController{
		input: InputState{},
	}
}

// SetInput устанавливает состояние ввода
func (c *MainDroneController) SetInput(input InputState) {
	c.input = input
}

// Update обновляет состояние главного дрона на основе ввода
func (c *MainDroneController) Update(drone *models.Drone, deltaTime float64) {
	direction := utils.Zero()

	// Горизонтальное движение
	if c.input.Forward {
		direction.Z -= 1
	}
	if c.input.Backward {
		direction.Z += 1
	}
	if c.input.Left {
		direction.X -= 1
	}
	if c.input.Right {
		direction.X += 1
	}

	// Вертикальное движение
	if c.input.Up {
		direction.Y += 1
	}
	if c.input.Down {
		direction.Y -= 1
	}

	// Нормализуем горизонтальную составляющую для диагонального движения
	horizontalLen := utils.NewVector3D(direction.X, 0, direction.Z).Length()
	if horizontalLen > 0 {
		direction.X /= horizontalLen * 1.414 // sqrt(2) для диагоналей
		direction.Z /= horizontalLen * 1.414
	}

	// Применяем направление к дрону
	drone.ApplyDirection(direction)

	// Ускоряем дрон
	drone.Accelerate(deltaTime)

	// Обновляем позицию
	drone.Update(deltaTime)

	// Ограничиваем Y координату (земля)
	drone.ClampY(GroundLevel)
}

// IsMoving возвращает true, если дрон движется
func (c *MainDroneController) IsMoving() bool {
	return c.input.Forward || c.input.Backward || c.input.Left || c.input.Right || c.input.Up || c.input.Down
}
