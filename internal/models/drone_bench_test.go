package models

import (
	"math/rand"
	"testing"
	"time"

	"drone/internal/utils"
)

func BenchmarkCalculateSphericalPosition(b *testing.B) {
	center := utils.NewVector3D(0, 0, 0)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CalculateSphericalPosition(center, 5, 10, rng)
	}
}

func BenchmarkCalculateSpherePoint(b *testing.B) {
	center := utils.NewVector3D(0, 0, 0)
	radius := 7.5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CalculateSpherePoint(center, radius, i%100, 100)
	}
}

func BenchmarkCalculateTruncatedSpherePosition(b *testing.B) {
	center := utils.NewVector3D(0, 10, 0)
	radius := 10.0
	groundLevel := 0.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CalculateTruncatedSpherePosition(center, radius, groundLevel, i%100, 100)
	}
}

func BenchmarkCalculateAdaptiveFormationTarget(b *testing.B) {
	center := utils.NewVector3D(0, 15, 0)
	radius := 10.0
	groundLevel := 0.0

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CalculateAdaptiveFormationTarget(center, i%100, 100, radius, groundLevel)
	}
}

func BenchmarkSmoothMove(b *testing.B) {
	current := utils.NewVector3D(0, 0, 0)
	target := utils.NewVector3D(10, 10, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SmoothMove(current, target, 0.2)
	}
}

func BenchmarkGetFormationStats(b *testing.B) {
	center := utils.NewVector3D(0, 0, 0)
	children := make([]*ChildDrone, 100)
	for i := 0; i < 100; i++ {
		children[i] = NewChildDrone(i, utils.NewVector3D(float64(i), 0, 0))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetFormationStats(center, children)
	}
}

func BenchmarkDroneGetPosition(b *testing.B) {
	drone := NewDrone(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = drone.GetPosition()
	}
}

func BenchmarkDroneSetPosition(b *testing.B) {
	drone := NewDrone(0)
	pos := utils.NewVector3D(1, 2, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		drone.SetPosition(pos)
	}
}

func BenchmarkChildDroneGetPosition(b *testing.B) {
	child := NewChildDrone(0, utils.NewVector3D(0, 0, 0))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = child.GetPosition()
	}
}

func BenchmarkChildDroneSetPosition(b *testing.B) {
	child := NewChildDrone(0, utils.NewVector3D(0, 0, 0))
	pos := utils.NewVector3D(1, 2, 3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		child.SetPosition(pos)
	}
}

func BenchmarkDroneAddChild(b *testing.B) {
	drone := NewDrone(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		child := NewChildDrone(i, utils.NewVector3D(float64(i), 0, 0))
		drone.AddChild(child)
	}
}

func BenchmarkDroneGetChildren(b *testing.B) {
	drone := NewDrone(0)
	for i := 0; i < 100; i++ {
		child := NewChildDrone(i, utils.NewVector3D(float64(i), 0, 0))
		drone.AddChild(child)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = drone.GetChildren()
	}
}

func BenchmarkDroneRangeChildren(b *testing.B) {
	drone := NewDrone(0)
	for i := 0; i < 100; i++ {
		child := NewChildDrone(i, utils.NewVector3D(float64(i), 0, 0))
		drone.AddChild(child)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		drone.RangeChildren(func(*ChildDrone) bool {
			return true
		})
	}
}

func BenchmarkDronePhysics(b *testing.B) {
	drone := NewDrone(0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		drone.Accelerate(0.016)
		drone.Update(0.016)
	}
}

func BenchmarkChildDronePhysics(b *testing.B) {
	child := NewChildDrone(0, utils.NewVector3D(0, 0, 0))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		child.Accelerate(0.016)
		child.Update(0.016)
	}
}
