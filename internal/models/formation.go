package models

import (
	"math"
	"math/rand"

	"drone/internal/utils"
)

// FormationType тип формирования дронов
type FormationType int

const (
	// FormationSphere полная сфера
	FormationSphere FormationType = iota
	// FormationTruncatedSphere срезанная сфера (полусфера со стороны земли)
	FormationTruncatedSphere
)

// FormationConfig конфигурация формирования
type FormationConfig struct {
	MinDistance       float64
	MaxDistance       float64
	MovementVariation float64
	SmoothingFactor   float64
	FormationType     FormationType
	GroundLevel       float64 // Уровень земли (Y координата)
}

// CalculateSphericalPosition рассчитывает позицию в сферическом формировании
func CalculateSphericalPosition(center *utils.Vector3D, minDist, maxDist float64) *utils.Vector3D {
	angle1 := rand.Float64() * 2 * math.Pi
	angle2 := rand.Float64() * math.Pi
	distance := minDist + rand.Float64()*(maxDist-minDist)

	x := distance * math.Sin(angle2) * math.Cos(angle1)
	y := distance * math.Sin(angle2) * math.Sin(angle1)
	z := distance * math.Cos(angle2)

	return utils.NewVector3D(
		center.X+x,
		center.Y+y,
		center.Z+z,
	)
}

// CalculateSpherePoint рассчитывает точку на сфере с использованием золотой спирали
// для более равномерного распределения точек
func CalculateSpherePoint(center *utils.Vector3D, radius float64, index, total int) *utils.Vector3D {
	// Золотое сечение для равномерного распределения
	goldenRatio := (1 + math.Sqrt(5)) / 2
	
	// Угол по вертикали
	y := 1 - (float64(index) / float64(total-1)) * 2 // от 1 до -1
	radiusAtY := math.Sqrt(1 - y*y)
	
	// Угол по горизонтали
	theta := 2 * math.Pi * float64(index) / goldenRatio
	
	x := math.Cos(theta) * radiusAtY
	z := math.Sin(theta) * radiusAtY
	
	return utils.NewVector3D(
		center.X+x*radius,
		center.Y+y*radius,
		center.Z+z*radius,
	)
}

// CalculateTargetPosition рассчитывает целевую позицию для поддержания формирования
func CalculateTargetPosition(parentPos *utils.Vector3D, currentPos *utils.Vector3D, config FormationConfig) *utils.Vector3D {
	angle1 := rand.Float64() * 2 * math.Pi
	angle2 := rand.Float64() * math.Pi
	distance := config.MinDistance + rand.Float64()*config.MovementVariation

	targetX := parentPos.X + distance*math.Sin(angle2)*math.Cos(angle1)
	targetY := parentPos.Y + distance*math.Sin(angle2)*math.Sin(angle1)
	targetZ := parentPos.Z + distance*math.Cos(angle2)

	targetPos := utils.NewVector3D(targetX, targetY, targetZ)

	// Плавное движение к целевой позиции
	return SmoothMove(currentPos, targetPos, config.SmoothingFactor)
}

// CalculateFormationTarget рассчитывает целевую позицию для дочернего дрона в формировании
// с использованием фиксированной точки на сфере для каждого дрона
func CalculateFormationTarget(parentPos *utils.Vector3D, droneIndex, totalDrones int, radius float64) *utils.Vector3D {
	return CalculateSpherePoint(parentPos, radius, droneIndex, totalDrones)
}

