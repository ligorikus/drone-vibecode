//go:build no_vis

package visualization

import (
	"context"

	"drone/internal/config"
	"drone/internal/models"
	"drone/internal/services"
)

// Manager управляет визуализацией (заглушка для no_vis)
type Manager struct {
	config    *config.Config
	mainDrone *models.Drone
}

// NewManager создаёт новый менеджер визуализации
func NewManager(cfg *config.Config, mainDrone *models.Drone) *Manager {
	return &Manager{
		config:    cfg,
		mainDrone: mainDrone,
	}
}

// Init инициализирует визуализатор (заглушка)
func (m *Manager) Init() error {
	return nil
}

// SetSimulation устанавливает сервис симуляции (заглушка)
func (m *Manager) SetSimulation(sim *services.SimulationService) {}

// Run запускает цикл визуализации (заглушка)
func (m *Manager) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Close закрывает визуализатор (заглушка)
func (m *Manager) Close() error {
	return nil
}

// GetVisualizer возвращает визуализатор (заглушка)
func (m *Manager) GetVisualizer() Visualizer {
	return nil
}
