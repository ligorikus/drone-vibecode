package models

import (
	"sync"
	"testing"

	"drone/internal/utils"
)

func TestNewDrone(t *testing.T) {
	drone := NewDrone(1)
	
	if drone.ID != 1 {
		t.Errorf("Expected ID=1, got %d", drone.ID)
	}
	
	pos := drone.GetPosition()
	if pos.X != 0 || pos.Y != 0 || pos.Z != 0 {
		t.Errorf("Expected position (0,0,0), got (%f,%f,%f)", pos.X, pos.Y, pos.Z)
	}
	
	if drone.GetChildCount() != 0 {
		t.Errorf("Expected 0 children, got %d", drone.GetChildCount())
	}
}

func TestNewChildDrone(t *testing.T) {
	pos := utils.NewVector3D(1, 2, 3)
	child := NewChildDrone(1, pos)
	
	if child.ID != 1 {
		t.Errorf("Expected ID=1, got %d", child.ID)
	}
	
	childPos := child.GetPosition()
	if childPos.X != 1 || childPos.Y != 2 || childPos.Z != 3 {
		t.Errorf("Expected position (1,2,3), got (%f,%f,%f)", childPos.X, childPos.Y, childPos.Z)
	}
	
// Проверяем, что позиция скопирована
	pos.X = 100
	childPos = child.GetPosition()
	if childPos.X == 100 {
		t.Error("Child drone position was not copied")
	}
}

func TestDroneAddChild(t *testing.T) {
	drone := NewDrone(0)
	child1 := NewChildDrone(1, utils.NewVector3D(1, 0, 0))
	child2 := NewChildDrone(2, utils.NewVector3D(2, 0, 0))
	
	drone.AddChild(child1)
	drone.AddChild(child2)
	
	if drone.GetChildCount() != 2 {
		t.Errorf("Expected 2 children, got %d", drone.GetChildCount())
	}
}

func TestDroneGetChildren(t *testing.T) {
	drone := NewDrone(0)
	child1 := NewChildDrone(1, utils.NewVector3D(1, 0, 0))
	child2 := NewChildDrone(2, utils.NewVector3D(2, 0, 0))
	
	drone.AddChild(child1)
	drone.AddChild(child2)
	
	children := drone.GetChildren()
	
	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}
	
// Проверяем, что возвращается копия слайса
	children = append(children, NewChildDrone(3, utils.NewVector3D(3, 0, 0)))
	if drone.GetChildCount() != 2 {
		t.Error("GetChildren returned reference instead of copy")
	}
}

func TestDroneSetPosition(t *testing.T) {
	drone := NewDrone(0)
	newPos := utils.NewVector3D(10, 20, 30)
	
	drone.SetPosition(newPos)
	
	pos := drone.GetPosition()
	if pos.X != 10 || pos.Y != 20 || pos.Z != 30 {
		t.Errorf("Expected position (10,20,30), got (%f,%f,%f)", pos.X, pos.Y, pos.Z)
	}
	
// Проверяем, что позиция копируется
	newPos.X = 100
	pos = drone.GetPosition()
	if pos.X == 100 {
		t.Error("Drone position was not copied")
	}
}

func TestDroneConcurrency(t *testing.T) {
	drone := NewDrone(0)
	
	var wg sync.WaitGroup
	
// Одновременно добавляем детей из разных горутин
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			child := NewChildDrone(id, utils.NewVector3D(float64(id), 0, 0))
			drone.AddChild(child)
		}(i)
	}
	
	wg.Wait()
	
	if drone.GetChildCount() != 100 {
		t.Errorf("Expected 100 children, got %d", drone.GetChildCount())
	}
}

func TestChildDroneSetPosition(t *testing.T) {
	child := NewChildDrone(1, utils.NewVector3D(0, 0, 0))
	newPos := utils.NewVector3D(10, 20, 30)
	
	child.SetPosition(newPos)
	
	pos := child.GetPosition()
	if pos.X != 10 || pos.Y != 20 || pos.Z != 30 {
		t.Errorf("Expected position (10,20,30), got (%f,%f,%f)", pos.X, pos.Y, pos.Z)
	}
}

func TestChildDroneConcurrency(t *testing.T) {
	child := NewChildDrone(1, utils.NewVector3D(0, 0, 0))
	
	var wg sync.WaitGroup
	
// Одновременно обновляем позицию из разных горутин
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val float64) {
			defer wg.Done()
			child.SetPosition(utils.NewVector3D(val, 0, 0))
		}(float64(i))
	}
	
	wg.Wait()
	
// Позиция должна быть установлена (какое-то значение)
	pos := child.GetPosition()
	if pos.X < 0 || pos.X > 99 {
		t.Errorf("Unexpected position X=%f", pos.X)
	}
}
