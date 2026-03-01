package models

import (
	"sync"

	"drone/internal/utils"
)

const (
	// GroundLevel уровень земли
	GroundLevel = 0.0
)

// DronePhysics представляет физические характеристики дрона
type DronePhysics struct {
	Velocity     *utils.Vector3D // Текущая скорость
	Direction    *utils.Vector3D // Вектор направления
	MaxSpeed     float64         // Максимальная скорость
	Acceleration float64         // Ускорение
	Drag         float64         // Сопротивление воздуха (0-1)
}

// Drone представляет главный дрон
type Drone struct {
	ID       int
	Position *utils.Vector3D
	Physics  *DronePhysics
	children []*ChildDrone
	mu       sync.RWMutex
}

// ChildDrone представляет дочерний дрон
type ChildDrone struct {
	ID       int
	Position *utils.Vector3D
	Physics  *DronePhysics
	mu       sync.RWMutex
}

// NewDronePhysics создаёт физические характеристики по умолчанию
func NewDronePhysics() *DronePhysics {
	return &DronePhysics{
		Velocity:     utils.Zero(),
		Direction:    utils.Zero(),
		MaxSpeed:     30.0,
		Acceleration: 15.0,
		Drag:         0.92,
	}
}

// NewDrone создаёт новый главный дрон в позиции (0, 0, 0)
func NewDrone(id int) *Drone {
	return &Drone{
		ID:       id,
		Position: utils.NewVector3D(0, 0, 0),
		Physics:  NewDronePhysics(),
		children: make([]*ChildDrone, 0),
	}
}

// NewDroneAtPosition создаёт главный дрон в указанной позиции
func NewDroneAtPosition(id int, position *utils.Vector3D) *Drone {
	return &Drone{
		ID:       id,
		Position: position.Copy(),
		Physics:  NewDronePhysics(),
		children: make([]*ChildDrone, 0),
	}
}

// NewChildDrone создаёт новый дочерний дрон
func NewChildDrone(id int, position *utils.Vector3D) *ChildDrone {
	return &ChildDrone{
		ID:       id,
		Position: position.Copy(),
		Physics:  NewDronePhysics(),
	}
}

// AddChild добавляет дочерний дрон
func (d *Drone) AddChild(child *ChildDrone) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.children = append(d.children, child)
}

// GetChildren возвращает копию слайса дочерних дронов
func (d *Drone) GetChildren() []*ChildDrone {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make([]*ChildDrone, len(d.children))
	copy(result, d.children)
	return result
}

// RangeChildren вызывает функцию для каждого дочернего дрона без создания слайса
func (d *Drone) RangeChildren(fn func(*ChildDrone) bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, child := range d.children {
		if !fn(child) {
			return
		}
	}
}

// GetPosition возвращает копию текущей позиции
func (d *Drone) GetPosition() *utils.Vector3D {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Position.Copy()
}

// SetPosition устанавливает новую позицию
func (d *Drone) SetPosition(pos *utils.Vector3D) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Position = pos.Copy()
}

// GetChildCount возвращает количество дочерних дронов
func (d *Drone) GetChildCount() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.children)
}

// GetPosition возвращает копию позиции дочернего дрона
func (cd *ChildDrone) GetPosition() *utils.Vector3D {
	cd.mu.RLock()
	defer cd.mu.RUnlock()
	return cd.Position.Copy()
}

// SetPosition устанавливает новую позицию дочернего дрона
func (cd *ChildDrone) SetPosition(pos *utils.Vector3D) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	cd.Position = pos.Copy()
}

// ApplyDirection применяет вектор направления к дрону
func (d *Drone) ApplyDirection(dir *utils.Vector3D) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Physics.Direction = dir.Copy()
}

// Accelerate ускоряет дрон в направлении вектора направления
func (d *Drone) Accelerate(deltaTime float64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Применяем ускорение в направлении движения
	if d.Physics.Direction.Length() > 0 {
		accel := d.Physics.Direction.Normalize().Multiply(d.Physics.Acceleration * deltaTime)
		d.Physics.Velocity = d.Physics.Velocity.Add(accel)
	}

	// Применяем сопротивление воздуха
	d.Physics.Velocity = d.Physics.Velocity.Multiply(d.Physics.Drag)

	// Ограничиваем максимальную скорость
	if d.Physics.Velocity.Length() > d.Physics.MaxSpeed {
		d.Physics.Velocity = d.Physics.Velocity.Normalize().Multiply(d.Physics.MaxSpeed)
	}
}

// Update обновляет позицию дрона на основе скорости
func (d *Drone) Update(deltaTime float64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Обновляем позицию на основе скорости
	displacement := d.Physics.Velocity.Multiply(deltaTime)
	d.Position = d.Position.Add(displacement)
}

// SetVelocity устанавливает скорость дрона
func (d *Drone) SetVelocity(vel *utils.Vector3D) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Physics.Velocity = vel.Copy()
}

// GetVelocity возвращает копию текущей скорости
func (d *Drone) GetVelocity() *utils.Vector3D {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Physics.Velocity.Copy()
}

// ClampY ограничивает Y координату минимальным значением
func (d *Drone) ClampY(minY float64) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.Position.Y < minY {
		d.Position.Y = minY
	}
	// Обнуляем вертикальную скорость при приземлении
	if utils.IsFloatEqual(d.Position.Y, minY) && d.Physics.Velocity.Y < 0 {
		d.Physics.Velocity.Y = 0
	}
}

// ApplyDirection применяет вектор направления к дочернему дрону
func (cd *ChildDrone) ApplyDirection(dir *utils.Vector3D) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	cd.Physics.Direction = dir.Copy()
}

// Accelerate ускоряет дочерний дрон
func (cd *ChildDrone) Accelerate(deltaTime float64) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	if cd.Physics.Direction.Length() > 0 {
		accel := cd.Physics.Direction.Normalize().Multiply(cd.Physics.Acceleration * deltaTime)
		cd.Physics.Velocity = cd.Physics.Velocity.Add(accel)
	}

	cd.Physics.Velocity = cd.Physics.Velocity.Multiply(cd.Physics.Drag)

	if cd.Physics.Velocity.Length() > cd.Physics.MaxSpeed {
		cd.Physics.Velocity = cd.Physics.Velocity.Normalize().Multiply(cd.Physics.MaxSpeed)
	}
}

// Update обновляет позицию дочернего дрона
func (cd *ChildDrone) Update(deltaTime float64) {
	cd.mu.Lock()
	defer cd.mu.Unlock()

	displacement := cd.Physics.Velocity.Multiply(deltaTime)
	cd.Position = cd.Position.Add(displacement)
}

// ClampY ограничивает Y координату дочернего дрона
func (cd *ChildDrone) ClampY(minY float64) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	if cd.Position.Y < minY {
		cd.Position.Y = minY
	}
	if utils.IsFloatEqual(cd.Position.Y, minY) && cd.Physics.Velocity.Y < 0 {
		cd.Physics.Velocity.Y = 0
	}
}

// SetVelocity устанавливает скорость дочернего дрона
func (cd *ChildDrone) SetVelocity(vel *utils.Vector3D) {
	cd.mu.Lock()
	defer cd.mu.Unlock()
	cd.Physics.Velocity = vel.Copy()
}

// GetVelocity возвращает копию текущей скорости дочернего дрона
func (cd *ChildDrone) GetVelocity() *utils.Vector3D {
	cd.mu.RLock()
	defer cd.mu.RUnlock()
	return cd.Physics.Velocity.Copy()
}
