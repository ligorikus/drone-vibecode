package services

import (
	"context"
	"log/slog"
	"time"

	"drone/internal/config"
	"drone/internal/models"
	"drone/internal/utils"
)

// DroneController контролирует поведение дронов
type DroneController struct {
	config *config.Config
	logger *slog.Logger
}

// NewDroneController создаёт новый контроллер дронов
func NewDroneController(cfg *config.Config) *DroneController {
	return &DroneController{
		config: cfg,
	}
}

// UpdatePosition обновляет позицию дочернего дрона
func (dc *DroneController) UpdatePosition(ctx context.Context, child *models.ChildDrone, parent *models.Drone) {
	// Получаем позицию родителя
	parentPos := parent.GetPosition()

	// Получаем текущую позицию дрона
	currentPos := child.GetPosition()

	// Рассчитываем целевую позицию в формировании (сфера вокруг главного дрона)
	// Используем среднюю дистанцию как радиус сферы
	avgRadius := (dc.config.MinDistance + dc.config.MaxDistance) / 2
	// Используем modulo для корректного индекса на сфере
	droneIndex := child.ID % dc.config.DroneCount
	targetPos := models.CalculateFormationTarget(parentPos, droneIndex, dc.config.DroneCount, avgRadius)

	// Вычисляем направление к целевой позиции
	direction := utils.NewVector3D(
		targetPos.X-currentPos.X,
		targetPos.Y-currentPos.Y,
		targetPos.Z-currentPos.Z,
	)

	// Нормализуем и применяем направление
	if direction.Length() > 0 {
		direction = direction.Normalize()
		child.ApplyDirection(direction)

		// Ускоряем дрон к цели
		child.Accelerate(0.016) // ~60 FPS

		// Обновляем позицию
		child.Update(0.016)

		// Ограничиваем высоту (земля)
		child.ClampY(0)
	}

	if dc.config.Debug {
		dc.logger.Debug("Позиция обновлена",
			"drone_id", child.ID,
			"current_position", currentPos,
			"target_position", targetPos,
		)
	}
}

// SetLogger устанавливает логгер для контроллера
func (dc *DroneController) SetLogger(logger *slog.Logger) {
	dc.logger = logger
}

// CalculateFormationDirection вычисляет направление для дочернего дрона к его позиции в формировании
func (dc *DroneController) CalculateFormationDirection(child *models.ChildDrone, parent *models.Drone) *utils.Vector3D {
	parentPos := parent.GetPosition()
	currentPos := child.GetPosition()

	avgRadius := (dc.config.MinDistance + dc.config.MaxDistance) / 2
	droneIndex := child.ID % dc.config.DroneCount
	targetPos := models.CalculateFormationTarget(parentPos, droneIndex, dc.config.DroneCount, avgRadius)

	direction := utils.NewVector3D(
		targetPos.X - currentPos.X,
		targetPos.Y - currentPos.Y,
		targetPos.Z - currentPos.Z,
	)

	if direction.Length() > 0 {
		return direction.Normalize()
	}
	return utils.Zero()
}

// MoveDroneTowards перемещает дрона к цели с заданной скоростью
func (dc *DroneController) MoveDroneTowards(child *models.ChildDrone, target *utils.Vector3D, deltaTime float64) {
	currentPos := child.GetPosition()
	
	direction := utils.NewVector3D(
		target.X - currentPos.X,
		target.Y - currentPos.Y,
		target.Z - currentPos.Z,
	)
	
	distance := direction.Length()
	
	// Если близко к цели, останавливаемся
	if distance < 0.1 {
		child.SetVelocity(utils.Zero())
		return
	}
	
	direction = direction.Normalize()
	child.ApplyDirection(direction)
	child.Accelerate(deltaTime)
	child.Update(deltaTime)
	child.ClampY(0)
}

// KeepFormation поддерживает формирование вокруг главного дрона
func (dc *DroneController) KeepFormation(ctx context.Context, child *models.ChildDrone, parent *models.Drone) {
	ticker := dc.config.UpdateInterval
	
	select {
	case <-ctx.Done():
		return
	case <-time.After(ticker):
		dc.UpdatePosition(ctx, child, parent)
	}
}
