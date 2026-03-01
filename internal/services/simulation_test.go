package services

import (
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"drone/internal/config"
)

func TestNewSimulationService(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	service := NewSimulationService(cfg, logger)

	if service == nil {
		t.Fatal("Expected service to be created")
	}
	if service.config != cfg {
		t.Error("Expected config to be set")
	}
	if service.rng == nil {
		t.Error("Expected RNG to be initialized")
	}
}

func TestSimulationServiceInit(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DroneCount = 5
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	err := service.Init()
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	mainDrone := service.GetMainDrone()
	if mainDrone == nil {
		t.Fatal("Expected main drone to be initialized")
	}

	pos := mainDrone.GetPosition()
	if pos.X != 0 || pos.Y != 10 || pos.Z != 0 {
		t.Errorf("Expected start position (0, 10, 0), got (%f, %f, %f)", pos.X, pos.Y, pos.Z)
	}
}

func TestSimulationServiceStart(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DroneCount = 5
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	if err := service.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	err := service.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Проверяем, что дочерние дроны созданы
	mainDrone := service.GetMainDrone()
	if mainDrone.GetChildCount() != cfg.DroneCount {
		t.Errorf("Expected %d children, got %d", cfg.DroneCount, mainDrone.GetChildCount())
	}

	service.Stop()
}

func TestSimulationServiceStop(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DroneCount = 5
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	if err := service.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if err := service.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Даём горутинам запуститься
	time.Sleep(50 * time.Millisecond)

	service.Stop()

	// Ждём завершения
	time.Sleep(100 * time.Millisecond)
}

func TestSimulationServiceGetFormationStats(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DroneCount = 10
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	if err := service.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if err := service.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer service.Stop()

	stats := service.GetFormationStats()

	if stats.Count != cfg.DroneCount {
		t.Errorf("Expected %d children, got %d", cfg.DroneCount, stats.Count)
	}

	if stats.AvgDistance < cfg.MinDistance || stats.AvgDistance > cfg.MaxDistance {
		t.Errorf("Expected avg distance in range [%.2f, %.2f], got %.2f",
			cfg.MinDistance, cfg.MaxDistance, stats.AvgDistance)
	}
}

func TestSimulationServiceGetConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	returnedCfg := service.GetConfig()
	if returnedCfg != cfg {
		t.Error("Expected GetConfig to return the same config")
	}
}

func TestSimulationServiceSetInput(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	input := InputState{
		Forward: true,
		Up:      true,
	}

	service.SetInput(input)

	if service.inputState.Forward != true {
		t.Error("Expected input state to be set")
	}
}

func TestSimulationServiceGetInputState(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	input := InputState{
		Backward: true,
		Left:     true,
	}

	service.SetInput(input)

	returnedInput := service.GetInputState()
	if !returnedInput.Backward || !returnedInput.Left {
		t.Error("Expected GetInputState to return the set input")
	}
}

func TestSimulationServiceUpdateMainDrone(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	if err := service.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	input := InputState{
		Forward: true,
	}
	service.SetInput(input)

	initialPos := service.GetMainDrone().GetPosition()
	service.UpdateMainDrone(0.016)
	newPos := service.GetMainDrone().GetPosition()

	// Проверяем, что позиция изменилась
	if newPos.Z >= initialPos.Z {
		t.Errorf("Expected Z to decrease when moving forward")
	}
}

func TestSimulationServiceContextCancellation(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DroneCount = 5
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	if err := service.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if err := service.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Отменяем контекст
	service.cancel()

	// Ждём завершения
	done := make(chan struct{})
	go func() {
		service.Stop()
		close(done)
	}()

	select {
	case <-done:
		// OK
	case <-time.After(2 * time.Second):
		t.Error("Service did not stop after context cancellation")
	}
}

func TestSimulationServiceWithDifferentConfigs(t *testing.T) {
	testCases := []struct {
		name       string
		droneCount int
		minDist    float64
		maxDist    float64
	}{
		{"Minimal", 1, 1.0, 2.0},
		{"Small", 10, 3.0, 5.0},
		{"Medium", 50, 5.0, 10.0},
		{"Large", 100, 8.0, 15.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{
				DroneCount:    tc.droneCount,
				MinDistance:   tc.minDist,
				MaxDistance:   tc.maxDist,
				UpdateInterval: 2 * time.Second,
			}
			logger := slog.New(slog.NewTextHandler(io.Discard, nil))

			service := NewSimulationService(cfg, logger)

			if err := service.Init(); err != nil {
				t.Fatalf("Init failed: %v", err)
			}

			if err := service.Start(); err != nil {
				t.Fatalf("Start failed: %v", err)
			}
			defer service.Stop()

			mainDrone := service.GetMainDrone()
			if mainDrone.GetChildCount() != tc.droneCount {
				t.Errorf("Expected %d children, got %d", tc.droneCount, mainDrone.GetChildCount())
			}
		})
	}
}

func TestSimulationServicePrintFormationStats(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DroneCount = 5
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	if err := service.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if err := service.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer service.Stop()

	// printFormationStats вызывается внутри Start, просто проверяем что статистика доступна
	stats := service.GetFormationStats()
	if stats.Count == 0 {
		t.Error("Expected formation stats to be available")
	}
}

func TestSimulationServiceGetRNG(t *testing.T) {
	cfg := config.DefaultConfig()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	rng := service.GetRNG()
	if rng == nil {
		t.Error("Expected RNG to be available")
	}

	// Проверяем, что RNG работает
	val := rng.Float64()
	if val < 0 || val >= 1 {
		t.Errorf("Expected RNG to return value in [0, 1), got %f", val)
	}
}

func TestSimulationServiceConcurrentAccess(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.DroneCount = 10
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	service := NewSimulationService(cfg, logger)

	if err := service.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	if err := service.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer service.Stop()

	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			_ = service.GetFormationStats()
			_ = service.GetConfig()
			_ = service.GetInputState()
			time.Sleep(time.Millisecond)
		}
		close(done)
	}()

	select {
	case <-done:
		// OK
	case <-time.After(5 * time.Second):
		t.Error("Concurrent access test timed out")
	}
}
