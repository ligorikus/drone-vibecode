package models

import (
	"math"
	"math/rand"

	"drone/internal/utils"
)

// FormationConfig конфигурация формирования
type FormationConfig struct {
	MinDistance       float64
	MaxDistance       float64
	MovementVariation float64
	SmoothingFactor   float64
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
