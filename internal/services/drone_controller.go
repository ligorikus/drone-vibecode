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

	// Рассчитываем целевую позицию с адаптивной формацией
	// При приближении к земле сфера автоматически превращается в срезанную сферу
	avgRadius := (dc.config.MinDistance + dc.config.MaxDistance) / 2
	droneIndex := child.ID % dc.config.DroneCount

	// Используем адаптивную целевую позицию с учетом уровня земли (groundLevel = 0)
	targetPos := models.CalculateAdaptiveFormationTarget(parentPos, droneIndex, dc.config.DroneCount, avgRadius, models.GroundLevel)

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
		child.Accelerate(utils.DefaultDeltaTime)

		// Обновляем позицию
		child.Update(utils.DefaultDeltaTime)

		// Ограничиваем высоту (земля)
		child.ClampY(models.GroundLevel)
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
	
	// Используем адаптивную целевую позицию
	targetPos := models.CalculateAdaptiveFormationTarget(parentPos, droneIndex, dc.config.DroneCount, avgRadius, 0)

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
	if distance < utils.FloatEpsilon {
		child.SetVelocity(utils.Zero())
		return
	}

	direction = direction.Normalize()
	child.ApplyDirection(direction)
	child.Accelerate(deltaTime)
	child.Update(deltaTime)
	child.ClampY(models.GroundLevel)
}

// MoveToFormationPosition перемещает дрона к позиции в адаптивной формации
func (dc *DroneController) MoveToFormationPosition(child *models.ChildDrone, parent *models.Drone, deltaTime float64) {
	parentPos := parent.GetPosition()
	avgRadius := (dc.config.MinDistance + dc.config.MaxDistance) / 2
	droneIndex := child.ID % dc.config.DroneCount

	targetPos := models.CalculateAdaptiveFormationTarget(parentPos, droneIndex, dc.config.DroneCount, avgRadius, models.GroundLevel)
	dc.MoveDroneTowards(child, targetPos, deltaTime)
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