// CalculateTruncatedSpherePosition рассчитывает позицию на срезанной сфере (полусфере)
// Срезанная сфера формируется над землей, дроны распределяются равномерно по видимой части
func CalculateTruncatedSpherePosition(center *utils.Vector3D, radius float64, groundLevel float64, index, total int) *utils.Vector3D {
	// Определяем, насколько центр сферы близок к земле
	centerHeight := center.Y - groundLevel

	// Если центр выше радиуса - используем полную сферу
	if centerHeight >= radius {
		return CalculateSpherePoint(center, radius, index, total)
	}

	// Золотое сечение для равномерного распределения
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Вычисляем максимальный угол от вертикали (от оси Y)
	// который обеспечивает нахождение точки над землей
	// Точка на сфере: y = center.Y + radius * cos(polarAngle)
	// Требуем: y >= groundLevel
	// => center.Y + radius * cos(polarAngle) >= groundLevel
	// => cos(polarAngle) >= (groundLevel - center.Y) / radius = -centerHeight / radius
	
	// Минимальное значение cos(polarAngle) для нахождения над землей
	minCosPolar := -centerHeight / radius
	
	// Ограничиваем диапазон [-1, 1]
	if minCosPolar < -1 {
		minCosPolar = -1
	}
	if minCosPolar > 1 {
		// Центр настолько глубоко, что вся сфера под землей
		// В этом случае размещаем дроны на самой верхней точке
		return utils.NewVector3D(center.X, center.Y+radius, center.Z)
	}

	// Максимальный полярный угол (от вертикали)
	maxPolarAngle := math.Acos(minCosPolar)

	// Распределяем точки только на видимой части сферы (выше земли)
	// Используем модифицированную золотую спираль
	t := float64(index) / float64(total)

	// Полярный угол: от 0 до maxPolarAngle (только верхняя часть)
	polarAngle := maxPolarAngle * t

	// Угол по горизонтали
	azimuthAngle := 2 * math.Pi * float64(index) / goldenRatio

	// Рассчитываем координаты
	y := math.Cos(polarAngle)
	radiusAtY := math.Sin(polarAngle)

	x := math.Cos(azimuthAngle) * radiusAtY
	z := math.Sin(azimuthAngle) * radiusAtY

	return utils.NewVector3D(
		center.X+x*radius,
		center.Y+y*radius,
		center.Z+z*radius,
	)
}

// CalculateAdaptiveFormationTarget рассчитывает целевую позицию с адаптивным выбором формации
// При приближении к земле плавно переходит от сферы к срезанной сфере
func CalculateAdaptiveFormationTarget(parentPos *utils.Vector3D, droneIndex, totalDrones int, radius float64, groundLevel float64) *utils.Vector3D {
	centerHeight := parentPos.Y - groundLevel

	// Если главный дрон высоко - используем обычную сферу
	if centerHeight >= radius*2 {
		return CalculateSpherePoint(parentPos, radius, droneIndex, totalDrones)
	}

	// Если близко к земле - используем срезанную сферу
	if centerHeight <= radius {
		return CalculateTruncatedSpherePosition(parentPos, radius, groundLevel, droneIndex, totalDrones)
	}

	// Плавный переход между сферой и срезанной сферой
	transitionFactor := (centerHeight - radius) / radius // от 0 до 1
	if transitionFactor > 1 {
		transitionFactor = 1
	}

	// Получаем позиции для обоих режимов
	spherePos := CalculateSpherePoint(parentPos, radius, droneIndex, totalDrones)
	truncatedPos := CalculateTruncatedSpherePosition(parentPos, radius, groundLevel, droneIndex, totalDrones)

	// Интерполируем между ними
	return utils.NewVector3D(
		spherePos.X*transitionFactor+truncatedPos.X*(1-transitionFactor),
		spherePos.Y*transitionFactor+truncatedPos.Y*(1-transitionFactor),
		spherePos.Z*transitionFactor+truncatedPos.Z*(1-transitionFactor),
	)
}

// SmoothMove плавно перемещает текущую позицию к целевой
func SmoothMove(current, target *utils.Vector3D, factor float64) *utils.Vector3D {
	return utils.NewVector3D(
		current.X*(1-factor)+target.X*factor,
		current.Y*(1-factor)+target.Y*factor,
		current.Z*(1-factor)+target.Z*factor,
	)
}

// GetFormationStats возвращает статистику формирования
func GetFormationStats(center *utils.Vector3D, children []*ChildDrone) FormationStats {
	if len(children) == 0 {
		return FormationStats{}
	}

	var totalDist, minDist, maxDist float64
	minDist = math.MaxFloat64

	for _, child := range children {
		pos := child.GetPosition()
		dist := center.Distance(pos)
		totalDist += dist

		if dist < minDist {
			minDist = dist
		}
		if dist > maxDist {
			maxDist = dist
		}
	}

	count := float64(len(children))
	return FormationStats{
		Count:       len(children),
		AvgDistance: totalDist / count,
		MinDistance: minDist,
		MaxDistance: maxDist,
	}
}

// FormationStats статистика формирования
type FormationStats struct {
	Count       int
	AvgDistance float64
	MinDistance float64
	MaxDistance float64
}
