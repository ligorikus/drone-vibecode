package visualization

import (
	"context"

	"drone/internal/models"
)

// Visualizer интерфейс для визуализации
type Visualizer interface {
	// Init инициализирует визуализатор
	Init() error
	
	// RenderLoop запускает цикл рендеринга
	RenderLoop(ctx context.Context) error
	
	// Close закрывает визуализатор
	Close() error
	
	// GetMainDrone возвращает главный дрон для визуализации
	GetMainDrone() *models.Drone
}
