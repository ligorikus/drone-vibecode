package services

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"drone/internal/config"
	"drone/internal/models"
	"drone/internal/utils"
)

// SimulationProvider интерфейс для доступа к симуляции из визуализатора
type SimulationProvider interface {
	GetMainDrone() *models.Drone
	GetInputState() InputState
	SetInput(InputState)
	UpdateMainDrone(deltaTime float64)
}

// SimulationService сервис симуляции дронов
type SimulationService struct {
	config     *config.Config
	mainDrone  *models.Drone
	logger     *slog.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	controller     *DroneController
	mainDroneCtrl  *MainDroneController
	inputState     InputState
	rng            *rand.Rand
}

// Проверка что SimulationService реализует SimulationProvider
var _ SimulationProvider = (*SimulationService)(nil)

// NewSimulationService создаёт новый сервис симуляции
func NewSimulationService(cfg *config.Config, logger *slog.Logger) *SimulationService {
	return &SimulationService{
		config:        cfg,
		logger:        logger,
		controller:    NewDroneController(cfg),
		mainDroneCtrl: NewMainDroneController(),
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Init инициализирует симуляцию
func (s *SimulationService) Init() error {
	// Позиция для главного дрона near origin
	startPosition := utils.NewVector3D(0, 10, 0) // Начальная высота 10 над землёй

	s.mainDrone = models.NewDroneAtPosition(0, startPosition)

	s.logger.Info("Инициализация симуляции",
		"drone_count", s.config.DroneCount,
		"min_distance", s.config.MinDistance,
		"max_distance", s.config.MaxDistance,
		"start_position", startPosition,
	)

	return nil
}

// Start запускает симуляцию
func (s *SimulationService) Start() error {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	fmt.Printf("Главный дрон %d запущен в позиции (%.2f, %.2f, %.2f)\n",
		s.mainDrone.ID,
		s.mainDrone.GetPosition().X,
		s.mainDrone.GetPosition().Y,
		s.mainDrone.GetPosition().Z,
	)
	
	// Создание дочерних дронов
	if err := s.createChildDrones(); err != nil {
		return fmt.Errorf("ошибка создания дочерних дронов: %w", err)
	}
	
	// Печать статистики формирования
	s.printFormationStats()
	
	s.logger.Info("Симуляция запущена", "children", s.mainDrone.GetChildCount())
	
	return nil
}

// createChildDrones создаёт и запускает дочерние дроны
func (s *SimulationService) createChildDrones() error {
	center := s.mainDrone.GetPosition()

	for i := 0; i < s.config.DroneCount; i++ {
		position := models.CalculateSphericalPosition(center, s.config.MinDistance, s.config.MaxDistance, s.rng)
		child := models.NewChildDrone(i, position)

		s.mainDrone.AddChild(child)

		s.wg.Add(1)
		go s.runChildDrone(child)
	}

	return nil
}

// GetRNG возвращает генератор случайных чисел для использования в других компонентах
func (s *SimulationService) GetRNG() *rand.Rand {
	return s.rng
}

// runChildDrone запускает поведение дочернего дрона
func (s *SimulationService) runChildDrone(child *models.ChildDrone) {
	defer s.wg.Done()

	// Обновляем чаще для плавного следования (16мс ~ 60 FPS)
	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	s.logger.Debug("Дочерний дрон запущен",
		"id", child.ID,
		"position", child.GetPosition(),
	)

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Debug("Дочерний дрон остановлен", "id", child.ID)
			return
		case <-ticker.C:
			s.controller.UpdatePosition(s.ctx, child, s.mainDrone)
		}
	}
}

// Stop останавливает симуляцию с graceful shutdown
func (s *SimulationService) Stop() {
	s.logger.Info("Остановка симуляции...")
	
	if s.cancel != nil {
		s.cancel()
	}
	
	// Ждём завершения всех горутин с таймаутом
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		s.logger.Info("Симуляция остановлена корректно")
	case <-time.After(5 * time.Second):
		s.logger.Warn("Таймаут остановки симуляции")
	}
}

// GetMainDrone возвращает главный дрон
func (s *SimulationService) GetMainDrone() *models.Drone {
	return s.mainDrone
}

// GetFormationStats возвращает статистику формирования
func (s *SimulationService) GetFormationStats() models.FormationStats {
	children := s.mainDrone.GetChildren()
	return models.GetFormationStats(s.mainDrone.GetPosition(), children)
}

// printFormationStats выводит статистику формирования
func (s *SimulationService) printFormationStats() {
	stats := s.GetFormationStats()
	
	fmt.Printf("\n=== Статистика формирования ===\n")
	fmt.Printf("Главный дрон: %d\n", s.mainDrone.ID)
	fmt.Printf("Дочерних дронов: %d\n", stats.Count)
	fmt.Printf("Средняя дистанция: %.2f\n", stats.AvgDistance)
	fmt.Printf("Минимальная дистанция: %.2f\n", stats.MinDistance)
	fmt.Printf("Максимальная дистанция: %.2f\n", stats.MaxDistance)
	fmt.Printf("===============================\n\n")
}

// GetConfig возвращает конфигурацию симуляции
func (s *SimulationService) GetConfig() *config.Config {
	return s.config
}

// SetInput устанавливает состояние ввода для главного дрона
func (s *SimulationService) SetInput(input InputState) {
	s.inputState = input
	s.mainDroneCtrl.SetInput(input)
}

// UpdateMainDrone обновляет позицию главного дрона на основе ввода
func (s *SimulationService) UpdateMainDrone(deltaTime float64) {
	s.mainDroneCtrl.Update(s.mainDrone, deltaTime)
}

// GetInputState возвращает текущее состояние ввода
func (s *SimulationService) GetInputState() InputState {
	return s.inputState
}
