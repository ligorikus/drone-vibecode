//go:build !no_vis

package visualization

import (
	"context"
	"fmt"

	"drone/internal/config"
	"drone/internal/models"
	"drone/internal/services"
)

// Manager управляет визуализацией
type Manager struct {
	visualizer Visualizer
	config     *config.Config
	mainDrone  *models.Drone
}

// NewManager создаёт новый менеджер визуализации
func NewManager(cfg *config.Config, mainDrone *models.Drone) *Manager {
	return &Manager{
		config:    cfg,
		mainDrone: mainDrone,
	}
}

// Init инициализирует визуализатор
func (m *Manager) Init() error {
	if !m.config.VisualizationEnabled {
		return nil
	}

	visualizer := NewOpenGLVisualizer(m.mainDrone, m.config)
	if err := visualizer.Init(); err != nil {
		return fmt.Errorf("ошибка инициализации OpenGL визуализатора: %w", err)
	}

	m.visualizer = visualizer
	return nil
}

// SetSimulation устанавливает сервис симуляции для визуализатора
func (m *Manager) SetSimulation(sim *services.SimulationService) {
	if ov, ok := m.visualizer.(*OpenGLVisualizer); ok {
		ov.SetSimulation(sim)
	}
}

// Run запускает цикл визуализации
func (m *Manager) Run(ctx context.Context) error {
	if !m.config.VisualizationEnabled || m.visualizer == nil {
		// Ждём контекст если визуализация отключена
		<-ctx.Done()
		return nil
	}

	return m.visualizer.RenderLoop(ctx)
}

// Close закрывает визуализатор
func (m *Manager) Close() error {
	if m.visualizer != nil {
		return m.visualizer.Close()
	}
	return nil
}

// GetVisualizer возвращает визуализатор
func (m *Manager) GetVisualizer() Visualizer {
	return m.visualizer
}
