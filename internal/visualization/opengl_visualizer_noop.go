//go:build no_vis

package visualization

import (
	"context"

	"drone/internal/config"
	"drone/internal/models"
	"drone/internal/services"
)

// OpenGLVisualizer - заглушка для сборки без визуализации
type OpenGLVisualizer struct {
	config *config.Config
}

// NewOpenGLVisualizer создаёт заглушку визуализатора
func NewOpenGLVisualizer(mainDrone *models.Drone, cfg *config.Config) *OpenGLVisualizer {
	return &OpenGLVisualizer{
		config: cfg,
	}
}

// SetSimulation устанавливает сервис симуляции (заглушка)
func (ov *OpenGLVisualizer) SetSimulation(sim *services.SimulationService) {}

// Init инициализирует ресурсы (заглушка)
func (ov *OpenGLVisualizer) Init() error {
	return nil
}

// RenderLoop запускает цикл рендеринга (заглушка)
func (ov *OpenGLVisualizer) RenderLoop(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Close закрывает ресурсы (заглушка)
func (ov *OpenGLVisualizer) Close() error {
	return nil
}

// GetMainDrone возвращает главный дрон
func (ov *OpenGLVisualizer) GetMainDrone() *models.Drone {
	return nil
}
